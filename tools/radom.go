package tools

import (
	"github.com/erDong01/micro-kit/base"
	"math/rand"
	"time"
)

type (
	Rand struct {
		*rand.Rand
	}
)

func (this *Rand) RandI(i int, n int) int {
	if i > n {
		base.Assert(false, "Rand::RandI: inverted range")
		return i
	}

	return int(i + this.Int()%(n-i+1))
}

func (this *Rand) RandF(i float32, n float32) float32 {
	if i > n {
		base.Assert(false, "Rand::RandF: inverted range")
		return i
	}

	return i + (n-i)*this.Float32()
}

var RAND = Rand{rand.New(rand.NewSource(time.Now().UnixNano()))}
