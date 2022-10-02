package cacheServer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {

	RedisClient := NewClient("localhost:8080", "localhost:6379", "", 0)
	MessageCache1 := NewCache("MessageCache", time.Second*5, "http://localhost:8080", RedisClient, DialClearAfterGet(true), DialSyncFromOtherCache(true))
	MessageCache2 := NewCache("MessageCache", time.Second*5, "http://localhost:8081", RedisClient, DialClearAfterGet(true), DialSyncFromOtherCache(true))
	go Start(":8080", MessageCache1)
	go Start(":8081", MessageCache2)
	go MockWebSocketMessage(MessageCache1)
	select {}
}

func MockWebSocketMessage(cache CacheServer) {
	rand.Seed(time.Now().UnixNano())
	var i int = 0
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		var receiver string
		if i%2 == 0 {
			receiver = "HYC"
		} else {
			receiver = "ZMN"
		}
		message := "send message " + strconv.Itoa(i) + " times"
		if err := cache.Set(receiver, message); err != nil {
			fmt.Printf("cache Set Error:%s \n", err)
		}
		i = i + 1
	}
}

func Start(Port string, cache CacheServer) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", health)
	mux.HandleFunc("/syncCache", syncCache(cache))
	mux.HandleFunc("/getCache", getCache(cache))
	svr := &http.Server{Addr: Port, Handler: mux}
	err := svr.ListenAndServe()
	return err
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server work")
}

func getCache(cache CacheServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := GetUrlArg(r, "key")
		value, err := cache.Get(key)
		if err != nil {
			fmt.Printf("cache Get Error:%s \n", err)
		}
		fmt.Fprintf(w, "Cache Key:%s,%v \n", key, value)
	}
}

func syncCache(cache CacheServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := cache.GetAll()
		if err != nil {
			fmt.Printf("cache GetAll Error:%s \n", err)
		}
		res, err := json.Marshal(value)
		if err != nil {
			fmt.Printf("cache json Error:%s \n", err)
		}
		w.Write(res)
	}
}

func GetUrlArg(r *http.Request, name string) string {
	var arg string
	values := r.URL.Query()
	arg = values.Get(name)
	return arg
}
