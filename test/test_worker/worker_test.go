package testworker

import (
	"fmt"
	"testing"

	"github.com/erdong01/kit/goWorker"
)

func TestWork(t *testing.T) {
	goworker := goWorker.New()
	goworker.Go(func() {
		fmt.Println("111")
	})
	goworker.Go(func() {
		fmt.Println("222")
	})

	goworker.Go(func() {
		fmt.Println("333")
	})

	goworker.Go(func() {
		fmt.Println("444")
	})

	goworker.Go(func() {
		fmt.Println("555")
	})
	goworker.Go(func() {
		fmt.Println("555")
	})
	goworker.Go(func() {
		fmt.Println("777")
	})
	goworker.Go(func() {
		fmt.Println("888")
	})
	goworker.Go(func() {
		fmt.Println("999")
	})

	goworker.Wait()
}
