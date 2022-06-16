package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jokereven/golang-log-collection/logagent/kafka"
	"github.com/jokereven/golang-log-collection/logagent/tailf"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"time"
)

// 日志收集客户端
// 类似的开源项目filebeat
// 将指定Path下的日志收集, 发送到kafka

type AppConfig struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address     string `ini:"address"`
	Topic       string `ini:"topic"`
	MsgChanSize int64  `ini:"msg_chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"log_file_path"`
}

// run 真正的业务逻辑
func run() (err error) {
	// Tails ->log -> Client -> kafka
	for {
		line, ok := <-tailf.Tails.Lines
		if !ok {
			logrus.Warning("tail file close reopen, filename:%s\n", tailf.Tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		// ? 具体业务逻辑, 将消息发送到kafka
		// 利用通道将同步的代码改成异步
		// 将每一行的消息发送到kafka
		msg := &sarama.ProducerMessage{}
		msg.Topic = "gnorev"
		msg.Value = sarama.StringEncoder(line.Text)
		// msg -> chan
		kafka.MsgChan <- msg
	}
	return
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
	// 1. 初始化连接kafka(做好准备工作)
	err = kafka.Init([]string{p.KafkaConfig.Address}, p.KafkaConfig.MsgChanSize)
	if err != nil {
		logrus.Error("init connection kafka failed, err: ", err)
		return
	}
	logrus.Info("init connection kafka success")

	// 2. 根据配置文件中的日志初始化tail
	err = tailf.Init(p.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("init tailf failed, err: ", err)
		return
	}
	logrus.Info("init tailf success")

	// 3. 把日志通过sarama发送包kafka
	err = run()
	if err != nil {
		logrus.Error("run failed err:", err)
		return
	}
}
