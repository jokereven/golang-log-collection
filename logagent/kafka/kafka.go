package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

func Init(address []string, ChanSize int64) (err error) {
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
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("producer closed, err:", err)
		fmt.Println("producer closed, err:", err)
		return
	}
	// 3. 初始化msgChan
	msgChan = make(chan *sarama.ProducerMessage, ChanSize)
	// 起一个后台的goroutine从msgChan中读取数据
	go SengMsg()
	return
}

// SengMsg 从msgChan中读取msg发送到kafka
func SengMsg() {
	for {
		select {
		case msg := <-msgChan:
			// 4. 发送消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed, err:", err)
				return
			}
			logrus.Infof("pid:%v offset:%v\n", pid, offset)
		}
	}
}

func MsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
