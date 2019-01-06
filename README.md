# drumstick
Implement crond by Golang

---

鼓槌(鸡腿)，golang定时任务包

### 功能特性:

 - 1.提供时间补偿机制，避免周期任务调度中时间损耗累计导致的长期使用后产生任务滞后现象
 - 2.参数可传入自定义函数, 和变长自定义参数
 - 3.可以设定任务的起始时间和任务的执行周期
 - 4.未完待续...

### 方法：

 - NewTask(startTime time.Time, time.Duration, function, ...param) (*Task, error)
   - 创建任务对象,当周期时间小于等于0时返回错误，否则返回nil
   - task, err := drumstick.NewTask(startTime, 2*time.Second, func1, "hello", 1 ,2)
 - (*Task) Start()
   - 启动任务
   - task.Start()
 - (*Task) Stop()
   - 停止任务继续生产，即已经启动的任务不会被结束，而是关闭他的调度器不再生产新任务
   - task.Stop()
 - (*Task) Reset(startTime time.Time, time.Duration)
   - 更新指定任务的周期时间
   - task.Reset(startTime, 1*time.Second)

### 示例:
```Go
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
```

