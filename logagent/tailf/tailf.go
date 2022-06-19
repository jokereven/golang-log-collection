package tailf

import (
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/jokereven/golang-log-collection/logagent/common"
	"github.com/jokereven/golang-log-collection/logagent/kafka"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	confChan chan []common.Config
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

func Init(allConf []*common.Config) (err error) {
	// 文件名
	// 配置信息
	// allConf 里面存了若干个日志的收集项
	// 对应没一个日志收集项创建一个对应Tails

	for _, conf := range allConf {
		tt := NewTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create Tails for path: %s, failed err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		// tt 创建成功就去收集任务
		go tt.run()
	}
	// 初始化confChan
	confChan = make(chan []common.Config) // 阻塞的channel
	newConf := <-confChan
	logrus.Infof("get new conf from etcd: %v", newConf)
	// 等待新配置的到来
	return
}

func SendNewConf(newConf []common.Config) {
	confChan <- newConf
}
