package tailf

import (
	"github.com/jokereven/golang-log-collection/logagent/common"
	"github.com/sirupsen/logrus"
)

// tailTask管理者

type TailTaskMgr struct {
	tailTackMap map[string]*tailTack
	ConfigList  []common.Config      // 所以配置项
	confChan    chan []common.Config // 等待新配置的通道
}

var (
	ttMgr *TailTaskMgr
)

func Init(allConf []common.Config) (err error) {
	// 文件名
	// 配置信息
	// allConf 里面存了若干个日志的收集项
	// 对应没一个日志收集项创建一个对应Tails

	// not := use =
	ttMgr = &TailTaskMgr{
		tailTackMap: make(map[string]*tailTack, 24),
		ConfigList:  allConf,
		confChan:    make(chan []common.Config), // 阻塞的channel
	}
	for _, conf := range allConf {
		tt := NewTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create Tails for path: %s, failed err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		// tt 创建成功就去收集任务
		// 在创建（NewTailTask时创建一个管理者对TailTask进行管理）
		ttMgr.tailTackMap[tt.path] = tt // 登记tailTask, 方便后续管理
		go tt.run()
	}
	go ttMgr.watch() // 等待新配置的到来
	return
}

func (t *TailTaskMgr) watch() {
	for {
		// 初始化confChan
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd, conf:%v, start manage tailTask...", newConf)
		// 新配置来了, 对tailTask 进行管理
		for _, conf := range newConf {
			// 1. 原来有现在也有的不要动
			if t.isExit(conf) {
				continue
			}
			// 2. 原来没有的现在有就创建一个 TailTask
			tt := NewTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create Tails for path: %s, failed err:%v", conf.Path, err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			// tt 创建成功就去收集任务
			// 在创建（NewTailTask时创建一个管理者对TailTask进行管理）
			t.tailTackMap[tt.path] = tt // 登记tailTask, 方便后续管理
			go tt.run()
			// 3. 原来有现在没有就停掉这个 TailTask
		}
	}
}

func (t *TailTaskMgr) isExit(conf common.Config) bool {
	_, ok := t.tailTackMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.Config) {
	ttMgr.confChan <- newConf
}
