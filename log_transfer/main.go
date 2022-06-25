package main

import (
	"fmt"
	"github.com/jokereven/log_transfer/es"
	"github.com/jokereven/log_transfer/kafka"
	"github.com/jokereven/log_transfer/model"
	"gopkg.in/ini.v1"
)

// log transfer
// 从 kafka消费日志数据, 发送到es

func main() {
	// 1. 加载配置文件
	var cfg = new(model.TransferConfig)
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		// %v 默认输出 %s 输出字符串
		fmt.Printf("loading config file failed, err: %v", err)
		return
	}
	fmt.Printf("loading config file success...\n")

	// 2. 连接es
	err = es.Init(cfg.EsConfig.Address, cfg.EsConfig.Index, cfg.EsConfig.MaxChanSize, cfg.EsConfig.GoroutineNum)
	if err != nil {
		fmt.Printf("connect to es failed, err: %v\n", err)
		return
	}
	fmt.Printf("init es success...\n")

	// 3. 连接kafka
	err = kafka.Init([]string{cfg.KafkaConfig.Address}, cfg.KafkaConfig.Topic)
	if err != nil {
		fmt.Printf("connect to kafka failed, err: %v\n", err)
		return
	}
	fmt.Printf("init kafka success...\n")

	select {}
}
