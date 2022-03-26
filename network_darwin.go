package utils

import "log"

/************************************************************/
/* 获取公网ip
/************************************************************/
// 获取默认路由对应的网卡名
func GetDefaultIface() (string, error) {
	log.Println("GetDefaultIface() TODO")
	return "", nil
}

// 方法1: 获取默认路由对应的网卡名, 通过网卡名获取出口ip
// 方法2: 发起一个udp连接，通过socket获取出口ip
// name 表示网卡名, 如: eth0
func GetIpOuter(name string) (string, error) {
	log.Println("GetIpOuter() TODO")
	return "", nil
}

/************************************************************/
/* 获取实时流量
/************************************************************/
// name 表示网卡名, 如: eth0
// io 表示输入输出, rx 表示输入(接收) tx表示输出(发送)
// 获取某个时间点网卡接收或发送的数据量, 单位bit
func GetNetDataSize(name, io string) (int, error) {
	log.Println("GetNetDataSize() TODO")
	return 0, nil
}

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
func GetNetBandwidth(name string) (string, error) {
	log.Println("GetNetBandwidth() TODO")
	return "", nil
}
