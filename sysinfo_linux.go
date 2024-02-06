package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

func UseShellCmd(cmd string) ([]byte, error) {
	//s := fmt.Sprintf("cat /sys/class/net/%s/statistics/%s", name, inout)
	c := exec.Command("/bin/sh", "-c", cmd)
	b, err := c.CombinedOutput()
	return b, err
}

func GetHostOs() string {
	s := fmt.Sprintf("cat /etc/redhat-release")
	b, err := UseShellCmd(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	s = strings.TrimSpace(string(b))
	return s
}

func GetKernel() string {
	s := fmt.Sprintf("uname -r")
	b, err := UseShellCmd(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	s = strings.TrimSpace(string(b))
	return s
}

func GetCpuInfo() string {
	s := fmt.Sprintf("cat /proc/cpuinfo | grep name | head -n 1 | cut -f2 -d:")
	b, err := UseShellCmd(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	s = strings.TrimSpace(string(b))
	return s
}

func GetCpuUsed() string {
	// 不加-b 会报错 top: failed tty get
	s := fmt.Sprintf("top -bn 1 | grep Cpu")
	b, err := UseShellCmd(s)
	if err != nil {
		log.Println(err)
		log.Println(string(b))
		return ""
	}
	s = strings.TrimSpace(string(b))
	s = strings.Replace(s, " ", "", -1)
	return s
}

func GetTaskInfo() string {
	// 不加-b 会报错 top: failed tty get
	s := fmt.Sprintf("top -bn 1 | grep Tasks")
	b, err := UseShellCmd(s)
	if err != nil {
		log.Println(err)
		log.Println(string(b))
		return ""
	}
	s = strings.TrimSpace(string(b))
	s = strings.Replace(s, " ", "", -1)
	return s
}

func GetLoadInfo() string {
	cmd := exec.Command("/bin/sh", "-c", `cat /proc/loadavg | awk '{print $1}'`)
	b, err := cmd.CombinedOutput()
	if err != nil || len(b) == 0 {
		log.Println(err, len(b))
		return ""
	}
	return string(b[:len(b)-1])
}

func GetMemInfo() (uint64, uint64) {
	var total, used uint64
	/*
		si := new(syscall.Sysinfo_t)
		err := syscall.Sysinfo(si)
		if err != nil {
			log.Println(err)
			return 0, 0
		}
		// total 单位为字节, 所以 除两次1024的到 MB
		total = si.Totalram
		used = total - si.Freeram
	*/
	return total / 1024 / 1024, used / 1024 / 1024
}

func GetDiskInfo(path string) (uint64, uint64) {
	var total, used uint64
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	// total 单位为字节, 所以 除两次1024的到 MB
	total = fs.Blocks * uint64(fs.Bsize)
	used = total - fs.Bfree*uint64(fs.Bsize)
	return total / 1024 / 1024, used / 1024 / 1024
}
