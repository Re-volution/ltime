//暂时弃用锁
package ltime

import (
	"fmt"
	"sync"
	"time"
)

var tGtime = new(timeT)

type timeT struct {
	l    sync.RWMutex
	ms   int64
	now  time.Time
	s    int64
	last int64
	nano int64
}

func init2() {
	tGtime.l.Lock()
	tGtime.now = time.Now()
	tGtime.ms = tGtime.now.UnixNano() / 1e6
	tGtime.s = tGtime.now.Unix()
	tGtime.last = tGtime.ms
	tGtime.nano = tGtime.now.UnixNano()
	tGtime.l.Unlock()
	go tGtime.update()
}

func (t *timeT) update() {
	fmt.Println("启动时间函数")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("时间函数出错,捕捉到错误:", err)
			go t.update()

		}
	}()
	for {
		n := time.Now()
		na := n.UnixNano()
		nt := na / 1e6
		if nt >= t.last+20 {
			t.l.Lock()
			t.nano = na
			t.last = nt
			t.ms = nt
			t.s = nt / 1e3
			t.now = n

			t.l.Unlock()
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
}

//获取秒数
func (t *timeT) GetS() int64 {
	t.l.RLock()
	nt := t.s
	t.l.RUnlock()
	return nt
}

//获取毫秒数
func (t *timeT) GetMS() int64 {
	t.l.RLock()
	nt := t.ms
	t.l.RUnlock()
	return nt
}

//获取时间戳
func (t *timeT) GetNow() time.Time {
	t.l.RLock()
	nt := t.now
	t.l.RUnlock()
	return nt
}

//获得纳秒
func (t *timeT) GetNano() int64 {
	t.l.RLock()
	nt := t.nano
	t.l.RUnlock()
	return nt
}
