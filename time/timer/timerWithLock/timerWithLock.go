package timerWithLock

import (
	"time"
	"sync"
	"github.com/TheStarBoys/gtools/time/timer"
)

// 加锁版

type TimerWrap struct {
	expirationTime time.Time
	*time.Timer
}

type Manager struct {
	Timers map[string]*TimerWrap
	sync.RWMutex
}

func NewManager() timer.Interface {
	return &Manager{
		Timers: make(map[string]*TimerWrap),
	}
}

func (manager *Manager) HasTimer(key string) bool {
	manager.RLock()
	_, ok := manager.Timers[key]
	manager.RUnlock()
	return ok
}

func (manager *Manager) CancelTimer(key string) bool {
	if !manager.HasTimer(key) {
		return false
	}
	timer := manager.getTimer(key)
	ok := timer.Stop()
	if ok {
		manager.Lock()
		delete(manager.Timers, key)
		manager.Unlock()
	}

	return ok
}

func (manager *Manager) AddTimer(key string, duration time.Duration, callback func()) {
	if manager.HasTimer(key) {
		manager.CancelTimer(key)
	}

	manager.setTimer(key, duration, func() {
		if manager.HasTimer(key) == false {
			return
		}
		manager.Lock()
		delete(manager.Timers, key)
		manager.Unlock()
		callback()
	})
}

func (manager *Manager) IncreaseLeftTime(key string, d time.Duration) {
	if manager.HasTimer(key) == false {
		return
	}
	wrap := manager.getTimer(key)
	manager.Lock()
	defer manager.Unlock()
	wrap.expirationTime = wrap.expirationTime.Add(d)
	duration := wrap.expirationTime.Sub(time.Now())
	wrap.Timer.Reset(duration)
}

func (manager *Manager) OnDoTimer(key string) {
	if manager.HasTimer(key) == false {
		return
	}
	wrap := manager.getTimer(key)
	manager.Lock()
	defer manager.Unlock()
	wrap.expirationTime = time.Now()
	wrap.Timer.Reset(0)
	delete(manager.Timers, key)
}

func (manager *Manager) Reset() {
	for key := range manager.Timers {
		manager.CancelTimer(key)
	}
}

func (manager *Manager) CycleCallFunc(key string, duration time.Duration, f func()) {
	var callback func()
	callback = func() {
		f()
		manager.AddTimer(key, duration, callback)
	}
	callback()
}

func (manager *Manager) getTimer(key string) *TimerWrap {
	manager.RLock()
	defer manager.RUnlock()
	return manager.Timers[key]
}

func (manager *Manager) setTimer(key string, duration time.Duration, callback func()) {
	timer := time.AfterFunc(duration, callback)
	manager.Lock()
	defer manager.Unlock()
	manager.Timers[key] = &TimerWrap{
		expirationTime: time.Now().Add(duration),
		Timer: timer,
	}
}
