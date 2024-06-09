package base

import (
	"context"

	"github.com/anden007/dp_clean_core/pkg"
)

type ICRUDRepository[T IDBModel] interface {
	IReadOnlyRepository[T]
	Add(entity T, ctx context.Context, transId string) (err error)
	DelByIds(ids []string, ctx context.Context, transId string) (err error)
	Edit(entity T, ctx context.Context, transId string) (err error)
	Updates(id string, fieldValues map[string]interface{}, ctx context.Context, transId string) (err error)
	Commit(transId string) (err error)
	Rollback(transId string) (err error)
}

type IReadOnlyRepository[T IDBModel] interface {
	GetAll() (result []T, err error)
	GetById(id string) (result T, err error)
	GetByIds(ids []string) (result []T, err error)
	GetByCondition(condition map[string]pkg.SearchCondition) (result []T, total int64, err error)
}
