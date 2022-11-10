package aop

import (
	"context"
	"testing"
)

type Order struct {
	Base
	OrderNo int
	Status  int8
}

func (o *Order) Handler() {
	o.Status = 1
	// fmt.Println("order_no", o.OrderNo)
	o.Set("order_no", o.OrderNo)
}
func (o *Order) After() {
	// fmt.Println("status", o.Status)
}

type Pay struct {
	Base
	OrderNo int
	Status  int8
}

func (o *Pay) Handler() {
	o.Get("order_no")
	// fmt.Println("order_no", order_no)
}
func (o *Pay) After() {
	// fmt.Println("status", o.Status)
}
func TestXxx(t *testing.T) {
	for i := 0; i < 100000000; i++ {
		New(context.Background(), &Order{OrderNo: 22}).Run()
		// Add()
	}
}

func Add() {
	var ctx context.Context = context.Background()
	ctx = context.WithValue(ctx, "order_no", 11)
	ctx.Value("order_no")
}
