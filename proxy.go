package plant

import (
	"reflect"

	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

type proxy struct {
	name        string
	full        bool
	productor   abstract.Productor
	subTimes    int
	lastMsg     string
	SetConsumer func(reflect.Value, interface{}) error
}

func reflectSetUp(name string, setUp SetUp) (p *proxy, mts []uintptr, cPtr uintptr) {
	productor, materials, consumer := setUp()
	for _, material := range materials {
		mv := reflect.ValueOf(material)
		mt := mv.Type()
		if mt.Kind() != reflect.Ptr {
			mt = reflect.PtrTo(mt)
		}
		mts = append(mts, unsafeTypeKey(mt))
	}
	if consumer != nil {
		// cv := reflect.ValueOf(consumer)
		// ct = cv.Type()
		// if ct.Kind() == reflect.Ptr {
		// 	ct = ct.Elem()
		// }
		// if ct.Kind() == reflect.Interface {
		// 	panic(fmt.Sprintf("%s setUp consumer type is interface", name))
		// }
		cv := reflect.ValueOf(consumer)
		ct := cv.Type()
		if ct.Kind() != reflect.Ptr {
			ct = reflect.PtrTo(ct)
		}
		cPtr = unsafeTypeKey(ct)
	}
	p = &proxy{
		name:      name,
		productor: productor,
	}
	return
}
