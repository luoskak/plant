package plant

import (
	"context"
	"fmt"

	"github.com/luoskak/mist"
)

const MiddlewareName = "plant"

var (
	_ mist.Runner = &Middleware{}
)

func init() {
	mist.DefaultManager.Register(MiddlewareName, &Middleware{})
}

type Middleware struct {
	proxies []*proxy
	opt     *mwOptions
}

func (m *Middleware) Inter(full bool) mist.Interceptor {
	return func(ctx context.Context, req interface{}, info *mist.ServerInfo, handler mist.Handler) (interface{}, error) {
		contain := &consumerContaining{
			opt:                m.opt,
			autoUnsubscription: true,
		}
		if full {
			contain.globalProxies = m.proxies
		} else {
			var proxies []*proxy
			for _, p := range m.proxies {
				if !p.full {
					proxies = append(proxies, p)
				}
			}
			contain.globalProxies = proxies
		}
		newCtx := context.WithValue(ctx, consumerKey, contain)
		res, err := handler(newCtx, req)
		// 取消订阅，释放内存
		if contain.autoUnsubscription {
			for _, sub := range contain.subs {
				sub.Unsubscribe()
			}
		}

		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (m *Middleware) Init(opt []mist.Option) {
	opts := defaultMwOptions
	for _, o := range opt {
		o.Apply(&opts)
	}
	for i, setUp := range opts.productorSetUps {
		p, mts, ct := reflectSetUp(opts.getName(i), setUp)
		for _, mt := range mts {
			opts.productorMaterialMap[mt] = i
		}
		if ct != 0 {
			opts.productorConsumerMap[ct] = i
		}

		m.proxies = append(m.proxies, p)
	}
	m.opt = &opts
}

func (m *Middleware) Run(ctx context.Context, info *mist.ServerInfo) error {
	for _, p := range m.proxies {
		err := p.productor.Run(ctx)
		if err != nil {
			panic(fmt.Sprintf("%s productor run got %v", p.name, err))
		}
	}
	return nil
}

func Put(ctx context.Context, material interface{}) error {
	m := mist.GetWm(MiddlewareName).(*Middleware)
	typ := unsafeKey(material)
	i, registered := m.opt.productorMaterialMap[typ]
	if !registered {
		return fmt.Errorf("material is not registered or productor is not support put")
	}
	return m.proxies[i].productor.Put(ctx, material)
}
