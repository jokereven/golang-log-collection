package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
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
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"},
		config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()

	// 3. 封装一个消息
	// 构造⼀个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "gnorev"
	msg.Value = sarama.StringEncoder("you gotta win once in your life")

	// 4. 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}
