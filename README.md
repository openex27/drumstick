# drumstick
Implement crond by Golang

---

鼓槌(鸡腿)，golang定时任务包

###功能特性:

 - 1.提供时间补偿机制，避免周期任务调度中的时间损耗累计，避免长期使用后产生任务滞后现象
 - 2.参数可传入自定义函数
 - 3.未完待续

###方法：

 - NewTask(time.Duration, function, ...param) (*Task, error)
  - 创建任务对象,当周期时间小于等于0时返回错误，否则返回nil
  - task, err := drumstick.NewTask(2*time.Second, func1, "hello", 1 ,2)
 - (*Task) Start()
  - 启动任务
  - task.Start()
 - (*Task) Stop()
  - 停止任务继续生产，即已经启动的任务不会被结束，而是关闭他的调度器不再生产新任务
  - task.Stop()
 - (*Task) Reset(time.Duration)
  - 更新指定任务的周期时间
  - task.Reset(1*time.Second)

###示例:
```
package main

import (
        "fmt"
        "time"

        drum "github.com/openex27/drumstick"
)

func sumEcho(s string, a, b int) {
        fmt.Printf("%s -> %d\n", s, a+b)
}

func main() {
        task, err := drum.NewTask(2*time.Second, sumEcho, "hello", 1, 5)
        if err != nil {
                panic(err)
        }
        task.Start()
        time.Sleep(5 * time.Second)
        task.Reset(1 * time.Second)
        time.Sleep(5 * time.Second)
        task.Stop()
        time.Sleep(1 * time.Second)
}
```

