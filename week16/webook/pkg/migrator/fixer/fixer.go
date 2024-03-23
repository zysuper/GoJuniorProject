package fixer

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator"
	"gitee.com/geekbang/basic-go/webook/pkg/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OverrideFixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	columns []string
}

func NewOverrideFixerV1[T migrator.Entity](base *gorm.DB, target *gorm.DB,
	columns []string) *OverrideFixer[T] {
	return &OverrideFixer[T]{base: base, target: target, columns: columns}
}

func NewOverrideFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) (*OverrideFixer[T], error) {
	rows, err := base.Model(new(T)).Order("id").Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	return &OverrideFixer[T]{base: base, target: target, columns: columns}, err
}

func (f *OverrideFixer[T]) Fix(ctx context.Context, id int64) error {
	// 最最粗暴的
	var t T
	err := f.base.WithContext(ctx).Where("id=?", id).First(&t).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return f.target.WithContext(ctx).Model(&t).Delete("id = ?", id).Error
	case nil:
		// upsert
		return f.target.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&t).Error
	default:
		return err
	}
}

func (f *OverrideFixer[T]) FixV1(evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeNEQ, events.InconsistentEventTypeTargetMissing:
		var t T
		err := f.base.Where("id=?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.Model(&t).Delete("id = ?", evt.ID).Error
		case nil:
			// upsert
			return f.target.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.Model(new(T)).Delete("id = ?", evt.ID).Error
	}
	return nil
}
