package finiteField

import (
	"fmt"
)

type Fq struct {
	n int
	q int
}

func mod(n int, q int) int {
	r := n % q
	if r < 0 {
		r += q
	}
	return r
}

func NewFieldElm(n int, q int) *Fq {
	f := Fq{mod(n, q), q}
	return &f
}

func (x *Fq) Set(n int) {
	x.n = mod(n, x.q)
}

func (x *Fq) Equals(y *Fq) bool {
	return x.q == y.q && x.n == y.n
}

func (x *Fq) Add(y *Fq) *Fq {
	n := x.n + y.n
	return NewFieldElm(n, x.q)
}

func (x *Fq) Sub(y *Fq) *Fq {
	n := x.n - y.n
	return NewFieldElm(n, x.q)
}

func (x *Fq) Mul(y *Fq) *Fq {
	n := x.n * y.n
	return NewFieldElm(n, x.q)
}

func Inverse(a, b int) int {
	x := 1
	y := 0
	x1 := 0
	y1 := 1
	a1 := a
	b1 := b
	for b1 != 0 {
		q := a1 / b1
		x, x1 = x1, x-q*x1
		y, y1 = y1, y-q*y1
		a1, b1 = b1, a1-q*b1
	}
	return x
}

func (x *Fq) Inv() *Fq {
	return NewFieldElm(Inverse(int(x.n), int(x.q)), x.q)
}

func (x *Fq) UnaryMinus() *Fq {
	return NewFieldElm(-int(x.n), x.q)
}

func (x *Fq) BitLen() int {
	return 16
	// return int(math.Ceil(math.Log2(float64(x.q))))
}

func (x *Fq) Bytes() []byte {
	b := make([]byte, x.BitLen()/8)
	b[0] = byte((x.n & 0xff00) >> 8)
	b[1] = byte(x.n & 0x00ff)
	return b
}

func NewFromBytes(b []byte, q int) *Fq {
	return NewFieldElm(int((int(b[0])<<8)+int(b[1])), q)
}

func (x *Fq) String() string {
	return fmt.Sprintf("%v", x.n)
}

func (x *Fq) Value() int {
	return x.n
}
