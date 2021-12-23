package abstract

type Callbacker interface{}

type NextCallbacker interface {
	Next(interface{})
}

type ErrorCallbacker interface {
	Err(error)
}

type CompleteCallbacker interface {
	Complete()
}

type StringCallbacker interface {
	Next(string)
}
