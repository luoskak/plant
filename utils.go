package plant

import "reflect"

// make sure cv is pointer
func canReciveThis(ct reflect.Type, t uintptr) bool {
	switch ct.Kind() {
	case reflect.Struct:
		cPtr := unsafeTypeKey(reflect.PtrTo(ct))
		if cPtr == t {
			return true
		}
		for i := 0; i < ct.NumField(); i++ {
			ft := ct.Field(i).Type
			cPtr := unsafeTypeKey(reflect.PtrTo(ft))
			if cPtr == t {
				return true
			}
		}
	case reflect.Slice:

	}
	return false
}
