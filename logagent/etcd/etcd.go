package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	Client *clientv3.Client
)

type Config struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

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

func GetConf(key string) (ConfigList []*Config, err error) {
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
