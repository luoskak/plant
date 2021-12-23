package abstract

import (
	"context"
)

type Productor interface {
	Put(ctx context.Context, material interface{}) error
	// BoM() interface{}
	Run(ctx context.Context) error
	// Consumer() Consumer
	Consume(consumer Consumer) Subscription
}
