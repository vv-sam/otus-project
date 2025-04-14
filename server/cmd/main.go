package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/service"
)

func main() {
	fmt.Println("Server is running...")

	ch := make(chan fmt.Stringer)

	// Чтобы программа не завершилась
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Каждые 60мс кладём в канал структурку
	go func() {
		for {
			service.GenerateStruct(ch)
			<-time.After(60 * time.Millisecond)
		}
	}()

	// эта функция сама завершится когда закроется канал
	go repository.PassStructs(ch)

	// Каждые 200мс проверяем обновления
	go func() {
		for {
			repository.CheckUpdates()
			<-time.After(200 * time.Millisecond)
		}
	}()

	wg.Wait()
	close(ch)
}
