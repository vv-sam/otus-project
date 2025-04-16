package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ch := make(chan fmt.Stringer)
	defer close(ch)

	// Чтобы программа не завершилась
	wg := &sync.WaitGroup{}
	wg.Add(3)
	defer wg.Wait()

	// Каждые 60мс кладём в канал структурку
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Generate struct: context is done")
				wg.Done()
				return
			case <-time.After(60 * time.Millisecond):
				service.GenerateStruct(ch)
			}

		}
	}()

	go service.ConsumeStructs(ctx, ch, wg)

	// Каждые 200мс проверяем обновления
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Check updates: context is done")
				wg.Done()
				return
			case <-time.After(60 * time.Millisecond):
				repository.CheckUpdates()
			}
		}
	}()
}
