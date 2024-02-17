package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
}

type GORMJobDAO struct {
	db *gorm.DB
	// 续约周期.
	refreshInterval time.Duration
	// 续约失败周期倍率.
	rate int64
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db, refreshInterval: time.Minute, rate: 2}
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		var j Job
		now := time.Now().UnixMilli()
		// 作业：这里是缺少找到续约失败的 JOB 出来执行
		// `没人抢`或者是 被标记为 `jobStatusRunning` 状态，但是很久没被续约 `utime` 了.
		err := db.Where("(next_time <? AND status = ?) or (utime <? and status = 1)",
			now, jobStatusWaiting, dao.checkRefreshInterval(now)).
			First(&j).Error
		if err != nil {
			return j, err
		}
		res := db.WithContext(ctx).Model(&Job{}).
			Where("id = ? AND version = ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			// 没抢到
			continue
		}
		return j, err
	}
}

func (dao *GORMJobDAO) checkRefreshInterval(now int64) int64 {
	// dao.rate 倍时间的续约周期，都没有被续约，说明获取锁的那支程序挂了.
	return now - dao.refreshInterval.Milliseconds()*dao.rate
}

func (dao *GORMJobDAO) Release(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

func (dao *GORMJobDAO) UpdateUtime(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"utime": now,
	}).Error
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, jid int64, t time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"utime":     now,
		"next_time": t.UnixMilli(),
	}).Error
}

type Job struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Name       string `gorm:"type:varchar(128);unique"`
	Executor   string
	Expression string
	Cfg        string
	// 状态来表达，是不是可以抢占，有没有被人抢占
	Status int

	Version int

	NextTime int64 `gorm:"index"`

	Utime int64
	Ctime int64
}

const (
	// jobStatusWaiting 没人抢
	jobStatusWaiting = iota
	// jobStatusRunning 已经被人抢了
	jobStatusRunning
	// jobStatusPaused 不再需要调度了
	jobStatusPaused
)
