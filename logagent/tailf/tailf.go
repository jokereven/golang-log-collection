package tailf

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var (
	Tails *tail.Tail
)

func Init(fileName string) (err error) {
	// 文件名
	// 配置信息
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	Tails, err = tail.TailFile(fileName, config)
	if err != nil {
		logrus.Error("tail file for path:%s failed, err:%v\n	", fileName, err)
		fmt.Printf("tail file for path:%s failed, err:%v\n", fileName, err)
		return
	}
	return
}
