package time

import (
	"sync"
	itime "time"
)

const (
	timerFormat      = "2006-01-02 15:04:05"
	infiniteDuration = itime.Duration(1<<63 - 1)
)

type TimerData struct {
	Key    string
	expire itime.Time
	fn     func()
	index  int
	next   *TimerData
}

func (td *TimerData) Delay() itime.Duration {
	return itime.Until(td.expire)
}

func (td *TimerData) ExpireString() string {
	return td.expire.Format(timerFormat)
}

type Timer struct {
	lock   sync.Mutex
	free   *TimerData
	timers []*TimerData
	signal *itime.Timer
	num    int
}

func NewTimer(num int) (t *Timer) {
	t = new(Timer)
	t.init(num)
	return t
}

func (t *Timer) Init(num int) {
	t.init(num)
}

func (t *Timer) init(num int) {
	t.signal = itime.NewTimer(infiniteDuration)
	t.timers = make([]*TimerData, 0, num)
	t.num = num
	t.grow()
	go t.start()
}

func (t *Timer) grow() {
	var (
		i   int
		td  *TimerData
		tds = make([]TimerData, t.num)
	)
	t.free = &(tds[0])
	td = t.free
	for i = 1; i < t.num; i++ {
		td.next = &(tds[i])
		td = td.next
	}
	td.next = nil
}

func (t *Timer) get() (td *TimerData) {
	if td = t.free; td == nil {
		t.grow()
		td = t.free
	}
	t.free = td.next
	return
}

func (t *Timer) put(td *TimerData) {
	td.fn = nil
	td.next = t.free
	t.free = td
}

func (t *Timer) Add(expire itime.Duration, fn func()) (td *TimerData) {
	t.lock.Lock()
	td = t.get()
	td.expire = itime.Now().Add(expire)
	td.fn = fn
	t.add(td)
	t.lock.Unlock()
	return
}

func (t *Timer) Del(td *TimerData) {
	t.lock.Lock()
	t.del(td)
	t.put(td)
	t.lock.Unlock()
}

func (t *Timer) add(td *TimerData) {
	var d itime.Duration
	td.index = len(t.timers)
	t.timers = append(t.timers, td)
	t.up(td.index)
	if td.index == 0 {
		d = td.Delay()
		t.signal.Reset(d)
	}
}

func (t *Timer) del(td *TimerData) {
	var (
		i    = td.index
		last = len(t.timers) - 1
	)
	if i < 0 || i > last || t.timers[i] != td {
		return
	}
	if i != last {
		t.swap(i, last)
		t.down(i, last)
		t.up(i)
	}
	t.timers[last].index = -1
	t.timers = t.timers[:last]
}

func (t *Timer) Set(td *TimerData, expire itime.Duration) {
	t.lock.Lock()
	t.del(td)
	td.expire = itime.Now().Add(expire)
	t.add(td)
	t.lock.Unlock()
}

func (t *Timer) start() {
	for {
		t.expire()
		<-t.signal.C
	}
}

func (t *Timer) expire() {
	var (
		fn func()
		td *TimerData
		d  itime.Duration
	)
	t.lock.Lock()
	for {
		if len(t.timers) == 0 {
			d = infiniteDuration
			break
		}
		td = t.timers[0]
		if d = td.Delay(); d > 0 {
			break
		}
		fn = td.fn
		// let caller put back
		t.del(td)
		t.lock.Unlock()
		if fn != nil {
			fn()
		}
		t.lock.Lock()
	}
	t.signal.Reset(d)

	t.lock.Unlock()
}

func (t *Timer) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i >= j || !t.less(j, i) {
			break
		}
		t.swap(i, j)
		j = i
	}
}

func (t *Timer) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 {
			break
		}
		j := j1
		if j2 := j1 + 1; j2 < n && !t.less(j1, j2) {
			j = j2
		}
		if !t.less(j, i) {
			break
		}
		t.swap(i, j)
		i = j
	}
}

func (t *Timer) less(i, j int) bool {
	return t.timers[i].expire.Before(t.timers[j].expire)
}

func (t *Timer) swap(i, j int) {
	t.timers[i], t.timers[j] = t.timers[j], t.timers[i]
	t.timers[i].index = i
	t.timers[j].index = j
}
