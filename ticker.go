package tickers

import (
	"time"
)

// ResetTicker is custom ticker type with reset available
type ResetTicker struct {
	// properties
	C        chan int
	Interval time.Duration
	Enabled  bool

	// fields
	current int
	reset   chan interface{}
}

// NewResetTicker - creates new ResetTicker
func NewResetTicker(interval time.Duration, start bool) (t *ResetTicker) {
	t = &ResetTicker{
		C:        make(chan int),
		Interval: interval,
		Enabled:  start,

		current: 0,
		reset:   make(chan interface{}),
	}
	if start {
		t.start()
	}
	return
}

// Start starts or continues ticker
func (t *ResetTicker) Start() {
	if t.Enabled {
		return
	}
	t.Enabled = true
	t.start()
}

// Reset is reset for ticker, nxt event will be fired after ticker interval
func (t *ResetTicker) Reset() {
	t.reset <- 0
}

// Stop current ticker
func (t *ResetTicker) Stop() {
	if !t.Enabled {
		return
	}
	t.Enabled = false
	t.current = 0
}

func (t *ResetTicker) start() {
	go func() {
		for t.Enabled {
			select {
			case <-time.After(t.Interval):
				t.current++
				t.C <- t.current
				continue
			case <-t.reset:
				t.current = 0
				continue
			}
		}
	}()
}
