package abstract

type Observable interface {
	Subscribe(NextCallbacker) Subscription
	Pipe(...OperatorFunc) Observable
}

type OperatorFunc func(interface{}) interface{}
