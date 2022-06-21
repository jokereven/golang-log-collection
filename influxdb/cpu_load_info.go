package main

import (
	"fmt"
	"github.com/shirou/gopsutil/load"
	"time"
)

// 获取CPU负载信息
func getCpuLoad() {
	info, _ := load.Avg()
	fmt.Printf("%v\n", info)
}

func main() {
	CpuLoad := time.Tick(1 * time.Second)
	for {
		<-CpuLoad
		getCpuLoad()
	}
}
