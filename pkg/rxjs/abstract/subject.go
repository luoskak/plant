package abstract

type Subject interface {
	Observable
	Next(v interface{})
}
