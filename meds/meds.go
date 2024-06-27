package meds

import (
	"fmt"
	"math"
	"meds/matrix"

	"golang.org/x/crypto/sha3"
)

var q, q_bitlen, n, m, k, s, t, w int
var l_tree_seed, l_sec_seed, l_pub_seed, l_salt, l_digest int
var l_f_mm, l_f_nn, l_G_i, l_sk, l_pk, l_path, l_sig int

func ParameterSetup(set int) {
	switch set {
	case 9923:
		q = 4093
		q_bitlen = 16
		n = 14
		m = 14
		k = 14
		s = 4
		t = 1152
		w = 14
		l_tree_seed = 16
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	case 13220:
		q = 4093
		q_bitlen = 16
		n = 14
		m = 14
		k = 14
		s = 5
		t = 192
		w = 20
		l_tree_seed = 16
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	case 41711:
		q = 4093
		q_bitlen = 16
		n = 22
		m = 22
		k = 22
		s = 4
		t = 608
		w = 26
		l_tree_seed = 24
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	case 69497:
		q = 4093
		q_bitlen = 16
		n = 22
		m = 22
		k = 22
		s = 5
		t = 160
		w = 36
		l_tree_seed = 24
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	case 134180:
		q = 2039
		q_bitlen = 16
		n = 30
		m = 30
		k = 30
		s = 5
		t = 192
		w = 52
		l_tree_seed = 32
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	case 167717:
		q = 2039
		q_bitlen = 16
		n = 30
		m = 30
		k = 30
		s = 6
		t = 112
		w = 66
		l_tree_seed = 32
		l_sec_seed = 32
		l_pub_seed = 32
		l_salt = 32
		l_digest = 32
		l_f_mm = m * m * (q_bitlen / 8)
		l_f_nn = n * n * (q_bitlen / 8)
		l_G_i = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
		l_sk = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
		l_pk = (s-1)*l_G_i + l_pub_seed
		l_path = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
		l_sig = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt
	default:
		fmt.Printf("Parameter selection error\n")
	}
}

func KeyGen() ([]byte, []byte) {
	delta := Randombytes(l_sec_seed)
	sigma_G_0 := make([]byte, l_pub_seed)
	sigma := make([]byte, l_sec_seed)
	xof := sha3.NewShake256()
	xof.Write(delta)
	xof.Read(sigma_G_0)
	xof.Read(sigma)
	G_0 := ExpandSystMat(sigma_G_0, q, k, m, n)
	sk := make([]byte, l_sk)
	pk := make([]byte, l_pk)

	sk_idx := 0
	pk_idx := 0
	addToKey(sk, delta, &sk_idx)
	addToKey(sk, sigma_G_0, &sk_idx)
	addToKey(pk, sigma_G_0, &pk_idx)

	offset := n * m * Bytelen(q)
	sk_A_idx := sk_idx
	sk_B_idx := sk_idx + (s-1)*offset
	I := matrix.Identity(m, q)
	for i := 1; i < s; i++ {
		var G *matrix.Matrix = nil
		var A_inv *matrix.Matrix
		var A, B_inv *matrix.Matrix = nil, nil
		for G == nil {
			for (A == nil && B_inv == nil) || !Invertable(A, I) || !Invertable(B_inv, I) {
				xof := sha3.NewShake256()
				xof.Write(sigma)
				sigma_a := make([]byte, l_sec_seed)
				sigma_T := make([]byte, l_sec_seed)
				xof.Read(sigma_a)
				xof.Read(sigma_T)
				xof.Read(sigma)
				T_i := ExpandInvMat(sigma_T, q, k)
				a_mm := ExpandFqs(sigma_a, 1, q)[0]
				G_0_prime := T_i.Mul(G_0)
				A, B_inv = Solve(G_0_prime, a_mm, m, n)
			}
			A_inv = Inverse(A)
			B := Inverse(B_inv)
			G = Pi(A, G_0, B)
			G = SF(G)
		}
		addToKey(pk, CompressG(G), &pk_idx)
		addToKey(sk, A_inv.Compress(), &sk_A_idx)
		addToKey(sk, B_inv.Compress(), &sk_B_idx)
	}
	return pk, sk
}

func addToKey(key, bs []byte, idx *int) {
	for i := 0; i < len(bs); i++ {
		key[i+(*idx)] = bs[i]
	}
	(*idx) += len(bs)
}

func Sign(sk, msg []byte) ([]byte, error) {
	f_sk := l_sec_seed
	sigma_G_0 := sk[f_sk : f_sk+l_pub_seed]
	f_sk += l_pub_seed
	G_0 := ExpandSystMat(sigma_G_0, q, k, m, n)
	A_inv := make([]*matrix.Matrix, s-1)
	B_inv := make([]*matrix.Matrix, s-1)
	for i := 0; i < s-1; i++ {
		A_inv[i] = matrix.Decompress(sk[f_sk:f_sk+l_f_mm], m, m, q)
		f_sk += l_f_mm
	}
	for i := 0; i < s-1; i++ {
		B_inv[i] = matrix.Decompress(sk[f_sk:f_sk+l_f_nn], n, n, q)
		f_sk += l_f_nn
	}
	delta := Randombytes(l_sec_seed)
	xof := sha3.NewShake256()
	xof.Write(delta)
	rho := make([]byte, l_tree_seed)
	alpha := make([]byte, l_salt)
	xof.Read(rho)
	xof.Read(alpha)
	seeds, err := SeedTree(rho, alpha, t)
	if err != nil {
		return []byte{}, err
	}
	G_tilde := make([]*matrix.Matrix, t)
	A_tilde := make([]*matrix.Matrix, t)
	B_tilde := make([]*matrix.Matrix, t)
	sigma_prime := make([]byte, l_salt+l_tree_seed+4)
	for i := 0; i < t; i++ {
		for G_tilde[i] == nil {
			idx := 0
			for j := 0; j < l_salt; j++ {
				sigma_prime[idx] = alpha[j]
				idx++
			}
			for j := 0; j < l_tree_seed; j++ {
				sigma_prime[idx] = seeds[i][j]
				idx++
			}
			x, err := ToBytes(int32(math.Pow(2, float64(1+int(math.Ceil(math.Log2(float64(t)))))))+int32(i), 4)
			if err != nil {
				return []byte{}, err
			}
			for j := 0; j < 4; j++ {
				sigma_prime[idx] = x[j]
				idx++
			}
			sigma_A_tilde := make([]byte, l_pub_seed)
			sigma_B_tilde := make([]byte, l_pub_seed)
			xof = sha3.NewShake256()
			xof.Write(sigma_prime)
			xof.Read(sigma_A_tilde)
			xof.Read(sigma_B_tilde)
			xof.Read(seeds[i])
			A_tilde[i] = ExpandInvMat(sigma_A_tilde, q, m)
			B_tilde[i] = ExpandInvMat(sigma_B_tilde, q, n)
			G_tilde[i] = Pi(A_tilde[i], G_0, B_tilde[i])
			G_tilde[i] = SF(G_tilde[i])
		}
	}
	H := sha3.NewShake256()
	for i := 0; i < t; i++ {
		H.Write(G_tilde[i].Submatrix(0, G_tilde[i].M, k, m*n).Compress())
	}
	H.Write(msg)
	d := make([]byte, l_digest)
	H.Read(d)

	h := ParseHash(s, t, w, d)
	msg_s := make([]byte, l_sig+len(msg))
	idx := 0
	for i := 0; i < t; i++ {
		if h[i] > 0 {
			mu := A_tilde[i].Mul(A_inv[h[i]-1])
			nu := B_inv[h[i]-1].Mul(B_tilde[i])
			for _, b := range mu.Compress() {
				msg_s[idx] = b
				idx++
			}
			for _, b := range nu.Compress() {
				msg_s[idx] = b
				idx++
			}
		}
	}

	p := SeedTreeToPath(w, t, h, rho, alpha)
	for i := 0; i < len(p); i++ {
		msg_s[idx] = p[i]
		idx++
	}
	for i := 0; i < len(d); i++ {
		msg_s[idx] = d[i]
		idx++
	}
	for i := 0; i < len(alpha); i++ {
		msg_s[idx] = alpha[i]
		idx++
	}
	for i := 0; i < len(msg); i++ {
		msg_s[idx] = msg[i]
		idx++
	}

	return msg_s, nil
}

func Verify(pk, msg_s []byte) []byte {
	sigma_G_0 := pk[:l_pub_seed]
	G_0 := ExpandSystMat(sigma_G_0, q, k, m, n)
	f_pk := l_pub_seed
	G := make([]*matrix.Matrix, s-1)
	for i := 0; i < s-1; i++ {
		G[i] = DecompressG(pk[f_pk:f_pk+l_G_i], q, m, n, k)
		f_pk += l_G_i
	}

	p := msg_s[l_sig-l_digest-l_salt-l_path : l_sig-l_digest-l_salt]
	d := msg_s[l_sig-l_digest-l_salt : l_sig-l_salt]
	alpha := msg_s[l_sig-l_salt : l_sig]
	msg := msg_s[l_sig:]
	h := ParseHash(s, t, w, d)
	seeds := PathToSeedTree(h, p, alpha, l_tree_seed)
	f_msg_s := 0
	I := matrix.Identity(m, q)
	G_hat := make([]*matrix.Matrix, t)
	for i := 0; i < t; i++ {
		if h[i] > 0 {
			mu := matrix.Decompress(msg_s[f_msg_s:f_msg_s+l_f_mm], m, m, q)
			nu := matrix.Decompress(msg_s[f_msg_s+l_f_mm:f_msg_s+l_f_mm+l_f_nn], n, n, q)
			f_msg_s += l_f_mm + l_f_nn
			if !Invertable(mu, I) || !Invertable(nu, I) {
				fmt.Print("Mu or Nu not invertable\n")
				return nil
			}
			G_hat[i] = Pi(mu, G[h[i]-1], nu)
			err := SF_on_submatrix(G_hat[i], 0, 0, G_hat[i].M, G_hat[i].N)
			if err != nil {
				fmt.Print("SF failed on G_hat\n")
				return nil
			}
		} else {
			for G_hat[i] == nil {
				sigma_prime := make([]byte, l_salt+l_tree_seed+4)
				sigma_A := make([]byte, l_pub_seed)
				sigma_B := make([]byte, l_pub_seed)
				x, err := ToBytes(int32(math.Pow(2, float64(1+int(math.Ceil(math.Log2(float64(t)))))))+int32(i), 4)
				if err != nil {
					fmt.Print("Failed getting byte value\n")
					return nil
				}
				idx := 0
				for j := 0; j < l_salt; j++ {
					sigma_prime[idx] = alpha[j]
					idx++
				}
				for j := 0; j < l_tree_seed; j++ {
					sigma_prime[idx] = seeds[i][j]
					idx++
				}
				for j := 0; j < 4; j++ {
					sigma_prime[idx] = x[j]
					idx++
				}
				xof := sha3.NewShake256()
				xof.Write(sigma_prime)
				xof.Read(sigma_A)
				xof.Read(sigma_B)
				xof.Read(seeds[i])
				A_hat := ExpandInvMat(sigma_A, q, m)
				B_hat := ExpandInvMat(sigma_B, q, n)
				G_hat[i] = Pi(A_hat, G_0, B_hat)
				G_hat[i] = SF(G_hat[i])
			}
		}
	}
	d_prime := make([]byte, l_digest)
	H := sha3.NewShake256()
	for i := 0; i < t; i++ {
		H.Write(G_hat[i].Submatrix(0, G_hat[i].M, k, m*n).Compress())
	}
	H.Write(msg)
	H.Read(d_prime)
	equal := true
	for i := 0; equal && i < l_digest; i++ {
		equal = d[i] == d_prime[i]
	}
	if equal {
		return msg
	}

	fmt.Print("Signature not valid\n")
	return nil
}
