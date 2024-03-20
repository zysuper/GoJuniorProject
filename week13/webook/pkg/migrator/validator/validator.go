package validator

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator/events"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Validator[T migrator.Entity] struct {
	// 你数据迁移，是不是肯定有
	base   *gorm.DB
	target *gorm.DB

	l         logger.LoggerV1
	producer  events.Producer
	direction string
	batchSize int
	utime     int64
	// <= 0 就认为中断
	// > 0 就认为睡眠
	sleepInterval time.Duration
	fromBase      func(ctx context.Context, offset int) ([]T, error)
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.LoggerV1,
	p events.Producer) *Validator[T] {
	res := &Validator[T]{base: base, target: target,
		l: l, producer: p, direction: direction, batchSize: 100}
	res.fromBase = res.fullFromBase
	return res
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	//err := v.validateBaseToTarget(ctx)
	//if err != nil {
	//	return err
	//}
	//return v.validateTargetToBase(ctx)

	var eg errgroup.Group
	eg.Go(func() error {
		return v.validateBaseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.validateTargetToBase(ctx)
	})
	return eg.Wait()
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := 0
	for {
		sources, err := v.fromBase(ctx, offset)
		// 超时或者取消了
		if err == context.DeadlineExceeded || err == context.Canceled {
			return nil
		}

		// 如果全部处理完了
		if err == nil && len(sources) == 0 {
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
			continue
		}

		if err != nil {
			// 查询出错了,这批都不处理了，下一批.
			v.l.Error("base -> target 查询 base 失败", logger.Error(err))
			offset += v.batchSize
			continue
		}

		v.validateTargetBatch(ctx, sources)
		offset += v.batchSize
	}
}

func (v *Validator[T]) validateTargetBatch(ctx context.Context, sources []T) error {
	var dst []T

	// 直接使用 in (IDS) 来一次查询出来.
	ids := slice.Map(sources, func(idx int, src T) int64 { return src.ID() })

	err := v.target.WithContext(ctx).
		Where("id in ?", ids).
		Find(&dst).Error

	if err != nil {
		ids_str := slice.Map(ids[:10], func(idx int, src int64) string {
			return strconv.FormatInt(src, 10)
		})

		v.l.Error("target 查询失败",
			logger.String("ids",
				// 省略太长的 ids.
				strings.Join(ids_str, ",")+" ..."),
			logger.Error(err))
		return err
	}

	// 创建 Map[ID]T， 方便和 src 数据的 ID 来查询比较对象。
	var targetIdxMap map[int64]T
	for _, d := range dst {
		targetIdxMap[d.ID()] = d
	}

	for idx, id := range ids {
		// 目标库能查询到时，需要真比较.
		if dst, ok := targetIdxMap[id]; ok {
			equal := sources[idx].CompareTo(dst)
			if !equal {
				// 要丢一条消息到 Kafka 上
				v.notify(sources[idx].ID(), events.InconsistentEventTypeNEQ)
			}
		} else {
			// 目标库查询不到，直接就是缺一条.
			v.notify(sources[idx].ID(), events.InconsistentEventTypeTargetMissing)
		}
	}

	return nil
}

func (v *Validator[T]) Full() *Validator[T] {
	v.fromBase = v.fullFromBase
	return v
}

func (v *Validator[T]) Incr() *Validator[T] {
	v.fromBase = v.incrFromBase
	return v
}

func (v *Validator[T]) Utime(t int64) *Validator[T] {
	v.utime = t
	return v
}

func (v *Validator[T]) SleepInterval(interval time.Duration) *Validator[T] {
	v.sleepInterval = interval
	return v
}

func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) ([]T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src []T
	err := v.base.WithContext(dbCtx).Order("id").
		Offset(offset).
		Limit(v.batchSize).
		Find(&src).Error
	return src, err
}

func (v *Validator[T]) incrFromBase(ctx context.Context, offset int) ([]T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src []T
	err := v.base.WithContext(dbCtx).
		Where("utime > ?", v.utime).
		Order("utime").
		Offset(offset).
		Limit(v.batchSize).
		Find(&src).Error
	return src, err
}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) error {
	offset := 0
	for {
		var ts []T
		err := v.target.WithContext(ctx).Select("id").
			//Where("utime > ?", v.utime).
			Order("id").Offset(offset).Limit(v.batchSize).
			Find(&ts).Error
		if err == context.DeadlineExceeded || err == context.Canceled {
			return nil
		}
		if err == gorm.ErrRecordNotFound || len(ts) == 0 {
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
			continue
		}
		if err != nil {
			v.l.Error("target => base 查询 target 失败", logger.Error(err))
			offset += len(ts)
			continue
		}
		// 在这里
		var srcTs []T
		ids := slice.Map(ts, func(idx int, t T) int64 {
			return t.ID()
		})
		err = v.base.WithContext(ctx).Select("id").
			Where("id IN ?", ids).Find(&srcTs).Error
		if err == gorm.ErrRecordNotFound || len(srcTs) == 0 {
			// 都代表。base 里面一条对应的数据都没有
			v.notifyBaseMissing(ts)
			offset += len(ts)
			continue
		}
		if err != nil {
			v.l.Error("target => base 查询 base 失败", logger.Error(err))
			// 保守起见，我都认为 base 里面没有数据
			// v.notifyBaseMissing(ts)
			offset += len(ts)
			continue
		}
		// 找差集，diff 里面的，就是 target 有，但是 base 没有的
		diff := slice.DiffSetFunc(ts, srcTs, func(src, dst T) bool {
			return src.ID() == dst.ID()
		})
		v.notifyBaseMissing(diff)
		// 说明也没了
		if len(ts) < v.batchSize {
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
		}
		offset += len(ts)
	}
}

func (v *Validator[T]) notifyBaseMissing(ts []T) {
	for _, val := range ts {
		v.notify(val.ID(), events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Type:      typ,
		Direction: v.direction,
	})
	if err != nil {
		v.l.Error("发送不一致消息失败",
			logger.Error(err),
			logger.String("type", typ),
			logger.Int64("id", id))
	}
}
