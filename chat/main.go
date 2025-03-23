// Package chat
// @Author NuyoahCh
// @Date 2025/2/12 23:13
// @Desc
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

// HTTP 服务，把 html 文件返回给请求方
func serveHome(w http.ResponseWriter, r *http.Request) {
	// 打印 URL 路径
	log.Println(r.URL)
	// check the URL
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 检查完毕跳转到这个静态页面
	http.ServeFile(w, r, "index.html")
}

func main() {
	flag.Parse()
	// 调用 hub
	hub := newHub()
	// # command-line-arguments
	//./main.go:39:9: undefined: newHub
	//./main.go:43:3: undefined: serveWs
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
