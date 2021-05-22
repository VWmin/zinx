package ziface

type IProperties interface {
	SetProperties(key string, val interface{})

	GetProperties(key string) (interface{}, bool)
}
