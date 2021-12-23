package abstract

type Observer interface {
	Next(interface{})
	Complete()
	Err(error)
}
