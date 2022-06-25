package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

// 连接es
// 将日志数据写入es

type EsClient struct {
	client     *elastic.Client  // es 真正的连接
	index      string           // 数据库
	LogMsgChan chan interface{} // 存放数据的channel
}

var (
	esClient = &EsClient{}
)

func Init(address string, index string, maxsize int64, gn int) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + address))
	if err != nil {
		// Handle error
		panic(err)
	}
	esClient.client = client
	esClient.index = index
	esClient.LogMsgChan = make(chan interface{}, maxsize)

	fmt.Println("connect to es success")

	// 从通道中取数据
	//  将日志文件写入topic
	// 将数据发送到es中
	for i := 0; i < gn; i++ {
		go SendToEs()
	}
	return
}

func SendToEs() {
	for msg := range esClient.LogMsgChan {
		//b, err := json.Marshal(msg)
		//if err != nil {
		//	fmt.Printf("json Marshal failed, err: %v", err)
		//	continue
		//}
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(msg).
			Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

// PutLogData 通过一个首字母大写的函数将包外的msg接受到channel中
func PutLogData(msg interface{}) {
	esClient.LogMsgChan <- msg
}
