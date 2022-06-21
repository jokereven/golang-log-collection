package main

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
)

// mem info
func getMemInfo() {
	memInfo, _ := mem.VirtualMemory()
	fmt.Printf("mem info:%v\n", memInfo)
}

func main() {
	getMemInfo()
}