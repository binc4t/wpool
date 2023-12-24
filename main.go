package main

import (
	"fmt"
	"time"

	"github.com/binc4t/playground/wpool/wpool"
)

func Sleep(name string, t time.Duration) func() {
	return func() {
		time.Sleep(t)
		fmt.Printf("%s sleep %s done\n", name, t)
	}
}

func main() {
	wpool.Run(Sleep("task 1", time.Second))
	fmt.Println("add task 1, now worker count: ", wpool.DefaultPool.WorkerCount())

	wpool.Run(Sleep("task 2", time.Second))
	fmt.Println("add task 2, now worker count: ", wpool.DefaultPool.WorkerCount())

	time.Sleep(time.Second * 3)
}
