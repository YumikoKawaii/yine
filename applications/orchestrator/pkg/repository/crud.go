package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/constants"
)

type IRepository[T any] interface {
	Get(ctx context.Context, filter IFilter) (T, error)
	List(ctx context.Context, filter IFilter) ([]T, error)
	Save(context.Context, *T) (T, error)
	SaveIgnoreConflicts(ctx context.Context, model *T) (T, error)
	Upsert(ctx context.Context, model *T) (T, error)
	UpsertMany(ctx context.Context, models []T) ([]T, error)
	Update(context.Context, *T) error
	SaveMany(context.Context, []T) ([]T, error)
	SaveManyIgnoreConflicts(context.Context, []T) ([]T, error)
	Exec(context.Context, string, ...interface{}) error
}

type IFilter interface {
	ApplyFilter(db *gorm.DB) *gorm.DB
}

type Repository[T any] struct {
	DB *gorm.DB
}

type Querier interface {
	ApplyQuery(db *gorm.DB) *gorm.DB
}

func New[T any](db *gorm.DB) IRepository[T] {
	return Repository[T]{DB: db}
}

func (c Repository[T]) Get(ctx context.Context, filter IFilter) (T, error) {
	query := c.DB.WithContext(ctx)
	if filter != nil {
		query = filter.ApplyFilter(query)
	}

	var record T
	return record, query.First(&record).Error
}

func (c Repository[T]) List(ctx context.Context, filter IFilter) ([]T, error) {
	query := c.DB.WithContext(ctx)
	if filter != nil {
		query = filter.ApplyFilter(query)
	}

	records := make([]T, 0)
	return records, query.Find(&records).Error
}

func (c Repository[T]) Save(ctx context.Context, model *T) (T, error) {
	err := c.DB.WithContext(ctx).Create(model).Error
	return *model, err
}

func (c Repository[T]) SaveIgnoreConflicts(ctx context.Context, model *T) (T, error) {
	err := c.DB.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(model).Error
	return *model, err
}

func (c Repository[T]) Upsert(ctx context.Context, model *T) (T, error) {
	err := c.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(model).Error
	return *model, err
}

func (c Repository[T]) UpsertMany(ctx context.Context, models []T) ([]T, error) {
	if len(models) == constants.Zero {
		return nil, nil
	}
	err := c.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(models).Error
	return models, err
}

func (c Repository[T]) SaveMany(ctx context.Context, models []T) ([]T, error) {
	if len(models) == constants.Zero {
		return nil, nil
	}
	err := c.DB.WithContext(ctx).Create(models).Error
	return models, err
}

func (c Repository[T]) SaveManyIgnoreConflicts(ctx context.Context, models []T) ([]T, error) {
	if len(models) == constants.Zero {
		return nil, nil
	}
	err := c.DB.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(models).Error
	return models, err
}

func (c Repository[T]) Update(ctx context.Context, model *T) error {
	err := c.DB.WithContext(ctx).Updates(model).Error
	return err
}

func (c Repository[T]) Exec(ctx context.Context, sql string, values ...interface{}) error {
	err := c.DB.WithContext(ctx).Exec(sql, values...).Error
	return err
}
