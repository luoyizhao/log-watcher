package main

import (
	"log"
	"os"

	watcher "log-watcher/lib"
)

func main() {
	// 要监听的文件路径
	if len(os.Args) < 2 {
		log.Fatal("请输入要监听的文件路径")
	}

	filePathList := os.Args[1:]

	w := watcher.GetWatcher()

	for _, filePath := range filePathList {
		go w.AddWatchFile(filePath)
	}

	// 阻塞主 goroutine
	<-make(chan struct{})
}
