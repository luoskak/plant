package rxjs

import (
	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

var (
	_ abstract.Subscription = &Subscription{}
)

type Subscription struct {
	Closed          bool
	received        bool
	nextFunc        func(interface{})
	completeFunc    func()
	errFunc         func(error)
	subscriptions   []*Subscription
	parentOrParents []*Subscription

	keepChannelClosed bool
	keepChannel       chan struct{}
}

func (sub *Subscription) IsClosed() bool {
	return sub.Closed
}

func (sub *Subscription) next(v interface{}) {
	sub.received = true
	if sub.nextFunc != nil {
		sub.nextFunc(v)
	}
	for _, child := range sub.subscriptions {
		if child.Closed {
			continue
		}
		child.next(v)
	}
}

func (sub *Subscription) complete() {
	if sub.completeFunc != nil {
		sub.completeFunc()
	}
	for _, child := range sub.subscriptions {
		if child.Closed {
			continue
		}
		child.complete()
	}
}

func (sub *Subscription) err(err error) {
	if sub.errFunc != nil {
		sub.errFunc(err)
	}
	for _, child := range sub.subscriptions {
		if child.Closed {
			continue
		}
		child.err(err)
	}
}

func (sub *Subscription) Unsubscribe() {
	sub.Closed = true
	if sub.keepChannel != nil {
		if !sub.keepChannelClosed {
			sub.keepChannelClosed = true
			close(sub.keepChannel)
		}
	}
	for _, parent := range sub.parentOrParents {
		parent.Unsubscribe()
	}
	for _, child := range sub.subscriptions {
		child.Unsubscribe()
	}
}

// Add teardown
// when subscription received the none argument function will be call directly
func (sub *Subscription) Add(teardown interface{}) {
	switch tf := teardown.(type) {
	case func(v interface{}):
		child := &Subscription{}
		child.nextFunc = tf
		sub.subscriptions = append(sub.subscriptions, child)
	case func():
		child := &Subscription{}
		child.completeFunc = tf
		sub.subscriptions = append(sub.subscriptions, child)
	}
}

func (sub *Subscription) Done() <-chan struct{} {
	d := make(chan struct{})
	child := &Subscription{}
	child.completeFunc = func() {
		sub.Unsubscribe()
		close(d)
	}
	sub.subscriptions = append(sub.subscriptions, child)
	return d
}

func (sub *Subscription) Once() <-chan struct{} {
	d := make(chan struct{}, 1)

	if sub.received {
		d <- struct{}{}
	}
	closed := false

	child := &Subscription{}
	child.nextFunc = func(v interface{}) {
		sub.Unsubscribe()
		if !closed {
			closed = true
			close(d)
		}
	}
	child.completeFunc = func() {
		sub.Unsubscribe()
		if !closed {
			closed = true
			close(d)
		}
	}
	sub.subscriptions = append(sub.subscriptions, child)
	return d
}

func (sub *Subscription) Try() <-chan struct{} {
	d := make(chan struct{}, 1)
	if sub.received {
		d <- struct{}{}
	}
	closed := false
	child := &Subscription{}
	child.nextFunc = func(v interface{}) {
		if !closed {
			closed = true
			close(d)
		}
	}
	child.completeFunc = func() {
		if !closed {
			closed = true
			close(d)
		}
	}
	sub.subscriptions = append(sub.subscriptions, child)
	return d
}

func (sub *Subscription) Keeper() <-chan struct{} {
	child := &Subscription{
		keepChannel: make(chan struct{}),
	}
	child.nextFunc = func(v interface{}) {
		child.keepChannel <- struct{}{}
	}
	child.completeFunc = func() {
		if !sub.keepChannelClosed {
			sub.keepChannelClosed = true
			close(child.keepChannel)
		}

	}
	sub.subscriptions = append(sub.subscriptions, child)
	return child.keepChannel
}

func (sub *Subscription) Remove(subscription abstract.Subscription) {
	for idx, child := range sub.subscriptions {
		if child == subscription {
			if idx == 0 {
				sub.subscriptions = sub.subscriptions[1:]
				break
			}
			if idx == len(sub.subscriptions)-1 {
				sub.subscriptions = sub.subscriptions[:idx]
				break
			}
			sub.subscriptions = append(sub.subscriptions[0:idx], sub.subscriptions[idx+1:]...)
		}
	}
}
