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

func main() {
	ticker := time.Tick(1 * time.Second)
	for {
		<-ticker
		run()
	}
}
