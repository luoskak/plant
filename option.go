package plant

type mwOptions struct {
	// productors []abstract.Productor
	productorConsumerMap map[uintptr]int
	productorMaterialMap map[uintptr]int
	productorNameMap     map[string]int
	productorSetUps      []SetUp
	productorFullMap     map[string]bool
}

func (mopt mwOptions) getName(index int) string {
	for n, i := range mopt.productorNameMap {
		if i == index {
			return n
		}
	}
	panic("no name")
}

var defaultMwOptions = mwOptions{
	productorMaterialMap: make(map[uintptr]int),
	productorConsumerMap: make(map[uintptr]int),
	productorNameMap:     make(map[string]int),
	productorFullMap:     make(map[string]bool),
}
