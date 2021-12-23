package rxjs

import "github.com/luoskak/plant/pkg/rxjs/abstract"

// Map observer
func Map(fi interface{}) func(interface{}) interface{} {
	switch f := fi.(type) {
	case func(interface{}):
		return func(v interface{}) interface{} {
			f(v)
			return nil
		}
	case func(interface{}) interface{}:
		return f
	}
	return nil
}

func SwitchMap(f func(interface{}) abstract.Observable) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		abstractOb := f(v)
		switch ob := abstractOb.(type) {
		case *observable:
			// 碾平
			var r interface{}
			ob.nextFunc = func(v interface{}) {
				r = v
			}
			ob.f()
			return r
		default:
			panic("unsupported Observable for rxjs SwitchMap")
		}
	}
}
