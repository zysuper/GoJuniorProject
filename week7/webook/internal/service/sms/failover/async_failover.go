package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/circuit_breaker"
	"gitee.com/geekbang/basic-go/webook/pkg/limiter"
	"github.com/google/uuid"
	"log"
	"time"
)

var LimitedError = errors.New("服务太忙，稍后再试")

type AsyncFailoverService struct {
	cb        cb.CircuitBreaker
	limiter   limiter.Limiter
	repo      repository.MsgRepository
	key       string
	retryTime time.Duration
	retryCnt  int
}

func NewAsyncFailoverService(
	cb cb.CircuitBreaker,
	limiter limiter.Limiter,
	repo repository.MsgRepository,
	retryTime time.Duration,
	retryCnt int,
) sms.Service {
	return &AsyncFailoverService{
		cb:        cb,
		limiter:   limiter,
		repo:      repo,
		retryTime: retryTime,
		retryCnt:  retryCnt,
		key:       "async-sms-service"}
}

func (a *AsyncFailoverService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := a.limiter.Limit(ctx, a.key)

	// redis 崩了，
	if err != nil {
		log.Println(err)
	}

	if limited {
		// 1. 触发限流，要进行异步发送.
		a.asyncSave(ctx, tplId, args, numbers)
		return LimitedError
	}

	// 要么 redis 崩了， 要么没有崩，都是直接访问了.
	err = a.cb.Do(ctx, tplId, args, numbers)
	if err != nil {
		a.asyncSave(ctx, tplId, args, numbers)
		return err
	}

	return nil
}

// asyncSave 异步保存到数据库，然后异步重试.
func (a *AsyncFailoverService) asyncSave(ctx context.Context, tplId string, args []string, numbers []string) {
	go func() {
		id, err := a.repo.Create(ctx, domain.Msg{
			Id:      uuid.New().String(),
			TplId:   tplId,
			Args:    args,
			Numbers: numbers})
		if err != nil {
			log.Println("创建记录失败")
			return
		}
		a.asyncSend(ctx, id)
	}()
}

// asyncSend 异步重试.
func (a *AsyncFailoverService) asyncSend(ctx context.Context, id int64) {
	timer := time.NewTimer(a.retryTime)
	go func(t *time.Timer, id int64, a *AsyncFailoverService) {
		var tryCnt = 0
		defer timer.Stop()
		for {
			<-t.C
			msg, err := a.repo.FindById(ctx, id)
			if err != nil {
				tryCnt++
				if tryCnt >= a.retryCnt {
					break
				}
				t.Reset(a.retryTime)
				continue
			}

			err = a.cb.Do(ctx, msg.TplId, msg.Args, msg.Numbers)
			if err != nil {
				tryCnt++
				if tryCnt >= a.retryCnt {
					break
				}
				t.Reset(a.retryTime)
				continue
			}
			// 干成功了.
			break
		}
	}(timer, id, a)
}
