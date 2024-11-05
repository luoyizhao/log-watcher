package service

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func Callback(event fsnotify.Event, content string) {
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	if content != "" {
		log.Println("content:", content)
		// do some thing
		// post(`{"msg_type": "text", "content": {"text": "` + content + `"}}`)
	}
}

func post(data string) {
	url := ""
	jsonData := []byte(data)

	// 发送 POST 请求
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("发送请求失败:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return
	}

	log.Println("响应状态:", response.Status)
	log.Println("响应内容:", string(body))
}
