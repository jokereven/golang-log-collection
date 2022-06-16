package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client sarama.SyncProducer
)

func Init(address []string) (err error) {
	// 1. 生产者配置
	// config
	config := sarama.NewConfig()
	// 设置 生产者 config
	// ACK应答机制 0|1|all
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 发送到那个分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 确认
	config.Producer.Return.Successes = true

	// 2. 连接kafka
	Client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("producer closed, err:", err)
		fmt.Println("producer closed, err:", err)
		return
	}
	return
}
