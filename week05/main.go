package main

import (
	"fmt"
	"math/rand"
	"time"
)

//1. 参考 Hystrix 实现一个滑动窗口计数器。
func main() {
	var windowSize time.Duration = 5000 * time.Millisecond //窗口大小5000毫秒
	var windowCount int = 4                                //窗口内允许最多调用4次
	Work := RateLimiterSlidingWindow(windowSize, windowCount, MockWorkFunc)
	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			Work() //模拟随机的调用 预计平均5秒调用10次，成功4次。
		}
	}()
	select {}
}

func RateLimiterSlidingWindow(windowSize time.Duration, windowCount int, workFunc func()) func() {
	windowCache := make([]time.Time, 0)
	startindex := 0
	return func() {
		for i := startindex; i < len(windowCache); i++ {
			if time.Now().Sub(windowCache[i]) < windowSize {
				startindex = i
				break
			}
		}
		fmt.Printf("startindex-->%d\n", startindex)
		if len(windowCache)-startindex < windowCount {
			workFunc()
			windowCache = append(windowCache, time.Now())
		} else {
			timeTemplate := "2006-01-02 15:04:05"
			fmt.Printf("Call WorkFunc Fail:%s\n", time.Now().Format(timeTemplate))
		}
	}
}

func MockWorkFunc() {
	timeTemplate := "2006-01-02 15:04:05"
	fmt.Printf("Call WorkFunc Success:%s\n", time.Now().Format(timeTemplate))
}
