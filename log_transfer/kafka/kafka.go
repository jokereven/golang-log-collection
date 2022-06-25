package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jokereven/log_transfer/es"
)

// 初始化kafka连接
// 从kafka中取出日志数据

func Init(address []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(address, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				//fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s", msg.Partition, msg.Offset, msg.Key, msg.Value)
				// 同步 -> 异步 将取出来的数据放到channel中
				//LogMsgChan <- msg
				fmt.Println(msg.Topic, string(msg.Value))
				var m1 map[string]interface{}
				err := json.Unmarshal(msg.Value, &m1)
				if err != nil {
					fmt.Printf("json Unmarshal failed, err: %v", err)
					continue
				}
				// 将m1发送到channel中
				es.PutLogData(m1)
			}
		}(pc)
	}
	return
}
