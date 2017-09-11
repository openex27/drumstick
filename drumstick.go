package drumstick

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
	period    time.Duration
	startTime time.Time
	count     int64
}

func (t *Task) nextTime() time.Duration {
	t.Lock()
	defer t.Unlock()
	nowTime := time.Now()
	if t.pChange == true {
		t.pChange = false
		t.startTime = nowTime
		t.count = 1
		return t.period
	}
	calResult := t.startTime.Add(time.Duration((t.count + 1) * int64(t.period))).Sub(nowTime)
	t.count++
	return calResult
}

func (t *Task) Start() {
	go t.fn.Call(t.args)
	newTimeValue := t.period
	t.startTime = time.Now()
	t.count = 1
	go func() {
		for {
			select {
			case <-t.Quit:
				return
			case <-time.After(newTimeValue):
				go t.fn.Call(t.args)
				newTimeValue = t.nextTime()
			}
		}
	}()
}

func (t *Task) Stop() {
	t.Quit <- struct{}{}
}

func (t *Task) Reset(newPeriod time.Duration) {
	t.Lock()
	defer t.Unlock()
	t.period = newPeriod
	t.pChange = true
}

func NewTask(period time.Duration, f interface{}, args ...interface{}) (*Task, error) {
	if period <= 0 {
		return nil, errors.New("period is 0,it will crazy running")
	}
	newTask := new(Task)
	newTask.fn = reflect.ValueOf(f)
	newTask.Quit = make(chan struct{}, 1)
	newTask.period = period
	tempA := []reflect.Value{}
	for _, temp := range args {
		tempA = append(tempA, reflect.ValueOf(temp))
	}
	newTask.args = tempA
	return newTask, nil
}
