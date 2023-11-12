package cb

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	Open int32 = iota
	Close
	HalfOpen
)

var CbCloseError = errors.New("熔断状态，直接异常")

type CircuitBreakerService struct {
	// 状态为 0 - 打开， 1 - 关闭， 2 - 半开半闭
	state int32
	// 熔断次数.
	cnt int32
	// 熔断上限.
	threshold int32
	// 熔断恢复时间间隔.
	recoverTime time.Duration
	// 具体干的事情，熔断器本身不关心.
	doWhat func(args ...any) error
}

func NewCircuitBreaker(threshold int32, recoverTime time.Duration, doWhat func(args ...any) error) CircuitBreaker {
	return &CircuitBreakerService{
		threshold:   threshold,
		recoverTime: recoverTime,
		doWhat:      doWhat,
	}
}

func (c *CircuitBreakerService) Do(args ...any) error {
	state := atomic.LoadInt32(&c.state)
	cnt := atomic.LoadInt32(&c.cnt)
	switch state {
	case Open:
		err := c.doWhat(args...)
		if err != nil {
			// 处理 open 状态下到 error 异常
			c.openStateError(cnt, state)
			return err
		}
		return nil
	case Close:
		return CbCloseError
		err := c.doWhat(args)
		if err != nil {
			// 直接恢复 close 状态.
			c.toCloseState(state)
			return err
		} else {
			if cnt+1 >= c.threshold {
				// 可以恢复到 Open 状态了.
				c.toPenState(state)
			} else {
				// 又成功了一次.
				atomic.AddInt32(&c.cnt, 1)
			}
		}
		return nil
	}
	return errors.New("原则是不会到这里，到这里就是一种未知熔断器的状态.")
}

func (c *CircuitBreakerService) openStateError(cnt int32, state int32) {
	if cnt+1 >= c.threshold {
		c.toCloseState(state)
	} else {
		// 增加错误计数器.
		atomic.AddInt32(&c.cnt, 1)
	}
}

func (c *CircuitBreakerService) toCloseState(state int32) {
	if atomic.CompareAndSwapInt32(&c.state, state, Close) {
		// 重置计数器,进入短路状态.
		atomic.StoreInt32(&c.cnt, 0)
		time.AfterFunc(c.recoverTime, func() {
			// 恢复为半开状态.
			atomic.CompareAndSwapInt32(&c.state, Close, HalfOpen)
		})
	}
}

func (c *CircuitBreakerService) toPenState(state int32) {
	if atomic.CompareAndSwapInt32(&c.state, state, Open) {
		// 重置计数器,进入短路状态.
		atomic.StoreInt32(&c.cnt, 0)
	}
}
