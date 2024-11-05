package watcher

import (
	"bufio"
	"io"
	"log"
	"log-watcher/service"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fileMap map[*fsnotify.Watcher]*os.File
}

func GetWatcher() *Watcher {
	return &Watcher{
		fileMap: make(map[*fsnotify.Watcher]*os.File),
	}
}

func (w *Watcher) AddWatchFile(filePath string) {
	wg := &sync.WaitGroup{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("打开文件失败:", err)
		return
	}
	defer file.Close()
	file.Seek(0, io.SeekEnd)

	// 监听文件
	err = watcher.Add(filePath)
	w.fileMap[watcher] = file
	if err != nil {
		log.Fatal(err)
	}

	log.Println("开始监听文件:", filePath)

	wg.Add(1)
	// 监听事件
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 当文件被写入时，打印日志
				if event.Op&fsnotify.Write == fsnotify.Write {
					readingFile := w.fileMap[watcher]
					if readingFile == nil {
						log.Fatal("找不到文件:", event.Name)
					}
					// 这里可以添加逻辑来读取文件内容
					w.readFileContent(event, readingFile)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("错误:", err)
			}
		}
	}()
	wg.Wait()
}

// 读取文件内容
func (w *Watcher) readFileContent(event fsnotify.Event, file *os.File) {

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		service.Callback(event, scanner.Text())
	}

	// 处理错误
	if err := scanner.Err(); err != nil {
		log.Println("读取文件失败:", err)
	}
}
