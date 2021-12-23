package plant

import (
	"reflect"
	"strconv"
	"sync"
	"unsafe"

	"github.com/modern-go/reflect2"
)

var (
	dynamicConsumerCache sync.Map
)

type eface struct {
	rtype unsafe.Pointer
	data  unsafe.Pointer
}

func unpackEFace(obj interface{}) *eface {
	return (*eface)(unsafe.Pointer(&obj))
}

func unsafeKey(v interface{}) uintptr {
	return uintptr(unpackEFace(v).rtype)
}

func unsafeTypeKey(typ reflect.Type) uintptr {
	return uintptr(unpackEFace(typ).data)
}

type consumer struct {
	setters map[uintptr]setter
}

func (consumer *consumer) Next(v interface{}) {
	key := unsafeKey(v)
	consumer.setters[key](v)
}

func createDynamicConsumer(containers ...interface{}) *consumer {
	if len(containers) == 0 {
		return nil
	}
	c := consumer{
		setters: make(map[uintptr]setter),
	}
	var (
		keyArr []uintptr
		key    = ""
	)
	for _, container := range containers {
		t1 := reflect.ValueOf(container).Type()
		key := unsafeTypeKey(t1.Elem())
		keyArr = append(keyArr, key)
	}

	for _, k := range keyArr {
		key += strconv.Itoa(int(k))
	}
	if i, ok := dynamicConsumerCache.Load(key); ok {
		return i.(*consumer)
	}
	for i, container := range containers {
		t1 := reflect.ValueOf(container).Type()
		t2 := reflect2.Type2(t1.Elem())
		container := containers[i]
		c.setters[keyArr[i]] = func(v interface{}) {
			t2.UnsafeSet(reflect2.PtrOf(container), reflect2.PtrOf(v))
		}
	}
	dynamicConsumerCache.Store(key, &c)
	return &c
}

type setter func(interface{})
