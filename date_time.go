package utils

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

//2022
func GetY() string {
	t := time.Now()
	return fmt.Sprintf("%d", t.Year())
}

//20220102
func GetYMD() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

//202201021030
func GetYMDHM() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

//20220102103020
func GetYMDHMS() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//2022年01月02日10:30:20
func GetYMDHMS0() string {
	t := time.Now()
	return fmt.Sprintf("%d年%02d月%02d日%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//2022-01-02T10:30:20.500
func GetYMDHMS1() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%03d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000)
}

func TimestampToTimestring(ts int64) string {
	t := time.Unix(ts, 0)
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func TimestringToTimestamp(s string) int64 {
	if len(s) < 14 {
		return 0
	}

	// s = 20201007143437
	ss := fmt.Sprintf("%s-%s-%s %s:%s:%s",
		s[:4], s[4:6], s[6:8], s[8:10], s[10:12], s[12:])

	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, ss, loc)
	if err != nil {
		log.Println(err)
		return 0
	}
	return theTime.Unix()
}

//func GetNowUnix(dw string) int64 {}
func GetTimestamp(dw string) int64 {
	switch dw {
	case "s":
		return time.Now().Unix() // 秒, 长度10位
	case "ms":
		return time.Now().UnixNano() / 1e6 // 毫秒, 长度13位
	case "ns":
		return time.Now().UnixNano() // 纳秒, 长度19位
	}
	return time.Now().Unix()
}

func TimespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
