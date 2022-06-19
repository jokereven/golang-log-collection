package main

import (
	"fmt"
	"github.com/jokereven/golang-log-collection/logagent/etcd"
	"github.com/jokereven/golang-log-collection/logagent/kafka"
	"github.com/jokereven/golang-log-collection/logagent/tailf"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// 日志收集客户端
// 类似的开源项目filebeat
// 将指定Path下的日志收集, 发送到kafka

type AppConfig struct {
	KafkaConfig   `ini:"kafka"`
	EtcdConfig    `ini:"etcd"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address     string `ini:"address"`
	MsgChanSize int64  `ini:"msg_chan_size"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

type CollectConfig struct {
	LogFilePath string `ini:"log_file_path"`
}

// run 真正的业务逻辑
/*func run() (err error) {
	// Tails ->log -> Client -> kafka
	for {
		line, ok := <-tailf.Tails.Lines
		if !ok {
			logrus.Warning("tail file close reopen, filename:%s\n", tailf.Tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		// 当日志文件为空行就不将消息发送到kafka
		if len(strings.Trim(line.Text, "\r")) == 0 {
			logrus.Info("出现空行直接跳过...")
			continue
		}
		// ? 具体业务逻辑, 将消息发送到kafka
		// 利用通道将同步的代码改成异步
		// 将每一行的消息发送到kafka
		msg := &sarama.ProducerMessage{}
		msg.Topic = "gnorev"
		msg.Value = sarama.StringEncoder(line.Text)
		// msg -> chan
		kafka.MsgChan(msg)
	}
	return
}*/

// 防止进程退出
func run() {
	select {}
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

	// 通过etcd加载日志收集配置信息
	// 初始化etcd连接
	// 从etcd拉取日志收集配置项
	err = etcd.Init([]string{p.EtcdConfig.Address})
	if err != nil {
		logrus.Error("init etcd failed, err: ", err)
		return
	}

	allConf, err := etcd.GetConf(p.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed, err:%v ", err)
		return
	}
	fmt.Println(allConf)

	// 监控 etcd p.EtcdConfig.CollectKey 配置项的变化
	go etcd.WatchConf(p.EtcdConfig.CollectKey)

	// 2. 根据配置文件中的日志初始化tail
	//err = tailf.Init(p.CollectConfig.LogFilePath)
	err = tailf.Init(allConf) // 把从etcd中读取到的配置传到Init中
	if err != nil {
		logrus.Error("init tailf failed, err: ", err)
		return
	}
	logrus.Info("init tailf success")

	// 3. 把日志通过sarama发送包kafka
	/*	err = run()
		if err != nil {
			logrus.Error("run failed err:", err)
			return
		}*/
	run()
}
