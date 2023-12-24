package main

import (
	"fmt"
	"time"

	"github.com/binc4t/playground/wpool/wpool"
)

func Print(str string) func() {
	return func() {
		fmt.Println(str)
	}
}

func main() {
	wpool.Run(Print("task 1"))
	wpool.Run(Print("task 2"))
	time.Sleep(time.Second * 3)
}
