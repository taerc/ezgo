package main

import (
	"fmt"
	"time"
)

func Work() {

	for {

		fmt.Println("xinfo ")
		time.Sleep(3)
	}
}

func main() {

	i := 1
	for {

		go Work()
		fmt.Println("D", i)
		i += 1
		time.Sleep(9)
	}

}
