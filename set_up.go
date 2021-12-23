package plant

import (
	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

// name the plant name
// full 是否需要全量权限
func RegisterSetUp(name string, full bool, setUp SetUp) {
	_, duplicated := defaultMwOptions.productorNameMap[name]
	if duplicated {
		panic("duplicated register setup")
	}
	defaultMwOptions.productorSetUps = append(defaultMwOptions.productorSetUps, setUp)
	defaultMwOptions.productorNameMap[name] = len(defaultMwOptions.productorSetUps) - 1
	defaultMwOptions.productorFullMap[name] = full
}

type SetUp func() (productor abstract.Productor, materials []interface{}, product interface{})
