package main

import (
	"fmt"
	"time"

	"github.com/vv-sam/otus-project/server/internal/repository"
	"github.com/vv-sam/otus-project/server/internal/service"
)

func main() {
	fmt.Println("Server is running...")
	n := 3

	for range n {
		structs := service.GenerateStructs()
		repository.PassStructs(structs)
		repository.PrintValues()

		<-time.After(3 * time.Second)
	}
}
