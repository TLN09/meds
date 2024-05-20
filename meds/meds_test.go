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

func TestVerify167717(test *testing.T) {
	ParameterSetup(167717)
	for i := 0; i < 10; i++ {
		pk, sk := KeyGen()
		msg_s, err := Sign(sk, msg)
		if err != nil {
			test.Errorf("%v\n", err)
		}
		// test.Logf("\nmsg: %v\nmsg_s:%v\n", string(msg), string(msg_s))
		msg_v := Verify(pk, msg_s)
		if msg_v == nil {
			test.Errorf("Invalid Signature MEDS-%v\n", 167717)
		}
	}
}

func BenchmarkKeygen9923(b *testing.B) {
	ParameterSetup(9923)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
		b.Elapsed()
	}
}
func BenchmarkKeygen13220(b *testing.B) {
	ParameterSetup(13220)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}
func BenchmarkKeygen41711(b *testing.B) {
	ParameterSetup(41711)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}
func BenchmarkKeygen69497(b *testing.B) {
	ParameterSetup(69497)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}
func BenchmarkKeygen134180(b *testing.B) {
	ParameterSetup(134180)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}
func BenchmarkKeygen167717(b *testing.B) {
	ParameterSetup(167717)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KeyGen()
	}
}

func BenchmarkSign9923(b *testing.B) {
	ParameterSetup(9923)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}
func BenchmarkSign13220(b *testing.B) {
	ParameterSetup(13220)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}
func BenchmarkSign41711(b *testing.B) {
	ParameterSetup(41711)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}
func BenchmarkSign69497(b *testing.B) {
	ParameterSetup(69497)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}
func BenchmarkSign134180(b *testing.B) {
	ParameterSetup(134180)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}
func BenchmarkSign167717(b *testing.B) {
	ParameterSetup(167717)
	_, sk := KeyGen()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sign(sk, msg)
		if err != nil {
			b.Fatal("Error is not nil")
		}
	}
}

func BenchmarkVerify9923(b *testing.B) {
	ParameterSetup(9923)
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
func BenchmarkVerify13220(b *testing.B) {
	ParameterSetup(13220)
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
func BenchmarkVerify41711(b *testing.B) {
	ParameterSetup(41711)
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
func BenchmarkVerify69497(b *testing.B) {
	ParameterSetup(69497)
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
func BenchmarkVerify134180(b *testing.B) {
	ParameterSetup(134180)
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
func BenchmarkVerify167717(b *testing.B) {
	ParameterSetup(167717)
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
