package timer

import (
	"time"
)

type Interface interface {
	IncreaseLeftTime(key string, d time.Duration) // 增加定时器剩余时间
	AddTimer(key string, d time.Duration, callback func()) // 添加定时器
	OnDoTimer(key string) // 定时时间到了，处理定时器的callback
	HasTimer(key string) bool // 是否有定时器
	CancelTimer(key string) bool // 撤消定时器
	Reset() // 重置定时器
}

type TimerWrap struct {
	expirationTime time.Time
	*time.Timer
}

type Manager struct {
	Timers map[string]*TimerWrap
}

func NewManager() Interface {
	return &Manager{
		Timers: make(map[string]*TimerWrap),
	}
}

func (manager *Manager) HasTimer(key string) bool {
	_, ok := manager.Timers[key]

	return ok
}

func (manager *Manager) CancelTimer(key string) bool {
	if !manager.HasTimer(key) {
		return false
	}
	timer := manager.getTimer(key)
	delete(manager.Timers, key)

	return timer.Stop()
}

func (manager *Manager) AddTimer(key string, duration time.Duration, callback func()) {
	if manager.HasTimer(key) {
		manager.CancelTimer(key)
	}

	manager.setTimer(key, duration, func() {
		callback()
		delete(manager.Timers, key)
	})
}

func (manager *Manager) IncreaseLeftTime(key string, d time.Duration) {
	if manager.HasTimer(key) == false {
		return
	}
	wrap := manager.getTimer(key)
	wrap.expirationTime = wrap.expirationTime.Add(d)
	duration := wrap.expirationTime.Sub(time.Now())
	wrap.Timer.Reset(duration)
}

func (manager *Manager) OnDoTimer(key string) {
	if manager.HasTimer(key) == false {
		return
	}
	wrap := manager.getTimer(key)
	wrap.expirationTime = time.Now()
	wrap.Timer.Reset(0)
	delete(manager.Timers, key)
}

func (manager *Manager) Reset() {
	for key := range manager.Timers {
		manager.CancelTimer(key)
	}
}

func (manager *Manager) getTimer(key string) *TimerWrap {
	return manager.Timers[key]
}

func (manager *Manager) setTimer(key string, duration time.Duration, callback func()) {
	timer := time.AfterFunc(duration, callback)
	manager.Timers[key] = &TimerWrap{
		expirationTime: time.Now().Add(duration),
		Timer: timer,
	}
}