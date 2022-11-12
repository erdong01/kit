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
	// fmt.Println("商品存在")
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
	// fmt.Println("创建订单")
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
	// fmt.Println("为订单号：", p.OrderNo, "创建支付订单")
}

type Stock struct {
	aop.Aop
	form Form
}

func (s *Stock) Handler() {
	// fmt.Println("减去库存数量：", s.form.GoodsNum)
}

func TestOrder(t *testing.T) {
	form := Form{GoodsNo: 11, GoodsNum: 1}
	for i := 0; i < 100000000; i++ {
		err := aop.New(context.Background(), &Order{}).SetBefore(&Goods{
			form: form,
		}).SetAfter(&Pay{}, &Stock{form: form}).Run()
		if err != nil {
			// fmt.Println("订单创建失败")
		}
	}
	fmt.Println("订单创建成功")
	// Create()
}

func Create(form Form) {
	var ctx context.Context = context.Background()

	goods := Goods{
		Aop:  aop.Aop{Ctx: ctx},
		form: form,
	}
	goods.Handler()
	order := Order{Aop: goods.Aop, form: form}
	order.Handler()
	order.After()
	pay := Pay{Aop: order.Aop}
	pay.Before()
	pay.Handler()
	stock := Stock{form: form}
	stock.Handler()
}
