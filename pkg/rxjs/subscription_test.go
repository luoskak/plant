package rxjs

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type keeyTryConsumer struct {
}

func (c *keeyTryConsumer) Next(v interface{}) {
	fmt.Println(v)
}

func TestKeepTry(t *testing.T) {

	ob := BehaviorSubject(context.Background(), 1)

	sub := ob.Subscribe(&keeyTryConsumer{})
	ctx := context.Background()
	go func() {
		<-sub.Try()
		fmt.Println("try done")
	}()
	go func() {
		i := 100
		keep := sub.Keeper()
		for {
			select {
			case _, ok := <-keep:
				if !ok {
					return
				}
				i++
				fmt.Printf("a keep %d\n", i)
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		i := 200
		for range sub.Keeper() {
			i++
			fmt.Printf("b keep %d\n", i)
			if i > 210 {
				sub.Unsubscribe()
			}
		}
	}()
	go func() {
		i := 300
		for range sub.Keeper() {
			i++
			fmt.Printf("c keep %d\n", i)
		}
	}()
	tk := time.NewTicker(time.Second)
	i := 0
	for {
		i++
		if i > 20 {
			return
		}
		<-tk.C
		ob.Next(i)
	}
}
