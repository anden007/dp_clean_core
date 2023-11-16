package pkg

import (
	"reflect"

	"github.com/google/uuid"
)

type UUIDTransformer struct {
}

func (m *UUIDTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(uuid.UUID{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				dst.Set(src)
			}
			return nil
		}
	}
	return nil
}
