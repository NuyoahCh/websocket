// Package main
// @Author NuyoahCh
// @Date 2025/2/12 23:13
// @Desc
package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

// 聊天室地址
var addr = flag.String("addr", ":8080", "http service address")
// 聊天室
var house sync.Map
var roomMutexes = make(map[string]*sync.Mutex)
var mutexForRoomMutexes = new(sync.Mutex)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	// 路由逻辑适配房间号参数
	r.HandleFunc("/{room}", serveHome)
	r.HandleFunc("/ws/{room}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomId := vars["room"]
		mutexForRoomMutexes.Lock()
		roomMutex, ok := roomMutexes[roomId]
		if ok {
			roomMutex.Lock()
		} else {
			roomMutexes[roomId] = new(sync.Mutex)
			roomMutexes[roomId].Lock()
		}
		mutexForRoomMutexes.Unlock()
		room, ok := house.Load(roomId)
		var hub *Hub
		if ok {
			hub = room.(*Hub)
		} else {
			hub = newHub(roomId)
			house.Store(roomId, hub)
			go hub.run()
		}
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
