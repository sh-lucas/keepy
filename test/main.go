package main

import (
	"fmt"
	"time"
)

func main() {
	count := 1
	for {
		fmt.Printf("Hello by the %dº time!\n", count)
		count++
		time.Sleep(2 * time.Second)
	}
}
