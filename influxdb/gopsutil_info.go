package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"log"
	"time"
)

var (
	cli                     client.Client
	lastNetIOStartTimeStamp int64    // 上一次获取网络io的时间
	lastNetIoInfo           *NetInfo // 上一次网络io的数据
)

const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

// SysInfo 系统信息
type SysInfo struct {
	InfoType string
	IP       string
	Data     interface{}
}

// CpuInfo cpu信息
type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

// MemInfo 内存信息
type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

// UsageStat 分区信息
type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	PartitionUsageStat map[string]*disk.UsageStat
}

// IOStat 网络io
type IOStat struct {
	BytesSent       uint64
	BytesRecv       uint64
	PacketsSent     uint64
	PacketsRecv     uint64
	BytesSentRate   float64 `json:"bytes_sent_rate"`
	BytesRecvRate   float64 `json:"bytes_recv_rate"`
	PacketsSentRate float64 `json:"packets_sent_rate"`
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}

// NetInfo 网络信息
type NetInfo struct {
	NetIOCountersStat map[string]*IOStat
}

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

// writesPointsDiskToDB
func writesPointsDiskToDB(info *DiskInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range info.PartitionUsageStat {
		tags := map[string]string{"path": k}
		fields := map[string]interface{}{
			"total":               int64(v.Total),
			"free":                int64(v.Free),
			"used":                int64(v.Used),
			"used_percent":        v.UsedPercent,
			"inodes_total":        int64(v.InodesTotal),
			"inodes_used":         int64(v.InodesUsed),
			"inodes_free":         int64(v.InodesFree),
			"inodes_used_percent": v.InodesUsedPercent,
		}
		pt, err := client.NewPoint("disk", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert disk info success")
}

// writesPointsNetToDB
func writesPointsNetToDB(info *NetInfo) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range info.NetIOCountersStat {
		tags := map[string]string{"name": k}
		fields := map[string]interface{}{
			"bytes_sent_rate":   v.BytesSentRate,
			"bytes_recv_rate":   v.BytesRecvRate,
			"packets_sent_rate": v.PacketsSentRate,
			"packets_recv_rate": v.PacketsRecvRate,
		}
		pt, err := client.NewPoint("net", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert net info success")
}

// cpu info（cpu信息）
func getCpuInfo() {
	var cpuInfo = new(CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v type of percent:%T \n", percent, percent)
	// 写入到influxdb数据库
	cpuInfo.CpuPercent = percent[0]
	writesPointsCpuToDB(cpuInfo)
}

//  mem info（内存信息）
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

// Disk info（磁盘信息）
func getDiskInfo() {
	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 16),
	}
	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		// 拿到每一个分区
		usageData, err := disk.Usage(part.Mountpoint) // 传挂载点
		if err != nil {
			fmt.Printf("get %s part mount_point failed, err:%v ", part.Mountpoint, err)
			continue
		}
		diskInfo.PartitionUsageStat[part.Mountpoint] = usageData
		writesPointsDiskToDB(diskInfo)
	}
}

// net info（网络信息）
func getNetInfo() {
	var netInfo = &NetInfo{
		NetIOCountersStat: make(map[string]*IOStat, 8),
	}
	// 获取当前时间
	currenttime := time.Now().Unix()
	netIOs, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("get net io failed, err: ", err)
		return
	}
	for _, netIO := range netIOs {
		var ioStat = new(IOStat)
		ioStat.BytesSent = netIO.BytesSent
		ioStat.BytesRecv = netIO.BytesRecv
		ioStat.PacketsSent = netIO.PacketsSent
		ioStat.PacketsRecv = netIO.PacketsRecv

		// 将具体网卡数据的isStat变量添加到map中
		netInfo.NetIOCountersStat[netIO.Name] = ioStat

		// 开启计算网卡相关数据
		// 第一次不进行计算
		if lastNetIOStartTimeStamp == 0 || lastNetIoInfo == nil {
			continue
		}
		// 时间间隔
		interval := currenttime - lastNetIOStartTimeStamp
		ioStat.BytesSentRate = (float64(ioStat.BytesSent) - float64(lastNetIoInfo.NetIOCountersStat[netIO.Name].BytesSent)) / float64(interval)
		ioStat.BytesRecvRate = (float64(ioStat.BytesRecv) - float64(lastNetIoInfo.NetIOCountersStat[netIO.Name].BytesRecv)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsSent) - float64(lastNetIoInfo.NetIOCountersStat[netIO.Name].PacketsSent)) / float64(interval)
		ioStat.PacketsSentRate = (float64(ioStat.PacketsRecv) - float64(lastNetIoInfo.NetIOCountersStat[netIO.Name].PacketsRecv)) / float64(interval)
	}
	// 更新全局记录的上一次采集网卡数据
	lastNetIOStartTimeStamp = currenttime // 更新时间
	lastNetIoInfo = netInfo
	// 发送到influxDB
	writesPointsNetToDB(netInfo)
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
		getMemInfo()
		getDiskInfo()
		getNetInfo()
	}
}
