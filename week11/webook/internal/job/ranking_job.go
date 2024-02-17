package job

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/loaddecider"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	l       logger.LoggerV1
	timeout time.Duration
	client  *rlock.Client
	key     string

	localLock *sync.Mutex
	lock      *rlock.Lock

	// 作业提示
	// 随机生成一个，就代表当前负载。你可以每隔一分钟生成一个
	decider loaddecider.Decider
}

func NewRankingJob(
	svc service.RankingService,
	l logger.LoggerV1,
	client *rlock.Client,
	decider loaddecider.Decider,
	timeout time.Duration) *RankingJob {
	return &RankingJob{svc: svc,
		key:       "job:ranking",
		l:         l,
		client:    client,
		localLock: &sync.Mutex{},
		decider:   decider,
		timeout:   timeout}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// go fun() { r.Run()}

func (r *RankingJob) Run() error {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 如果负载比拼胜出
	if r.decider.IsVictory(cxt) {
		return r.doRun()
	} else {
		// 如果没有胜出，且有持有分布式锁，需要主动释放.
		r.localLock.Lock()
		defer r.localLock.Unlock()

		lock := r.lock
		if lock != nil {
			ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			// 释放分布式锁.
			err := lock.Unlock(ctx2)
			r.lock = nil

			if err != nil {
				r.l.Error("释放分布式锁失败", logger.Error(err))
				return err
			}
		}
	}
	return nil
}

func (r *RankingJob) doRun() error {
	r.localLock.Lock()
	lock := r.lock
	if lock == nil {
		// 抢分布式锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout,
			&rlock.FixIntervalRetry{
				Interval: time.Millisecond * 100,
				Max:      3,
				// 重试的超时
			}, time.Second)
		if err != nil {
			r.localLock.Unlock()
			r.l.Warn("获取分布式锁失败", logger.Error(err))
			return nil
		}
		r.lock = lock
		r.localLock.Unlock()
		go func() {
			// 并不是非得一半就续约
			er := lock.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				// 续约失败了
				// 你也没办法中断当下正在调度的热榜计算（如果有）
				r.localLock.Lock()
				r.lock = nil
				//lock.Unlock()
				r.localLock.Unlock()
			}
		}()
	}
	// 这边就是你拿到了锁
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

//func (r *RankingJob) Run() error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
//	defer cancel()
//	lock, err := r.client.Lock(ctx, r.key, r.timeout,
//		&rlock.FixIntervalRetry{
//			Interval: time.Millisecond * 100,
//			Max:      3,
//			// 重试的超时
//		}, time.Second)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		// 解锁
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//		defer cancel()
//		er := lock.Unlock(ctx)
//		if er != nil {
//			r.l.Error("ranking job释放分布式锁失败", logger.Error(er))
//		}
//	}()
//	ctx, cancel = context.WithTimeout(context.Background(), r.timeout)
//	defer cancel()
//
//	return r.svc.TopN(ctx)
//}
