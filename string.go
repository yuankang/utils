package utils

import (
	"log"
	"path"
	"strconv"
	"strings"
)

//清除字符串中所有空字符
//空字符 包括: 空格, Tab(\t), 回车(\r), 换行(\n)
//方法1:
//str = strings.Replace(str, " ", "", -1)  //去除空格
//str = strings.Replace(str, "\n", "", -1) //去除换行
//方法2:
//strings.Fields()	按空字符把字符串分割, 返回字符串数组
//strings.Join()	将字符串数组重新拼接
//strings.Join(strings.Fields(input), "")

//清除字符串两端的空格
//TrimSpace() 只去除两端的空格, 中间的无法去除
//s1 := strings.TrimSpace(s0)

//string to int
func Str2Int() {
	s := "66"

	//第一个参数10是base，表示十进制
	//第二个参数0/64表示bitSize，也就是整数类型
	//string转int
	n, err := strconv.ParseInt(s, 10, 0)
	log.Println(n, err)

	//string转uint64
	n, err = strconv.ParseInt(s, 10, 64)
	log.Println(n, err)

	//Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
	//func Atoi(s string) (int, error)
	i, err := strconv.Atoi(s)
	log.Println(i, err)
}

//从路径中获取相关信息
///usr/local/sliveserver/streamlog/GSP3bnx69BgxI-avEc0oE4C4/publish_rtmp_20230113.log
func GetPathInfo(pfn string) {
	//获取文件名称带后缀
	fn := path.Base(pfn)
	//获取文件的后缀(文件类型)
	fType := path.Ext(fn)
	//获取文件后缀名
	fExt := path.Ext(fn)
	//得到文件名不带后缀
	ofn := strings.TrimSuffix(fn, fExt)

	log.Println(pfn)
	log.Println(fn, fType, fExt, ofn)
}
