package loaddecider

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"math/rand"
	"sync/atomic"
	"time"
)

// Gauge 测量系统的负载
type Gauge interface {
	// Load 获取当前机器的负载
	Load() int32
	Close()
}

// RandomLoadGauge 随机生成负载器
type RandomLoadGauge struct {
	// 作业提示
	// 随机生成一个，就代表当前负载。你可以每隔一分钟生成一个
	load int32
	// 定时器
	timer *time.Ticker
	// 是否关闭模拟的负载获取的 goroutine
	done chan bool
	// logger
	logger logger.LoggerV1
	// Reporter to redis
	reporter Reporter
	// identity
	serverId string
	closed   atomic.Bool
}

func (r *RandomLoadGauge) Close() {
	if r.closed.CompareAndSwap(false, true) {
		r.done <- true
	}
}

func NewDefaultLoadGauge(logger logger.LoggerV1, reporter Reporter, id Identity) Gauge {
	return NewRandomLoadGauge(logger, reporter, id, time.Second*30)
}

func NewRandomLoadGauge(logger logger.LoggerV1, reporter Reporter, id Identity, checkInterval time.Duration) Gauge {
	rg := RandomLoadGauge{
		timer:    time.NewTicker(checkInterval),
		done:     make(chan bool),
		logger:   logger,
		reporter: reporter,
		serverId: id.Id(),
	}
	rg.initLoadUpdater()
	return &rg
}

// initLoadUpdater 开启 goroutine 定时处理负载更新的模拟.
func (r *RandomLoadGauge) initLoadUpdater() {
	// 上来先看看负载,然后汇报一次.
	r.updateRandomLoad()
	r.reporter.Report(context.Background(), r.serverId, r.load)

	// 定期更新负载信息.
	go func() {
		for {
			select {
			case <-r.timer.C:
				r.updateRandomLoad()
				// 性能指标汇报到 reporter
				r.reporter.Report(context.Background(), r.serverId, r.load)
			case <-r.done:
				return
			}
		}
	}()
}

// updateRandomLoad 模拟更新当前节点的负载.
func (r *RandomLoadGauge) updateRandomLoad() {
	// r.logger.Info(">> generate load!")
	atomic.StoreInt32(&r.load, rand.Int31n(100))
}

func (r *RandomLoadGauge) Load() int32 {
	return r.load
}
