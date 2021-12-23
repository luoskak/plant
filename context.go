package plant

import (
	"context"
	"reflect"
	"sync"

	"github.com/luoskak/logger"
	"github.com/luoskak/mist"
	"github.com/luoskak/plant/pkg/rxjs"
	"github.com/luoskak/plant/pkg/rxjs/abstract"
)

var (
	consumerKey    mist.ContextKey = "consumeKey"
	containerCache sync.Map
)

type consumerContaining struct {
	globalProxies      []*proxy
	opt                *mwOptions
	subs               []abstract.Subscription
	autoUnsubscription bool
}

// Consume 从生产者处消费数据，请不要对指针类型数据进行修改
func Consume(ctx context.Context, containers ...interface{}) (abstract.OnceSubscription, bool) {
	if len(containers) == 0 {
		return nil, false
	}
	return consume(ctx, containers, true)
}

// KeepConsume 持续从生产者处消费数据，请不要对指针类型数据进行修改
func KeepConsume(ctx context.Context, containers ...interface{}) (abstract.Subscription, bool) {
	if len(containers) == 0 {
		return nil, false
	}
	return consume(ctx, containers, false)
}

// Consume args[0] bool set autoUnsubscription
func consume(ctx context.Context, containers []interface{}, args ...interface{}) (abstract.Subscription, bool) {
	for _, container := range containers {
		if reflect.ValueOf(container).Type().Kind() != reflect.Ptr {
			panic("container can not be set value")
		}
	}
	containing := ctx.Value(consumerKey).(*consumerContaining)
	if len(args) > 0 {
		switch a1 := args[0].(type) {
		case bool:
			containing.autoUnsubscription = a1
		}
	}
	var subs []abstract.Subscription
	var dataContainers []interface{}
	for _, container := range containers {
		if consumer, ok := container.(abstract.Consumer); ok {
			var proxies []*proxy
			key := unsafeKey(consumer)
			if i, ok := containerCache.Load(key); ok {
				proxies = i.([]*proxy)
			} else {
				ct := reflect.ValueOf(container).Elem().Type()

				for t, i := range containing.opt.productorConsumerMap {
					if canReciveThis(ct, t) {
						p := containing.globalProxies[i]
						proxies = append(proxies, p)
					}
				}
				containerCache.Store(key, proxies)
			}
			for _, p := range proxies {
				p.lastMsg = mist.WriteRuntimeMsg()
				p.subTimes++
				sub := p.productor.Consume(consumer)
				subs = append(subs, sub)
			}
		} else {
			dataContainers = append(dataContainers, container)
		}
	}
	// 自动合并数据容器
	if dynamicConsumer := createDynamicConsumer(dataContainers...); dynamicConsumer != nil {
		var proxies []*proxy
		key := unsafeKey(dynamicConsumer)
		if i, ok := containerCache.Load(key); ok {
			proxies = i.([]*proxy)
		} else {
			for _, container := range dataContainers {
				containerType := unsafeTypeKey(reflect.ValueOf(container).Type())
				for t, i := range containing.opt.productorConsumerMap {
					if containerType == t {
						p := containing.globalProxies[i]
						proxies = append(proxies, p)
					}
				}
			}
			containerCache.Store(key, proxies)
		}
		for _, p := range proxies {
			p.lastMsg = mist.WriteRuntimeMsg()
			p.subTimes++
			sub := p.productor.Consume(dynamicConsumer)
			subs = append(subs, sub)
		}

	}
	if len(subs) > 0 {
		var (
			sub abstract.Subscription
		)
		if containing.autoUnsubscription {
			sub = rxjs.Observable(func(observer rxjs.Observer) {
				for _, sub := range subs {
					<-sub.Once()
				}
				observer.Next(nil)
				observer.Complete()
			}).Subscribe(&emptySubscribe{})
		} else {
			sub = rxjs.Observable(func(observer rxjs.Observer) {
				for _, sub := range subs {
					<-sub.Try()
				}
				observer.Next(nil)
				// support for keep
				var keeps []<-chan struct{}
				for _, sub := range subs {
					keeps = append(keeps, sub.Keeper())
				}
				stashArr := make([]bool, len(keeps))
				got := func() bool {
					for _, stash := range stashArr {
						if !stash {
							return false
						}
					}
					return true
				}
				wg := sync.WaitGroup{}
				for i, keep := range keeps {
					wg.Add(1)
					go func(kept <-chan struct{}, index int) {
						for range kept {
							stashArr[index] = true
							if got() {
								observer.Next(nil)
							}
						}
						wg.Done()
					}(keep, i)
				}
				wg.Wait()
				observer.Complete()
			}).Subscribe(&emptySubscribe{})
		}
		containing.subs = append(containing.subs, sub)
		return sub, true
	}
	logger.Error("无法消费，是否已注册productor或者非全量权限")
	return nil, false
}

type emptySubscribe struct{}

func (s *emptySubscribe) Next(v interface{}) {}
