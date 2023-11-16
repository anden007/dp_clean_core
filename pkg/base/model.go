package base

type IDBModel interface {
	TableName() string
	ModelName() string
}

type IViewModel interface {
	ToDBModel() error
	ToViewModel() error
}

type BaseViewModel struct {
}

func (m *BaseViewModel) ToDBModel() (err error) {
	return nil
}

func (m *BaseViewModel) ToViewModel() (err error) {
	return nil
}
