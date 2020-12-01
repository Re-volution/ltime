//修改版，不用锁，改为原子操作
package ltime

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type timeT2 struct {
	ms   *int64
	now  *unsafe.Pointer
	s    *int64
	last int64
	nano *int64
}

var gtime *timeT2

var bo = false

func init() {
	gtime = new(timeT2)
	gtime.ms = new(int64)
	gtime.s = new(int64)
	gtime.last = 0
	gtime.nano = new(int64)
	gtime.now = new(unsafe.Pointer)

	go gtime.update()
}

func newpointer(data *time.Time) unsafe.Pointer {
	return unsafe.Pointer(data)
}

func (t *timeT2) update() {
	fmt.Println("启动时间函数2")
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("时间函数2出错,捕捉到错误:", err)
			bo = false
			go t.update()
		}
	}()
	var on = new(sync.Once)

	for {
		n := time.Now()
		nt := n.UnixNano() / 1e6
		if nt >= t.last+20 {
			t.last = nt
			atomic.StoreInt64(t.nano, n.UnixNano())
			atomic.StoreInt64(t.ms, nt)
			atomic.StoreInt64(t.s, nt/1e3)
			atomic.StorePointer(t.now, newpointer(&n))

		} else {
			time.Sleep(time.Millisecond * 10)
		}
		on.Do(func() {
			bo = true
		})
	}
}

//获取秒数
func GetS() int64 {
	if bo {
		return atomic.LoadInt64(gtime.s)
	} else {
		return time.Now().Unix()
	}

}

//获取毫秒数
func GetMS() int64 {
	if bo {
		return atomic.LoadInt64(gtime.ms)
	} else {
		return time.Now().UnixNano() / 1e6
	}
}

//获取时间戳
func GetNow() time.Time {
	if bo {
		return *(*time.Time)(atomic.LoadPointer(gtime.now))
	} else {
		return time.Now()
	}
}

//获得纳秒
func GetNano() int64 {
	if bo {
		return atomic.LoadInt64(gtime.nano)
	} else {
		return time.Now().UnixNano()
	}
}
