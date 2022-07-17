package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {
	group, ctx := errgroup.WithContext(context.Background())
	signalchan := GetSignalChan()
	svrList := InitMockServer()

	for _, svr := range svrList {
		group.Go(GetStartFunc(svr))
	}

	mocksignalchan := make(chan string, 1)
	go Mocksignalerror(mocksignalchan)

	var err error
	select {
	case <-signalchan:
		err = errors.New("signal error")
	case <-ctx.Done():
		err = ctx.Err()
	case <-mocksignalchan:
		err = errors.New("mock signal receive")
	}
	fmt.Printf("Close reason:%s\n", err)
	for _, svr := range svrList {
		group.Go(GetShutDownFunc(svr, ctx))
	}
	group.Go(func() error {
		time.Sleep(time.Second * 60) // 60秒还关不完就不管了，强制关闭
		return errors.New("ShutDown time out")
	})

	err = group.Wait()
	if err != nil {
		fmt.Printf("ShutDown Reason:%s\n", err)
	} else {
		fmt.Printf("All Server ShutDown\n")
	}
	fmt.Printf("ShutDown End\n")
}

func GetStartFunc(svr *http.Server) func() error {
	return func() error {
		err := StartHttpServer(svr)
		if err != nil {
			fmt.Printf("http server error,Port %s,error: %s\n", svr.Addr, err)
		}
		return err
	}
}

func GetShutDownFunc(svr *http.Server, ctx context.Context) func() error {
	return func() error {
		//模拟关闭超时
		//time.Sleep(time.Second * 80)
		err := svr.Shutdown(ctx)
		if err != nil {
			fmt.Printf("ShutDown ,Port %s, error: %s\n", svr.Addr, err)
		} else {
			fmt.Printf("ShutDown success,Port %s\n", svr.Addr)
		}
		return nil //关闭服务时不希望一个服务关闭失败时直接全部退出，尽量等待所有服务关闭。
	}
}

//模拟错误发生
func Mocksignalerror(mocksignalchan chan string) {
	time.Sleep(time.Second * 10)
	fmt.Printf("mock sign send\n")
	mocksignalchan <- "mock sign send"
}

//启动服务
func StartHttpServer(svr *http.Server) error {
	fmt.Printf("http server start,Port %s\n", svr.Addr)
	err := svr.ListenAndServe()
	return err
}

//linux signal信号注册
func GetSignalChan() chan os.Signal {
	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan)
	return signalchan
}

//模拟两个不同的服务
func InitMockServer() []*http.Server {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", home1)
	mux1.HandleFunc("/hello", helloServer1)

	mux2 := http.NewServeMux()
	mux2.HandleFunc("/", home2)
	mux2.HandleFunc("/hello", helloServer2)

	svr1 := &http.Server{Addr: ":8080", Handler: mux1}
	svr2 := &http.Server{Addr: ":8081", Handler: mux2}

	svrList := []*http.Server{svr1, svr2}
	return svrList
}

func helloServer1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world 1")
}

func home1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home page 1")
}

func helloServer2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world 2")
}

func home2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home page 2")
}
