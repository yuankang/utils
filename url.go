package utils

import (
	"fmt"
	"net/http"
	"strings"
)

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
