package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

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
func writesPointsToDB(percent []float64) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": percent[0],
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
	log.Println("insert success")
}

// 获取 cpu 的信息
func getCpuInfo() {
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v type of percent:%T \n", percent, percent)
	// 写入到influxdb数据库
	writesPointsToDB(percent)
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
		getCpuInfo()
	}
}
