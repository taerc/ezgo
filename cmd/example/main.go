package main

import (
	"fmt"
	"runtime"
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

func main() {
	// 启动goroutine监控，每隔1秒打印一次当前goroutine数量
	go monitorGoroutines(1 * time.Second)

	// 创建500个goroutine
	for i := 0; i < 500; i++ {
		go func() {
			for {
				// 模拟工作负载
				time.Sleep(2 * time.Second)
				fmt.Println("T")
			}
		}()
		time.Sleep(1 * time.Second)
	}
}
