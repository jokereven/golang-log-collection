package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jokereven/golang-log-collection/logagent/common"
	"github.com/jokereven/golang-log-collection/logagent/tailf"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	Client *clientv3.Client
)

func Init(addr []string) (err error) {
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	return
}

func GetConf(key string) (ConfigList []*common.Config, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := Client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key: %s failed, err: %v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get conf len:0 from etcd by key: %s", key)
		return
	}
	fmt.Println(resp.Kvs)
	ret := resp.Kvs[0]
	//ret.Value json格式数据
	err = json.Unmarshal(ret.Value, &ConfigList)
	if err != nil {
		logrus.Error("json Unmarshal failed err: ", err)
		return
	}
	return
}

func WatchConf(key string) {
	// watch key: key change
	rch := Client.Watch(context.Background(), key) // <-chan WatchResponse
	var newConf []common.Config
	for resp := range rch {
		for _, ev := range resp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			err := json.Unmarshal(ev.Kv.Value, &newConf)
			if err != nil {
				logrus.Errorf("json Unmarshal new conf failed, err: %s", err)
				continue
			}
			tailf.SendNewConf(newConf)
		}
	}
}
