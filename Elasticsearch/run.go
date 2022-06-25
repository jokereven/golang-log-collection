package main

import (
	"os/exec"
	"time"
)

func run() {
	dataPath := "./es.exe"
	cmd := exec.Command("cmd.exe", "/c", "start "+dataPath)
	err := cmd.Run()
	if err != nil {
		return
	}
}

// Random 生成随机中文名添加到es中
func Random() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker
		run()
	}
}
