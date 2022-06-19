package tailf

import (
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/jokereven/golang-log-collection/logagent/kafka"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tailTack struct {
	path  string
	topic string
	Tails *tail.Tail
}

func NewTailTask(path, topic string) *tailTack {
	tt := &tailTack{
		path:  path,
		topic: topic,
	}
	return tt
}

func (t *tailTack) Init() (err error) {
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.Tails, err = tail.TailFile(t.path, cfg)
	return
}

// 读取日志发往kafka
func (t *tailTack) run() {
	// Tails ->log -> Client -> kafka
	logrus.Infof("path %s is running", t.path)
	for {
		line, ok := <-t.Tails.Lines
		if !ok {
			logrus.Warning("tail file close reopen, path:%s\n", t.path)
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
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)
		// msg -> chan
		kafka.MsgChan(msg)
	}
	return
}
