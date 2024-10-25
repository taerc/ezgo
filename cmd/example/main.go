package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func Work() {

	for {

		fmt.Println("xinfo ")
		time.Sleep(3)
	}
}

func main_() {

	i := 1
	for {

		// go Work()
		fmt.Println("Done", i)
		i += 1
		time.Sleep(3)
	}

}

func monitorGoroutines(interval time.Duration) {
	var initialGoroutines int
	initialGoroutines = runtime.NumGoroutine()
	fmt.Printf("Initial goroutines: %d\n", initialGoroutines)

	ticker := time.NewTicker(interval)
	for ; true; <-ticker.C {
		currentGoroutines := runtime.NumGoroutine()
		fmt.Printf("Current goroutines: %d\n", currentGoroutines)
		if currentGoroutines-initialGoroutines >= 50 {
			fmt.Println("Warning: High number of goroutines")
		}
	}
}

func createGoroutines(wg *sync.WaitGroup, num int) {
	for i := 0; i < num; i++ {
		go func() {
			// 模拟工作负载
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()
	}
}

func main() {
	// 启动goroutine监控，每隔1秒打印一次当前goroutine数量
	go monitorGoroutines(1 * time.Second)

	var wg sync.WaitGroup
	// 创建500个goroutine
	createGoroutines(&wg, 500)
	wg.Wait() // 等待所有goroutine完成工作
}
