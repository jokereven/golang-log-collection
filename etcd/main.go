package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"net"
	"strings"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

func GetLocalCpIp() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:23")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println("local computer ip is:", strings.Split(localAddr.String(), ":")[0])
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func main() {
	// 获取本地ip
	ip, err := GetLocalCpIp()
	if err != nil {
		fmt.Println("get local ip failed,err: ", err)
		return
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			return
		}
	}(cli)

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	var key = fmt.Sprintf("etcd_collect_%s_conf", ip)
	// 向etcd中发送json格式数据 json在线工具（json.cn）
	//str := `[{"path":"./log01.log","topic":"web01_log"},{"path":"./log02.log","topic":"web02_log"}]`
	//str := `[{"path":"./log01.log","topic":"web01_log"},{"path":"./log02.log","topic":"web02_log"},{"path":"./log03.log","topic":"web03_log"}]`
	str := `[{"path":"./log01.log","topic":"web01_log"},{"path":"./log02.log","topic":"web02_log"},{"path":"./log03.log","topic":"web03_log"},{"path":"./log04.log","topic":"web04_log"}]`
	_, err = cli.Put(ctx, key, str)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "etcd_collect_conf")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
