package main

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
)

func getNetInfo() {
	info, _ := net.IOCounters(true)
	for index, v := range info {
		fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	}
}

func main() {
	getNetInfo()
}
