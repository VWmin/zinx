package znet

import "sync"

type BaseProperties struct {
	lock sync.RWMutex
	kv map[string] interface{}
}

func (b *BaseProperties) SetProperties(key string, val interface{}) {
	b.lock.Lock()
	defer b.lock.Unlock()

	// fixme: 怎样做到不创建对象，使继承该类的类能够直接使用这些方法
	if b.kv == nil {
		b.kv = make(map[string] interface{})
	}

	b.kv[key] = val
}

func (b *BaseProperties) GetProperties(key string) (interface{}, bool) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.kv == nil {
		b.kv = make(map[string] interface{})
	}

	s, ok := b.kv[key]
	return s, ok
}



