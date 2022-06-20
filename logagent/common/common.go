package common

import (
	"fmt"
	"net"
	"strings"
)

// Config 要收集日志配置项的结构体
type Config struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func GetLocalCpIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:23")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("local computer ip is:", strings.Split(localAddr.String(), ":")[0])
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
