package aop

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/erDong01/micro-kit/aop"
)

type Form struct {
	GoodsNo  int //商品编号
	GoodsNum int //商品数量
}

type Goods struct {
	aop.Aop
	form Form
}

func (g *Goods) Handler() {
	if g.form.GoodsNo <= 0 {
		g.Break(errors.New("商品不存在")) //结束执行
		return
	}
	fmt.Println("商品存在")
}

type Order struct {
	aop.Aop
	form    Form
	OrderNo int
	Status  int8
}

func (o *Order) Handler() {
	o.OrderNo = 123456
	o.Status = 1
	fmt.Println("创建订单")
}
func (o *Order) After() {
	o.Set("order_no", o.OrderNo)
}

type Pay struct {
	aop.Aop
	OrderNo int
}

func (p *Pay) Before() {
	p.OrderNo = p.Get("order_no").(int)
}

func (p *Pay) Handler() {
	fmt.Println("为订单号：", p.OrderNo, "创建支付订单")
}
func TestOrder(t *testing.T) {
	aop.New(context.Background(), &Order{}).SetBefore(&Goods{
		form: Form{GoodsNo: 11, GoodsNum: 1},
	}).SetAfter(&Pay{}).Run()
	// Add()
}

func Add() {
	var ctx context.Context = context.Background()
	ctx = context.WithValue(ctx, "order_no", 11)
	ctx.Value("order_no")
}
