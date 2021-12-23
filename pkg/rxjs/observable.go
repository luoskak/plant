package rxjs

import (
	"sync/atomic"

	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

var (
	_ abstract.Observable = &observable{}
)

type observable struct {
	// base
	subscriptions []*Subscription
	running       int32
	// make func called only once
	completeFuncCalled int32
	f                  func()
	nextFunc           func(interface{})
	errFunc            func(error)
	completeFunc       func()
	// broadcast
	parent   *observable
	children []*observable
	// operators
	mapFuncs []abstract.OperatorFunc
}

func (o *observable) next(v interface{}) {
	r := v
	for _, mapFunc := range o.mapFuncs {
		r = mapFunc(r)
	}
	for _, sub := range o.subscriptions {
		if sub.Closed {
			continue
		}
		sub.next(r)
	}
	o.freeClosedSubscriptions()
	for _, child := range o.children {
		child.next(r)
	}

}

func (o *observable) err(err error) {
	for _, sub := range o.subscriptions {
		if sub.Closed || sub.errFunc == nil {
			continue
		}
		sub.err(err)
	}
	for _, child := range o.children {
		child.err(err)
	}

}

func (o *observable) complete() {
	if !atomic.CompareAndSwapInt32(&o.completeFuncCalled, 0, 1) {
		return
	}
	for _, sub := range o.subscriptions {
		if sub.Closed || sub.completeFunc == nil {
			continue
		}
		sub.complete()
	}
	for _, child := range o.children {
		child.complete()
	}
}

func (o *observable) freeClosedSubscriptions() {
	var subs []*Subscription
	for _, sub := range o.subscriptions {
		if !sub.Closed {
			subs = append(subs, sub)
		}
	}
	o.subscriptions = subs
}

type Observer abstract.Observer

type observer struct {
	holder *observable
}

func (os observer) Next(v interface{}) {
	os.holder.nextFunc(v)
}

func (os observer) Complete() {

	os.holder.completeFunc()
}

func (os observer) Err(err error) {
	os.holder.errFunc(err)
}

func newObservable() *observable {
	o := &observable{}
	o.nextFunc = o.next
	o.completeFunc = o.complete
	o.errFunc = o.err
	return o
}

// Observable 可通过chan无限发送数据，但注意配合context小心使用
func Observable(f func(observer Observer)) *observable {
	o := newObservable()
	observer := observer{holder: o}
	o.f = func() { f(observer) }
	return o
}

func (o *observable) Subscribe(callback abstract.NextCallbacker) abstract.Subscription {
	sub := &Subscription{}
	sub.nextFunc = callback.Next
	if cb, ok := callback.(abstract.ErrorCallbacker); ok {
		sub.errFunc = cb.Err
	}
	if cb, ok := callback.(abstract.CompleteCallbacker); ok {
		sub.completeFunc = cb.Complete
	}
	o.subscriptions = append(o.subscriptions, sub)
	var (
		root = o
	)
	for root.parent != nil {
		root = root.parent
	}
	if atomic.CompareAndSwapInt32(&root.running, 0, 1) {
		go func() {
			root.f()
			o.completeFunc()
		}()
	}
	return sub
}

func (o *observable) Pipe(fs ...abstract.OperatorFunc) abstract.Observable {
	child := &observable{parent: o}
	child.mapFuncs = append(child.mapFuncs, fs...)
	o.children = append(o.children, child)
	return child
}
