package utils

import (
	"log"
	"regexp"
)

func SayHi() {
	log.Println("hi, this is utils")
}

func RegexpFindString(s, regex string) string {
	re := regexp.MustCompile(regex)
	str := re.FindString(s)
	//log.Println(s, regex, str)
	return str
}
