package utils

import (
	"log"
	"net"
	"os/exec"
)

/************************************************************/
/* 获取公网ip
/************************************************************/
// 获取默认路由对应的网卡名
func GetDefaultIface() (string, error) {
	// route 返回有可能慢，大约要10秒左右
	cmd := exec.Command("/bin/sh", "-c", `route | grep default | awk '{print $8}'`)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return "", err
	}

	iface := string(b[:len(b)-1])
	return iface, nil
}

// 方法1: 获取默认路由对应的网卡名, 通过网卡名获取出口ip
// 方法2: 发起一个udp连接，通过socket获取出口ip
// name 表示网卡名, 如: eth0
func GetIpOuter(name string) (string, error) {
	var err error
	if name == "" {
		name, err = GetDefaultIface()
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Println(err)
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		log.Println(err)
		return "", err
	}

	var ip string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ipnet.IP.To4() != nil && ok {
			ip = ipnet.IP.String()
			break
		}
	}
	return ip, nil
}

/************************************************************/
/* 获取实时流量
/************************************************************/
// name 表示网卡名, 如: eth0
// io 表示输入输出, rx 表示输入(接收) tx表示输出(发送)
// 获取某个时间点网卡接收或发送的数据量, 单位bit
/*
func GetNetDataSize(name, io string) (int, error) {
	s := fmt.Sprintf("cat /sys/class/net/%s/statistics/%s_bytes", name, inout)

	cmd := exec.Command("/bin/sh", "-c", s)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	ret, err := strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 返回的是字节数 是大B，所以要乘上8 转换为小b
	return ret * 8, nil
}
*/

// 返回值 是 单位时间内网络流量 Mbps
func GetNetFlowRate(t1, t2, sec int) int {
	if t1 == 0 || t2 == 0 {
		return 0
	}
	fr := (t2 - t1) / 1024 / sec
	return fr
}

/************************************************************/
/* 获取网卡带宽 Mb
/************************************************************/
// name 表示网卡名, 如: eth0
/*
func GetNetBandwidth(name string) error {
	iface, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range iface {
		if (v.Flags & net.FlagUp) == 0 {
			continue
		}

		addrs, _ := v.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ipnet.IP.To4() != nil && ok {
				log.Printf("name:%s, ip:%s, mac:%s", v.Name, v.HardwareAddr.String(), ni.Mac)
			}
		}
	}
	return nil
}
*/
