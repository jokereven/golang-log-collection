package main

import (
	"fmt"
	"github.com/jokereven/golang-log-collection/logagent/kafka"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// 日志收集客户端
// 类似的开源项目filebeat
// 将指定Path下的日志收集, 发送到kafka

type AppConfig struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type CollectConfig struct {
	LogFilePath string `ini:"log_file_path"`
}

func main() {
	// 0. 读取配置文件 `go-ini`
	//（1）通过直接加载配置文件的方法
	/*	cfg, err := ini.Load("./conf/config.ini")
		if err != nil {
			logrus.Error("read config file failed, err: ", err)
			fmt.Println("read config file failed, err: ", err)
			return
		}
		KafkaAddress := cfg.Section("kafka").Key("address").String()
		LogFilePath := cfg.Section("collect").Key("log_file_path").String()
		fmt.Println(KafkaAddress)
		fmt.Println(LogFilePath)*/
	//（2）通过结构体映射的方法
	p := new(AppConfig)
	err := ini.MapTo(p, "./conf/config.ini")
	if err != nil {
		logrus.Error("read config file failed, err: ", err)
		fmt.Println("read config file failed, err: ", err)
		return
	}
	fmt.Printf("p: %#v\n", p)
	// 连接kafka
	err = kafka.Init([]string{p.KafkaConfig.Address})
	if err != nil {
		logrus.Error("init connection kafka failed, err: ", err)
		return
	}
	logrus.Info("init connection kafka success")
	// 初始化tail
}
