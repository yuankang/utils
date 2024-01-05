package utils

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var SecretKey = "a5c02a3b57a78ef78fbca4f650029323"

type UrlInfo struct {
	remoteaddr string
	url        string
	stream_id  string
	timestamp  string
	auth       string
}

func ParseUrl(r *http.Request) (UrlInfo, error) {
	var ui UrlInfo
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return ui, err
	}

	if len(r.Form["stream"]) < 1 {
		err = fmt.Errorf("can't find stream in url")
		log.Println(err)
		return ui, err
	}
	if len(r.Form["timestamp"]) < 1 {
		err = fmt.Errorf("can't find timestamp in url")
		log.Println(err)
		return ui, err
	}
	if len(r.Form["auth"]) < 1 {
		err = fmt.Errorf("can't find auth in url")
		log.Println(err)
		return ui, err
	}

	ui.remoteaddr = r.RemoteAddr
	ui.url = r.URL.String()
	ui.stream_id = r.Form["stream"][0]
	ui.timestamp = r.Form["timestamp"][0]
	ui.auth = r.Form["auth"][0]
	return ui, nil
}

func CheckUrl(url UrlInfo) (bool, error) {
	md5str := Md5Sum(url.timestamp + SecretKey)
	if url.auth == md5str {
		log.Println(url.auth, md5str)
		return true, nil
	}
	return false, fmt.Errorf("%s(calc md5) != %s(url md5)", md5str, url.auth)
}

/*
大端字节序：高位字节在前，低位字节在后，这是人类读写数值的方法。
小端字节序：低位字节在前，高位字节在后，即以0x1122形式储存。
举例来说，数值0x2211使用两个字节储存：高位字节是0x22，低位字节是0x11。

首先，为什么会有小端字节序？
答案是，计算机电路先处理低位字节，效率比较高，因为计算都是从低位开始的。所以，计算机的内部处理都是小端字节序。
但是，人类还是习惯读写大端字节序。所以，除了计算机的内部处理，其他的场合几乎都是大端字节序，比如网络传输和文件储存。

计算机处理字节序的时候，不知道什么是高位字节，什么是低位字节。它只知道按顺序读取字节，先读第一个字节，再读第二个字节。
如果是大端字节序，先读到的就是高位字节，后读到的就是低位字节。小端字节序正好相反。
理解这一点，才能理解计算机如何处理字节序。

字节序的处理，就是一句话：
"只有读取的时候，才必须区分字节序，其他情况都不用考虑。"
处理器读取外部数据的时候，必须知道数据的字节序，将其转成正确的值。然后，就正常使用这个值，完全不用再考虑字节序。
即使是向外部设备写入数据，也不用考虑字节序，正常写入一个值即可。外部设备会自己处理字节序的问题。
*/
func IsLittleEndianBak() bool {
	var i int32 = 0x01020304

	// 下面这两句是为了将int32类型的指针转换为byte类型的指针
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)

	b := *pb // 取得pb位置对应的值

	// 由于b是byte类型的,最多保存8位,那么只能取得开始的8位
	// 小端: 04 03 02 01
	// 大端: 01 02 03 04
	return (b == 0x04)
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
