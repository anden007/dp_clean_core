package base

import (
	"context"

	"github.com/anden007/dp_clean_core/pkg"
)

type ICRUDUsecase[T IDBModel] interface {
	Add(entity T, ctx context.Context, transId string) (err error)
	DelByIds(ids []string, ctx context.Context, transId string) (err error)
	Edit(entity T, ctx context.Context, transId string) (err error)
	Updates(id string, fieldValues map[string]interface{}, ctx context.Context, transId string) (err error)
	GetAll() (result []T, err error)
	GetById(id string) (result T, err error)
	GetByIds(ids []string) (result []T, err error)
	GetByCondition(condition map[string]pkg.SearchCondition) (result []T, total int64, err error)
	Commit(transId string) (err error)
	Rollback(transId string) (err error)
}
