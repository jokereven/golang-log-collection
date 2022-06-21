package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
)

// host info
func getHostInfo() {
	hInfo, _ := host.Info()
	fmt.Printf("host info:%v uptime:%v boottime:%v\n", hInfo, hInfo.Uptime, hInfo.BootTime)
}

func main() {
	getHostInfo()
}
