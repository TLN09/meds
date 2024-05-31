package finiteField

import (
	"math/rand"
	"testing"
)

const q int = 4093

func TestNewFieldElm(t *testing.T) {
	elm := NewFieldElm(0, q)
	if elm.n != 0 {
		t.Errorf("expected: 0\tResult: %v", elm.n)
	}
}

func TestAdd(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := NewFieldElm(rand.Int(), q)
		y := NewFieldElm(rand.Int(), q)
		e := mod(int(x.n)+int(y.n), q)

		r := x.Add(y)
		if r.n != e {
			t.Errorf("e: %v\tr: %v", e, r.n)
		}
	}
}

func TestSub(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := NewFieldElm(rand.Int(), q)
		y := NewFieldElm(rand.Int(), q)
		e := mod(int(x.n)-int(y.n), q)

		r := x.Sub(y)
		if r.n != e {
			t.Errorf("e: %v\tr: %v", e, r.n)
		}
	}
}

func TestMul(t *testing.T) {
	for i := 0; i < 1000; i++ {
		x := NewFieldElm(rand.Int(), q)
		y := NewFieldElm(rand.Int(), q)
		e := mod(int(x.n)*int(y.n), q)

		r := x.Mul(y)
		if r.n != e {
			t.Errorf("e: %v\tr: %v", e, r.n)
		}
	}
}

func TestInv(t *testing.T) {
	x := NewFieldElm(3351, q)
	e := NewFieldElm(1, q)
	inv := x.Inv()

	r := x.Mul(inv)

	if !r.Equals(e) {
		t.Errorf("x: %v\nr: %v\ne: %v\ninv: %v", x, r, e, inv)
	}
}
