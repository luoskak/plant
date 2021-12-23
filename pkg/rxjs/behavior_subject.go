package rxjs

import (
	"context"

	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

type behaviorSubject struct {
	*subject
	behavior interface{}
	self     abstract.Subscription
}

func (o *behaviorSubject) Subscribe(callback abstract.NextCallbacker) abstract.Subscription {
	sub := o.subject.Subscribe(callback)
	o.subscriptions[len(o.subscriptions)-1].next(o.behavior)
	return sub
}

func BehaviorSubject(ctx context.Context, init interface{}) abstract.Subject {
	b := new(behaviorSubject)
	b.subject = Subject(ctx).(*subject)
	subscriber := &behaviorSubscriber{b: b}
	b.self = b.Subscribe(subscriber)
	b.behavior = init
	return b
}

type behaviorSubscriber struct {
	b *behaviorSubject
}

func (bs *behaviorSubscriber) Next(v interface{}) {
	bs.b.behavior = v
}
