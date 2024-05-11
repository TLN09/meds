package meds

import (
	"testing"
)

var msg = []byte("This is my message")
var parameterSets = []int{9923, 13220, 41711, 69497, 134180, 167717}

func empty(x []byte) bool {
	empty := true
	for i := 0; empty && i < len(x); i++ {
		empty = x[i] == 0
	}
	return empty
}

func TestKeyGen(test *testing.T) {
	for _, p := range parameterSets {
		test.Logf("MEDS-%v\n", p)
		ParameterSetup(p)
		sk, pk := KeyGen()
		if empty(sk) || empty(pk) {
			test.Errorf("Keys are empty\n")
		}
	}
}

func TestSign(test *testing.T) {
	msg := []byte("This is my message")
	for _, p := range parameterSets {
		test.Logf("MEDS-%v\n", p)
		ParameterSetup(p)
		_, sk := KeyGen()
		signed, err := Sign(sk, msg)
		if err != nil || empty(signed[:l_sig]) {
			test.Errorf("%v\n", err)
			test.Errorf("signed: %v\n", signed)
		}
	}
}

func TestVerify(test *testing.T) {
	for _, p := range parameterSets {
		test.Logf("MEDS-%v\n", p)
		ParameterSetup(p)
		pk, sk := KeyGen()
		msg_s, err := Sign(sk, msg)
		if err != nil {
			test.Errorf("%v\n", err)
		}
		// test.Logf("\nmsg: %v\nmsg_s:%v\n", string(msg), string(msg_s))
		msg_v := Verify(pk, msg_s)
		if msg_v == nil {
			test.Errorf("Invalid Signature MEDS-%v\n", p)
		}
	}
}

// func TestVerifyMEDS_134180(test *testing.T) {
// 	p := 134180
// 	test.Logf("MEDS-%v\n", p)
// 	ParameterSetup(p)
// 	pk, sk := KeyGen()
// 	msg_s, err := Sign(sk, msg)
// 	if err != nil {
// 		test.Errorf("%v\n", err)
// 	}
// 	// test.Logf("\nmsg: %v\nmsg_s:%v\n", string(msg), string(msg_s))
// 	msg_v := Verify(pk, msg_s)
// 	if msg_v == nil {
// 		test.Errorf("Invalid Signature MEDS-%v\n", p)
// 	}
// }

func BenchmarkKeygen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}

func BenchmarkSign(b *testing.B) {
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}

func BenchmarkVerify(b *testing.B) {
	pk, sk := KeyGen()
	signed, err := Sign(sk, msg)
	if err != nil {
		b.Fatal("Error is not nil")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := Verify(pk, signed)
		if msg == nil {
			b.Fatal("signature is invalid")
		}
	}
}
