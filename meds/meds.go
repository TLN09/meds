package meds

import (
	"math"
	"meds/matrix"

	"golang.org/x/crypto/sha3"
)

const q int = 4093
const q_bitlen int = 16
const n int = 14
const m int = 14
const k int = 14
const s int = 4
const t int = 1152
const w int = 14
const l_tree_seed int = 16
const l_sec_seed int = 32
const l_pub_seed int = 32
const l_salt int = 32
const l_digest int = 32
const l_f_mm int = m * m * (q_bitlen / 8)
const l_f_nn int = n * n * (q_bitlen / 8)
const l_G_i int = (((k-2)*(m*n-k) + n) * q_bitlen) / 8
const l_sk int = (s-1)*(l_f_mm+l_f_nn) + l_sec_seed + l_pub_seed
const l_pk int = (s-1)*l_G_i + l_pub_seed

var l_path int = (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * l_tree_seed
var l_sig int = l_digest + w*(l_f_mm+l_f_nn) + l_path + l_salt

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
		addToKey(sk, A_inv.Compress(), &sk_idx)
		addToKey(sk, B_inv.Compress(), &sk_idx)
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
		// fmt.Printf("%v\n", len(sk[f_sk:f_sk+l_f_mm]))
		A_inv[i] = matrix.Decompress(sk[f_sk:f_sk+l_f_mm], m, m, q)
		f_sk += l_f_mm
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
	// fmt.Printf("alpha: %v\n", alpha)
	seeds, err := SeedTree(rho, alpha, t)
	// fmt.Printf("seeds: %v\n", seeds)
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
			xof.Read(sigma_prime)
			A_tilde[i] = ExpandInvMat(sigma_A_tilde, q, m)
			B_tilde[i] = ExpandInvMat(sigma_B_tilde, q, n)
			G_tilde[i] = Pi(A_tilde[i], G_0, B_tilde[i])
			G_tilde[i] = SF(G_tilde[i])
		}
	}
	H := sha3.NewShake256()
	for i := 0; i < t; i++ {
		H.Write(G_tilde[i].Submatrix(0, G_tilde[i].M, k, m*n).Compress())
		// fmt.Printf("G_tilde[i]: %v\n", G_tilde[i].Submatrix(0, G_tilde[i].M, k, m*n).Compress())
	}
	H.Write(msg)
	d := make([]byte, l_digest)
	H.Read(d)
	// fmt.Printf("d: %v\n", d)

	h := ParseHash(s, t, w, d)
	msg_s := make([]byte, l_sig+len(msg))
	idx := 0
	for i := 0; i < t; i++ {
		if h[i] > 0 {
			// fmt.Printf("h[i]: %v\nÃ: %v\nA^-1: %v\nB^-1: %v\nB~: %v\n", h[i], len(A_tilde), len(A_inv), len(B_inv), len(B_tilde))
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

	// fmt.Printf("%v\n", len(v) == 2*w)

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

	// fmt.Printf("%v\n", msg_s)

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
	// fmt.Printf("alpha: %v\nd: %v\n", alpha, d)
	msg := msg_s[l_sig:]
	// fmt.Printf("msg: %v\n", string(msg))
	h := ParseHash(s, t, w, d)
	seeds := PathToSeedTree(h, p, alpha, l_tree_seed)
	// fmt.Printf("seeds: %v\n", seeds)
	f_msg_s := 0
	I := matrix.Identity(m, q)
	G_hat := make([]*matrix.Matrix, t)
	for i := 0; i < t; i++ {
		if h[i] > 0 {
			mu := matrix.Decompress(msg_s[f_msg_s:f_msg_s+l_f_mm], m, m, q)
			nu := matrix.Decompress(msg_s[f_msg_s+l_f_mm:f_msg_s+l_f_mm+l_f_nn], n, n, q)
			f_msg_s += l_f_mm + l_f_nn
			if !Invertable(mu, I) || !Invertable(nu, I) {
				return nil
			}
			G_hat[i] = Pi(mu, G[h[i]-1], nu)
			G_hat[i] = SF(G_hat[i])
			if G_hat[i] == nil {
				return nil
			}
		} else {
		LINE_24:
			sigma_prime := make([]byte, l_salt+l_tree_seed+4)
			sigma_A := make([]byte, l_pub_seed)
			sigma_B := make([]byte, l_pub_seed)
			x, err := ToBytes(int32(math.Pow(2, float64(1+int(math.Ceil(math.Log2(float64(t)))))))+int32(i), 4)
			if err != nil {
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
			xof.Read(sigma_prime)
			A_hat := ExpandInvMat(sigma_A, q, m)
			B_hat := ExpandInvMat(sigma_B, q, n)
			G_hat[i] = Pi(A_hat, G_0, B_hat)
			G_hat[i] = SF(G_hat[i])
			if G_hat[i] == nil {
				goto LINE_24
			}
		}
	}
	d_prime := make([]byte, l_digest)
	H := sha3.NewShake256()
	for i := 0; i < t; i++ {
		H.Write(G_hat[i].Submatrix(0, G_hat[i].M, k, m*n).Compress())
		// fmt.Printf("G_hat[i]: %v\n", G_hat[i].Submatrix(0, G_hat[i].M, k, m*n).Compress())
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
	// fmt.Printf("d: %v\nd_prime: %v", d, d_prime)
	return nil
}
