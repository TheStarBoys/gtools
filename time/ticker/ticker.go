package ticker

import (
	"time"
)

type Interface interface {
	HasTicker(key string) bool // 是否有ticker
	CancelTicker(key string) // 取消ticker
	AddTicker(key string, duration time.Duration, cycleFunc func()) // 添加ticker
	Reset() // 重置所有ticker
}

type TickerWrap struct {
	cycleFunc func()
	*time.Ticker
}

func newTickerWrap(duration time.Duration, cycleFunc func()) *TickerWrap {
	return &TickerWrap{
		cycleFunc: cycleFunc,
		Ticker: time.NewTicker(duration),
	}
}

func (t *TickerWrap) Run() {
	for {
		select {
		case <- t.Ticker.C:
			t.cycleFunc()
		}
	}
}

type Manager struct {
	tickers map[string]*TickerWrap
}

func NewManager() Interface {
	return &Manager{
		tickers: make(map[string]*TickerWrap),
	}
}

func (manager *Manager) HasTicker(key string) bool {
	_, ok := manager.tickers[key]

	return ok
}

func (manager *Manager) CancelTicker(key string) {
	if manager.HasTicker(key) == false {
		return
	}
	wrap := manager.getTicker(key)
	wrap.Stop()
	delete(manager.tickers, key)
}

func (manager *Manager) AddTicker(key string, duration time.Duration, cycleFunc func()) {
	wrap := newTickerWrap(duration, cycleFunc)
	manager.setTicker(key, wrap)
}

func (manager *Manager) Reset() {
	for key := range manager.tickers {
		manager.CancelTicker(key)
	}
}

func (manager *Manager) getTicker(key string) *TickerWrap {
	return manager.tickers[key]
}

func (manager *Manager) setTicker(key string, wrap *TickerWrap) {
	manager.tickers[key] = wrap
	go wrap.Run()
}