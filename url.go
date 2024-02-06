package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//http://172.0.0.1:8888/hello?admin=admin&age=18
func GetUrlArgTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s", r.Method)
	log.Printf("path: %s", r.URL.Path)
	log.Printf("args: %s", r.URL.RawQuery)
	// 注意FormValue 可以获取get参数还有post参数
	log.Printf("admin值: %s, age值: %s", r.FormValue("admin"), r.FormValue("age"))
	// PostFormValue 只会获取post传递过来的参数
	log.Printf("username值: %s, age值: %s", r.PostFormValue("username"), r.PostFormValue("age"))
}

func GetFullUrl(r *http.Request, host, port string) (Url string) {
	/*
		log.Println(r.Proto)      // output:HTTP/1.1
		log.Println(r.TLS)        // output: <nil>
		log.Println(r.Host)       // output: localhost:9090
		log.Println(r.RequestURI) // output: /index?id=1
	*/
	//ss := strings.Split(r.Host, ":")
	t := fmt.Sprintf("%s:%s", host, port)

	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	//http://localhost:9090/index?id=1
	return strings.Join([]string{scheme, t, r.RequestURI}, "")
}
