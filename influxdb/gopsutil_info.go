package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

type SysInfo struct {
	InfoType string
	IP       string
	Data     interface{}
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

var (
	cli client.Client
)

// 使用 gopustuil 监控cpu的情况并写入influxdb

// 连接 influxdb
func connInfluxDB() (err error) {
	cli, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	return
}

// insert

//  writesPointsCpuToDB
func writesPointsCpuToDB(info *CpuInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": info.CpuPercent,
	}

	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert cpu info success")
}

// writesPointsMemToDB
func writesPointsMemToDB(info *MemInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"mem_total":        int64(info.Total),
		"mem_available":    int64(info.Available),
		"mem_used":         int64(info.Used),
		"mem_used_percent": info.UsedPercent,
	}

	pt, err := client.NewPoint("memory", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert mem info success")
}

// 获取 cpu 的信息
func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v type of percent:%T \n", percent, percent)
	// 写入到influxdb数据库
	cpuInfo.CpuPercent = percent[0]
	writesPointsCpuToDB(cpuInfo)
}

// mem info
func getMemInfo() {
	var MemInfo = new(MemInfo)
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("get mem info failed, err:", err)
		return
	}

	MemInfo.Total = info.Total
	MemInfo.Available = info.Available
	MemInfo.Used = info.Used
	MemInfo.UsedPercent = info.UsedPercent
	writesPointsMemToDB(MemInfo)
}

func main() {
	err := connInfluxDB()
	if err != nil {
		fmt.Println("connect influxdb failed, err: ", err)
		return
	}
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker
		//getCpuInfo()
		getMemInfo()
	}
}
