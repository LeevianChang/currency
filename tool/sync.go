package tool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	Timeout = time.Second * 5
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func Go(f func(ctx context.Context) error) {
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Sprintf("panic %+v\n", err))
			}
			wg.Done()
		}()
		c, cancel := context.WithCancel(ctx)
		defer cancel()
		if err := f(c); err != nil {
			fmt.Println(err)
		}
	}()
}

func Ticker(f func(ctx context.Context) error, d time.Duration) {
	Go(func(ctx context.Context) error {
		t := time.NewTicker(d)
		for {
			select {
			case <-t.C:
				{
					f(ctx)
				}
			case <-ctx.Done():
				{
					return nil
				}
			}
		}
	})
}

func Close() {
	cancel()
	c, _ := context.WithTimeout(context.Background(), Timeout)
	finish := make(chan struct{})
	go func() {
		wg.Wait()
		finish <- struct{}{}
	}()
	select {
	case <-c.Done():
		{
			fmt.Println(fmt.Sprintf("ctx err %+v", ctx))
			return
		}
	case <-finish:
		{
			return
		}
	}
}
