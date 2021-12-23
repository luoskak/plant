package rxjs

import (
	"context"
	"sync/atomic"

	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

type subject struct {
	*observable
	observer Observer
}

func (o *subject) Next(v interface{}) {
	o.observer.Next(v)
}

func (o *subject) Subscribe(callback abstract.NextCallbacker) abstract.Subscription {
	sub := &Subscription{}
	sub.nextFunc = callback.Next
	if cb, ok := callback.(abstract.ErrorCallbacker); ok {
		sub.errFunc = cb.Err
	}
	o.subscriptions = append(o.subscriptions, sub)
	return sub
}

func Subject(ctx context.Context) abstract.Subject {
	s := new(subject)
	o := Observable(func(observer Observer) {
		s.observer = observer
		<-ctx.Done()
	})
	if atomic.CompareAndSwapInt32(&o.running, 0, 1) {
		go o.f()
	}
	s.observable = o
	return s
}
