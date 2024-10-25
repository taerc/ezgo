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

		// go Work()
		fmt.Println("Done", i)
		i += 1
		time.Sleep(3)
	}

}
