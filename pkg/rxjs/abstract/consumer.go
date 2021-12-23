package abstract

type Consumer interface {
	Next(v interface{})
}

type ArgsConsumer interface {
	Args() []interface{}
}
