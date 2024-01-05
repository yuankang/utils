package utils

import (
	"encoding/json"
	"log"
	"net"
	"unsafe"
)

/*
## 彻底弄懂字节序
为什么要有字节序? 什么情况下要区分字节序?
首先大家好明白, 只有读取的时候, 才必须区分字节序, 其他情况都不用考虑
处理器读取数据时, 必须知道写入数据时的字节序, 否则读取的值和原值不一致

计算机处理字节的时候, 不知道什么是高位字节, 什么是低位字节, 它只知道按顺序读取
如果用大端字节序读数据, 先读到的就是高位字节, 后读到的就是低位字节
如果用小端字节序读数据, 先读到的就是低位字节, 后读到的就是高位字节
举例来说, 数值 0x11223344 使用四个字节储存
如果用大端字节序读数据, 0x11为高位 0x44位低位 结果为0x11223344
如果用小端字节序读数据, 0x11为低位 0x44位高位 结果为0x44332211

那么写入的时候 真正的数据是 0x11223344 还是 0x44332211 呢?
这就要求, 读取数据 要和 写入数据 的顺序一致, 否则读取的结果就是错误的
以小端方式写, 就要以小端方式读; 以大端方式写, 就要以大端方式读;
一般情况下 电脑主机的cpu都是以 小端方式读取数据的
一般情况下 网络设备的cpu都是以 大端方式读取数据的

centos判断cpu大小端的命令
[www:~]# lscpu | grep -i byte
Byte Order:            Little Endian
*/
func IsLittleEndian() bool {
	//这个函数 原理是什么?
	var i int32 = 0x11223344
	//将int32类型的指针转换为byte类型的指针
	ip := unsafe.Pointer(&i)
	bp := (*byte)(ip)
	//b==0x11表示用大端读; b==0x44表示用小端读;
	b := *bp
	return (b == 0x44)
}

/************************************************************/
/* 通过文件获取网络信息
/************************************************************/
/* 虚拟机网卡上之后内网ip 没有公网ip, 只能通过配置文件获取公网ip
配置文件 建议放到 cat /etc/net.conf
{
    "IpOuter":"192.168.0.106",
    "IpInner":"192.168.0.106",
    "===BandwidthMax===":"单位是 Mb 兆比特",
	"BandwidthMax":500
} */
type NetInfo struct {
	IpOuter   string
	IpInner   string
	Bandwidth int
}

func GetNetInfoFromFile(file string) (NetInfo, error) {
	var ni NetInfo

	if file == "" {
		file = "/etc/net.conf"
	}

	b, err := ReadAllFile(file)
	if err != nil {
		log.Println(err)
		return ni, err
	}

	err = json.Unmarshal(b, &ni)
	if err != nil {
		log.Println(err)
		return ni, err
	}
	return ni, nil
}

/************************************************************/
/* 获取内网ip
/************************************************************/
func IsInnerIp(ip string) bool {
	IP := net.ParseIP(ip)
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}

	ip4 := IP.To4()
	if ip4 != nil {
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

func GetIpInner() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		log.Println(err)
		return "", err
	}

	var ip string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ipnet.IP.To4() != nil && ok {
			ip = ipnet.IP.String()
			if IsInnerIp(ip) {
				break
			}
		}
	}
	return ip, nil
}

/************************************************************/
/* 获取所有网卡信息
/************************************************************/
/* 2022/03/26 13:18:42 name:lo0, ip:127.0.0.1, mac:
2022/03/26 13:18:42 name:lo0, ip:::1, mac:
2022/03/26 13:18:42 name:lo0, ip:fe80::1, mac:
2022/03/26 13:18:42 name:en6, ip:fe80::aede:48ff:fe00:1122, mac:ac:de:48:00:11:22
2022/03/26 13:18:42 name:en0, ip:fe80::867:f9af:7e88:1dd4, mac:f0:18:98:0b:43:bd
2022/03/26 13:18:42 name:en0, ip:192.168.0.107, mac:f0:18:98:0b:43:bd
2022/03/26 13:18:42 name:awdl0, ip:fe80::7c58:a4ff:fe0c:e439, mac:7e:58:a4:0c:e4:39
2022/03/26 13:18:42 name:llw0, ip:fe80::7c58:a4ff:fe0c:e439, mac:7e:58:a4:0c:e4:39
2022/03/26 13:18:42 name:utun0, ip:10.3.206.1, mac:
2022/03/26 13:18:42 name:utun1, ip:fe80::8b85:544a:fdc9:5baa, mac:
2022/03/26 13:18:42 name:utun2, ip:fe80::ba70:80df:c3ab:6aaf, mac:
2022/03/26 13:18:42 name:utun7, ip:0.0.1.1, mac: */
func GetNetInterface() {
	iface, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range iface {
		if (v.Flags & net.FlagUp) == 0 {
			continue
		}

		addrs, _ := v.Addrs()
		for _, addr := range addrs {
			ipnet, _ := addr.(*net.IPNet)
			//if ipnet.IP.To4() != nil && ok {
			log.Printf("name:%s, ip:%s, mac:%s", v.Name, ipnet.IP.String(), v.HardwareAddr.String())
			//}
		}
	}
}
