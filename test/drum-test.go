package main

import (
	"fmt"
	"time"

	drum "github.com/bkzy-wangjp/drumstick"
)

func sumEcho(s string, a, b int) {
	fmt.Printf("%s -> %d,now=%s\n", s, a+b, time.Now())
}

func main() {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-01-06 21:10:00", time.Local)
	task, err := drum.NewTask(startTime, 5*time.Second, sumEcho, "hello", 1, 5)
	if err != nil {
		panic(err)
	}
	fmt.Println("Start:")
	task.Start()
	time.Sleep(30 * time.Second)

	fmt.Println("Reset:")
	task.Reset(startTime, 4*time.Second)
	time.Sleep(30 * time.Second)

	fmt.Println("Stop.")
	task.Stop()
	time.Sleep(5 * time.Second)
}
