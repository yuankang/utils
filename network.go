package utils

import (
	"encoding/json"
	"log"
	"net"
	"unsafe"
)

/*
## 彻底弄明白字节序 ByteOrder
什么情况下要注意字节序?
处理器读取数据的时候, 必须指定数据的字节序, 否则读取的结果不是真实数据的值
其他情况都不用考虑 即使是向外发送数据时 也不用考虑字节序

为什么会有字节序?
计算机处理字节的时候, 不知道什么是高位字节, 什么是低位字节, 它只知道按顺序读取字节
如果是按大端字节序读取, 先读到的就是高位字节, 后读到的就是低位字节
如果是按小端字节序读取, 先读到的就是低位字节, 后读到的就是高位字节
举例说明: 对4字节的数据 0x11223344 来说
按大端字节序读取, 0x11是高位 0x44是低位, 结果就是 0x11223344 (符合人类阅读习惯)
按小端字节序读取, 0x11是低位 0x44是高位, 结果就是 0x44332211 (有利于cpu运算)

那么 写入的真实数据是 0x11223344 还是 0x44332211 呢?
这就要求 读取数据的字节序 要和 写入数据的字节序 一致, 否则读取的结果不是真实数据的值
以大端字节序写入 就要用大端字节序读取; 以小端字节序写入 就要用小端字节序读取;
调用系统函数读取或写入数据时 一般都可以指定字节序
对于网络数据, 一般发送方和接收方会 约定字节序, 比如: rtmp协议 发收数据 默认使用大端字节序

通常电脑主机使用intel cpu, 默认使用小端字节序, intel cpu 计算能力强 功耗高
通常网络设备使用  arm cpu, 默认使用大端字节序, arm   cpu 计算能力弱 功耗低

centos获取cpu默认字节序的命令
[www:~]# lscpu | grep -i byte
Byte Order:            Little Endian

大端小端名词的由来? BigEndian LittleEndian
大端小端名词来源于一个有趣的故事, 故事出自Jonathan Swift的《格利佛游记》
两个国家为了 吃鸡蛋时, 先打破鸡蛋的哪一端, 而爆发了36个月的战争 故事其实在讽刺当时英国和法国之间持续的冲突。
Danny Cohen一位网络协议的开创者，第一次使用这两个术语指代字节顺序，后来就大家广泛接受。
*/
func IsLittleEndian() bool {
	//底层原理是什么?
	var i int32 = 0x11223344
	//将int32类型的指针转换为byte类型的指针
	ip := unsafe.Pointer(&i)
	//bp := (*[n]byte)(ip) //???
	//return (bp[0] == 0x44)
	bp := (*byte)(ip)
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
