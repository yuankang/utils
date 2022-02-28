package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type LiveNet struct {
	IpOuter string
	IpInner string
	BwTotal uint64
}

/************************************************************/
/* ip address
/************************************************************/
// 通过配置文件 /etc/net.config 获取 公网和内网ip
// [centos78:livepower]# cat net.config
// IpOuter=62.138.21.55
// IpInner=192.168.2.200
// BwTotal=10240
func GetIpFromFile0(file string) (string, string, uint64) {
	var ipo, ipi string
	var bwt int64

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return "", "", 0
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		str := s.Text()
		if strings.Contains(str, "IpOuter") {
			ss := strings.Split(str, "=")
			if len(ss) != 2 {
				continue
			}
			ipo = ss[1]
		} else if strings.Contains(str, "IpInner") {
			ss := strings.Split(str, "=")
			if len(ss) != 2 {
				continue
			}
			ipi = ss[1]
		} else if strings.Contains(str, "BwTotal") {
			ss := strings.Split(str, "=")
			if len(ss) != 2 {
				continue
			}
			bwt, _ = strconv.ParseInt(ss[1], 10, 64)
		}
	}
	if err := s.Err(); err != nil {
		log.Println(err)
	}

	//log.Println(ipo, ipi, bwt)
	return ipo, ipi, uint64(bwt)
}

func GetIpFromFile(file string) (string, string, uint64) {
	var ln LiveNet
	s, err := ReadAllFile(file)
	if err != nil {
		log.Println(err)
		return "", "", 0
	}

	err = json.Unmarshal([]byte(s), &ln)
	if err != nil {
		log.Println(err)
		return "", "", 0
	}
	log.Println(ln)
	return ln.IpOuter, ln.IpInner, ln.BwTotal
}

// get outer(public) ip
// 方法1: 获取默认路由的iface, 听过iface获取出口ip
// 方法2: 发起一个udp连接，通过socket获取出口ip
func GetDefaultIface() string {
	// route 返回有可能慢，大约要10秒左右
	cmd := exec.Command("/bin/sh", "-c", `route | grep default | awk '{print $8}'`)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b[:len(b)-1])
}

func GetIpOuter(name string) string {
	var ip string
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Println(err)
		return ""
	}
	addrs, _ := iface.Addrs()
	for _, addr := range addrs {
		ipnet, _ := addr.(*net.IPNet)
		if ipnet.IP.To4() != nil {
			ip = ipnet.IP.String()
			break
		}
	}
	return ip
}

// get inner ip
func IsInnerIp(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		if ip4[0] == 10 {
			return true
		} else if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		} else if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	}
	return false
}

func GetIpInner() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		fmt.Println(err)
		return ""
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ipnet.IP.To4() != nil && ok && IsInnerIp(ipnet.IP) {
			ip = ipnet.IP.String()
			break
		}
	}
	return ip
}

/************************************************************/
/* bandwidth
/************************************************************/
func GetBwInfo(name, inout string) int {
	// rx_packets  rx_bytes  rx_dropped
	s := fmt.Sprintf("cat /sys/class/net/%s/statistics/%s_bytes", name, inout)
	cmd := exec.Command("/bin/sh", "-c", s)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return 0
	}
	ret, err := strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		log.Println(err)
		return 0
	}
	// 返回的是字节数 是大B，所以要乘上8 转换为小b
	return ret * 8
}

func CalcBwRate(n1, n2, time int) int {
	if n1 == 0 || n2 == 0 {
		return 0
	}
	r := (n2 - n1) / 1024 / time
	return r
}

// 注意：程序调用时用 awk '{printf $2}'; 命令行执行时用 awk '{printf $3}';
// 虚拟机无法使用 ethtool 和 lspci, 公司的虚拟机网卡 默认10G(小b) 默认不限速
func GetBwTotal(name string) string {
	s := fmt.Sprintf("ethtool %s | grep Speed | awk '{printf $2}' | awk -F M '{printf $1}'", name)
	cmd := exec.Command("/bin/sh", "-c", s)
	b, err := cmd.CombinedOutput()
	if err != nil || len(b) == 0 {
		log.Println(err, len(b))
		return ""
	}
	return string(b[:len(b)]) + "000"
}

/************************************************************/
/* interface
/************************************************************/
type NetInterface struct {
	Name string
	Mac  string
	Ip   string
}

func GetNetInterface() (NetInterface, error) {
	var ni NetInterface

	iface, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return ni, err
	}

	for _, v := range iface {
		if (v.Flags & net.FlagUp) != 0 {
			addrs, _ := v.Addrs()

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ni.Name = v.Name
						ni.Mac = v.HardwareAddr.String()
						ni.Ip = ipnet.IP.String()
						return ni, nil
					}
				}
			}
		}
	}
	return ni, err
}
