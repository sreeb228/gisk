package gisk

import "sync"

type DslGetter struct {
	Getter DslGetterInterface
	Dsl    sync.Map
}

type DslGetterInterface interface {
	GetDsl(elementType ElementType, key string, version string) (string, error)
}

func (getter *DslGetter) GetDsl(elementType ElementType, key string, version string) (string, error) {
	//做缓存处理，保证dsl的一致性
	k := string(elementType) + "_" + key + "_" + version
	v, ok := operationMap.Load(k)
	if ok {
		return v.(string), nil
	}

	vv, err := getter.Getter.GetDsl(elementType, key, version)
	if err == nil {
		getter.Dsl.Store(k, vv)
	}
	return vv, err
}
