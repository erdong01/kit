package aop

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type Goods struct {
	Aop
	GoodsNo int
}

func (g *Goods) Handler() {
	if g.GoodsNo <= 0 {
		g.Break(errors.New("商品不存在"))
	}
	fmt.Println("商品存在")
}

type Order struct {
	Aop
	OrderNo int
	Status  int8
}

func (o *Order) Handler() {
	o.Status = 1
	fmt.Println("order_no", o.OrderNo)
	o.Set("order_no", o.OrderNo)
}
func (o *Order) After() {
	order_no := o.Get("order_no")
	fmt.Println("Order order_no", order_no)
	fmt.Println("status", o.Status)
	o.Break(errors.New("错了"))
}

type Pay struct {
	Aop
	OrderNo int
	Status  int8
}

func (o *Pay) Handler() {
	order_no := o.Get("order_no")
	fmt.Println("Pay order_no", order_no)
}
func TestOrder(t *testing.T) {
	New(context.Background(), &Order{OrderNo: 22}).SetBefore(&Goods{}).SetAfter(&Pay{}).Run()
	// Add()
}

func Add() {
	var ctx context.Context = context.Background()
	ctx = context.WithValue(ctx, "order_no", 11)
	ctx.Value("order_no")
}
