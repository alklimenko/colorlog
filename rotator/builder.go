package rotator

import (
	"os"
	"time"

	"github.com/djherbis/times"
)

const (
	defaultPrefix      = "log_"
	defaultExt         = ".log"
	defaultCheckPeriod = time.Second
)

type builder struct {
	r Rotator
}

func (b *builder) WithStrategy(strategy Strategy) *builder {
	if strategy == StrategyPeriod && !times.HasBirthTime {
		panic("os not supported file birthtime")
	}
	b.r.options.Strategy = strategy
	return b
}

func (b *builder) WithSize(size int64) *builder {
	b.r.options.Size = size
	return b
}

func (b *builder) WithPeriod(period time.Duration) *builder {
	b.r.options.Period = period
	return b
}

func (b *builder) WithFilePrefix(filePrefix string) *builder {
	b.r.options.FilePrefix = filePrefix
	return b
}

func (b *builder) WithCount(count int) *builder {
	b.r.options.Count = count
	return b
}

func (b *builder) WithExt(ext string) *builder {
	b.r.options.Ext = ext
	return b
}

func (b *builder) WithDirname(dirname string) *builder {
	b.r.options.Dirname = dirname
	return b
}

func (b *builder) WithCheckPeriod(checkPeriod time.Duration) *builder {
	b.r.options.CheckPeriod = checkPeriod
	return b
}

func (b *builder) Build() *Rotator {
	return &b.r
}

func NewBuilder() *builder {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &builder{
		r: Rotator{
			options: Options{
				Strategy:    StrategyNone,
				Size:        1 * 1024 * 1024,
				Period:      24 * time.Hour,
				FilePrefix:  defaultPrefix,
				Count:       10,
				Dirname:     wd,
				Ext:         defaultExt,
				CheckPeriod: defaultCheckPeriod,
			},
			filename: getFilename(wd, defaultPrefix, defaultExt, 0),
		},
	}
}
