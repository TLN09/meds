package meds

import "testing"

func TestKeyGen(test *testing.T) {
	pk, sk := KeyGen()
	test.Errorf("pk: %v\nsk: %v\n", pk, sk)
}

func TestSign(test *testing.T) {
	_, sk := KeyGen()
	msg := []byte("This is my message")
	signed, err := Sign(sk, msg)
	if err != nil {
		test.Errorf("%v\n", err)
	}
	test.Errorf("%v\nsigned[l_sig:]: %v\n", signed[l_sig-l_salt-l_digest:], signed[l_sig:])
	test.Errorf("%v\n", signed)
}

func TestVerify(test *testing.T) {
	pk, sk := KeyGen()
	msg := []byte("This is my message")
	msg_s, err := Sign(sk, msg)
	if err != nil {
		test.Errorf("%v\n", err)
	}
	// test.Logf("\nmsg: %v\nmsg_s:%v\n", string(msg), string(msg_s))
	msg_v := Verify(pk, msg_s)
	if msg_v == nil {
		test.Errorf("Invalid Signature\n")
		return
	}
}
