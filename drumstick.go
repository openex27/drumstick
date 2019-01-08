package drumstick

/*
1.提供时间补偿机制，避免周期任务调度中时间损耗累计导致的长期使用后产生任务滞后现象
2.参数可传入自定义函数, 和变长自定义参数
3.未完待续...
*/
import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type funcer struct {
	fn   reflect.Value
	args []reflect.Value
}

type Task struct {
	funcer
	Quit    chan struct{}
	pChange bool
	sync.Mutex
	setStartTime time.Time     //用户设定的起始时间
	startTime    time.Time     //本次启动程序后任务首次调度的起始时间
	period       time.Duration //执行周期
	count        int64         //程序启动后的执行次数计数
}

/*
通过用户设定的起始时间和执行周期计算当前周期的起始时间点
如果用户设定的时间大于当前时间，则返回设定的起始时间
*/
func prevTime(setstarttime time.Time, period time.Duration) (time.Time, error) {
	if period <= 0 {
		return setstarttime, errors.New("period is 0,it must be greater than 0")
	}
	tmp := setstarttime //临时时间点
	if tmp.Add(period).Before(time.Now()) {
		dur := time.Now().Sub(tmp)
		cnt := dur.Nanoseconds() / period.Nanoseconds()
		tmp = setstarttime.Add(time.Duration(cnt * period.Nanoseconds()))
	}
	return tmp, nil
}

func (t *Task) nextTime(doing <-chan struct{}) time.Duration {
	<-doing
	t.Lock()
	defer t.Unlock()
	nowTime := time.Now()
	if t.pChange == true {
		t.pChange = false
		t.startTime, _ = prevTime(t.setStartTime, t.period)
		t.count = 1
		return t.period
	}
	calResult := t.startTime.Add(time.Duration((t.count + 1) * int64(t.period))).Sub(nowTime)
	t.count++
	return calResult
}

/*
(*Task) Start()
	启动任务
	task.Start()
*/
func (t *Task) Start() {
	doing := make(chan struct{})
	go func() {
		doing <- struct{}{}
		t.fn.Call(t.args)

	}()
	<-doing
	t.startTime, _ = prevTime(t.setStartTime, t.period)
	newTimeValue := t.startTime.Add(t.period).Sub(time.Now())
	t.count = 1
	go func() {
		for {
			select {
			case <-t.Quit:
				return
			case <-time.After(newTimeValue):
				go func() {
					doing <- struct{}{}
					t.fn.Call(t.args)
				}()
				newTimeValue = t.nextTime(doing)

			}
		}
	}()
}

/*
(*Task) Stop()
	停止任务继续生产，即已经启动的任务不会被结束，而是关闭他的调度器不再生产新任务
	task.Stop()
*/
func (t *Task) Stop() {
	t.Quit <- struct{}{}
}

/*
(*Task) Reset(time.Duration)
更新指定任务的周期时间
task.Reset(1*time.Second)
*/
func (t *Task) Reset(newStartTime time.Time, newPeriod time.Duration) {
	t.Lock()
	defer t.Unlock()
	t.setStartTime = newStartTime
	t.period = newPeriod
	t.pChange = true
}

/*
NewTask(startTime,time.Duration, function, ...param) (*Task, error)
创建任务对象,当周期时间小于等于0时返回错误，否则返回nil
task, err := drumstick.NewTask(startTime,2*time.Second, func1, "hello", 1 ,2)
*/
func NewTask(startTime time.Time, period time.Duration, f interface{}, args ...interface{}) (*Task, error) {
	if period <= 0 {
		return nil, errors.New("period is 0,it will crazy running")
	}
	newTask := new(Task)
	if reflect.TypeOf(f).Kind() != reflect.Func {
		return nil, errors.New("void interface delivered")
	}
	newTask.fn = reflect.ValueOf(f)
	newTask.Quit = make(chan struct{}, 1)
	newTask.setStartTime = startTime
	newTask.period = period
	tempA := []reflect.Value{}
	for _, temp := range args {
		tempA = append(tempA, reflect.ValueOf(temp))
	}
	newTask.args = tempA
	return newTask, nil
}
