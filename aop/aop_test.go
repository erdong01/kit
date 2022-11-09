package aop

import (
	"context"
	"fmt"
	"testing"
)

type Order struct {
	Base
	OrderNo int
}

func (o *Order) Handler() {
	fmt.Println("order_no", o.OrderNo)
}

func TestXxx(t *testing.T) {
	New(context.Background(), &Order{OrderNo: 11}).Run()
}
