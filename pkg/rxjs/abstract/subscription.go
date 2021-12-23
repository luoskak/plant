package abstract

type OnceSubscription interface {
	Unsubscribe()
	Add(teardown interface{})
	// Done 某些生产者只能在程序结束了才会获得该channel
	// 调用该接口并获取到channel后将同时取消订阅
	Done() <-chan struct{}
	// Once 执行一次就获得channel,可确保获取到至少1次数据
	// 调用该接口并获取到channel后将同时取消订阅
	Once() <-chan struct{}
	// Try 执行过一次就获得channel,可确保获取到至少1次数据
	// 获取后不取消订阅
	Try() <-chan struct{}

	Remove(Subscription)
	IsClosed() bool
}

type Subscription interface {
	OnceSubscription
	// Keeper 一直尝试，直到取消订阅或该订阅结束
	// 只有接收到调用接口后的数据才触发
	Keeper() <-chan struct{}
}
