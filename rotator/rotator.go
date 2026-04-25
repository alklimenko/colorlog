package rotator

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/djherbis/times"
)

type Strategy int

const (
	StrategyNone   Strategy = 0
	StrategySize   Strategy = 1
	StrategyPeriod Strategy = 2
)

type Options struct {
	Strategy    Strategy
	Size        int64
	Period      time.Duration
	FilePrefix  string
	Ext         string
	Count       int
	Dirname     string
	CheckPeriod time.Duration
}

type Rotator struct {
	options   Options
	out       io.WriteCloser
	lock      sync.RWMutex
	lastCheck time.Time
	filename  string
}

func (r *Rotator) Write(p []byte) (n int, err error) {
	if r.out == nil {
		_, err = os.Stat(r.filename)
		if os.IsNotExist(err) {
			r.out, err = os.Create(r.filename)
			if err != nil {
				return 0, err
			}
		} else {
			r.out, err = os.OpenFile(r.filename, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return 0, err
			}
		}
	}
	r.check()

	r.lock.Lock()
	defer r.lock.Unlock()
	return r.out.Write(p)
}

func (r *Rotator) check() {
	if r.options.Strategy == StrategyNone || time.Since(r.lastCheck) < r.options.CheckPeriod || r.out == os.Stderr {
		return
	}
	r.lastCheck = time.Now()

	if r.options.Strategy == StrategySize {
		stat, err := os.Stat(r.filename)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
			return
		}
		if stat.Size() > r.options.Size {
			r.shift()
		}
		return
	}
	if r.options.Strategy == StrategyPeriod {
		t, err := times.Stat(r.filename)
		if err != nil {
			_, _ = os.Stderr.Write([]byte(err.Error()))
			return
		}
		if time.Now().Sub(t.BirthTime()) < r.options.Period {
			return
		}
		r.shift()
	}
}

func (r *Rotator) shift() {
	r.lock.Lock()
	defer r.lock.Unlock()
	_ = r.out.Close()
	var err error
	r.filename, err = getNext(r.options.Dirname, r.options.FilePrefix, r.options.Ext, r.options.Count)
	if err != nil {
		r.out = os.Stdout
		_, _ = os.Stderr.Write([]byte(err.Error()))
		return
	}
	r.out, err = os.Create(r.filename)
	if err != nil {
		r.out = os.Stderr
		_, _ = os.Stderr.Write([]byte(err.Error()))
		return
	}
}
