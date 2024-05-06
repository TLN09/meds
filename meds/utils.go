package meds

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"meds/finiteField"
	"meds/matrix"
	"meds/seedTree"

	"golang.org/x/crypto/sha3"
)

func Randombytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// CompressG compresses the given matrix.Matrix to a []byte
// Returns: []byte with length as specified in MEDS specification document
func CompressG(G *matrix.Matrix) []byte {
	l_g_prime := (k-2)*(G.N-k) + n
	G_prime := matrix.New(1, l_g_prime, G.Q)
	idx := 0
	for i := 0; i < n; i++ {
		G_prime.Set(0, idx, G.Get(1, m*n-n+i))
		idx++
	}
	for i := 2; i < k; i++ {
		for j := k; j < m*n; j++ {
			G_prime.Set(0, idx, G.Get(i, j))
			idx++
		}
	}

	return G_prime.Compress()
}

// DecompressG Decompresses the []byte to a matrix.Matrix
// Returns: *matrix.Matrix
func DecompressG(b []byte, q int, m, n, k int) *matrix.Matrix {
	l_g_prime := (k-2)*(m*n-k) + n
	G_prime := matrix.Decompress(b, 1, l_g_prime, q)
	G := matrix.New(k, m*n, q)
	// Making the first k by k submatrix the identity matrix
	for i := 0; i < G.M; i++ {
		G.Get(i, i).Set(1)
	}

	for i := 1; i < m; i++ {
		G.Get(0, i*(n+1)).Set(1)
	}
	for i := 1; i < m-1; i++ {
		G.Get(1, i*(n+1)+1).Set(1)
	}
	idx := 0
	for i := 0; i < n; i++ {
		G.Set(1, m*n-n+i, G_prime.Get(0, idx))
		idx++
	}
	for i := 2; i < k; i++ {
		for j := k; j < m*n; j++ {
			G.Set(i, j, G_prime.Get(0, idx))
			idx++
		}
	}

	return G
}

// ExpandFqs generates field elements from the given seed
// Returns: An array of field elements
func ExpandFqs(seed []byte, l, q int) []*finiteField.Fq {
	shake := sha3.NewShake256()
	shake.Write(seed)
	a := make([]*finiteField.Fq, l)
	for i := 0; i < l; i++ {
		a[i] = expandFqs(shake, q)
	}

	return a
}

func expandFqs(shake sha3.ShakeHash, q int) *finiteField.Fq {
	m := int(math.Pow(2, float64(Bitlen(q))))
	byte_len := Bytelen(q)
	buf := make([]byte, 1)
EXPAND_LOOP:
	a := 0
	for j := 0; j < byte_len; j++ {
		shake.Read(buf)
		a += int(buf[0]) * int(math.Pow(2, float64(8*j)))
	}
	a %= m
	if a >= q {
		goto EXPAND_LOOP
	}
	return finiteField.NewFieldElm(a, q)
}

// ExpandSystMat generates a matrix in systematic form from the given seed
// Returns: Matrix $M \in F_{q}^{k \times mn}$
func ExpandSystMat(seed []byte, q, k, m, n int) *matrix.Matrix {
	a := ExpandFqs(seed, k*(m*n-k), q)
	f_a := 0
	M := matrix.New(k, m*n, q)

	for i := 0; i < k; i++ {
		M.Get(i, i).Set(1)
		for j := k; j < m*n; j++ {
			M.Set(i, j, a[f_a])
			f_a++
		}
	}

	return M
}

// RowsToMatricies takes a matrix $M \in F_{q}^{k \times mn}$
// as input and returns an array of matricies $P_1, \dots, P_{k-1} \in F_{q}^{m \times n}$
func RowsToMatricies(M *matrix.Matrix, m, n int) []*matrix.Matrix {
	P := make([]*matrix.Matrix, M.M)

	for i := 0; i < M.M; i++ {
		P[i] = matrix.New(m, n, M.Q)
		for j := 0; j < M.N; j++ {
			P[i].Set(j/n, j%n, M.Get(i, j))
		}
	}

	return P
}

func MatriciesToRows(P []*matrix.Matrix) *matrix.Matrix {
	k := len(P)
	M := matrix.New(k, P[0].M*P[0].N, P[0].Q)

	for i := 0; i < k; i++ {
		for j := 0; j < M.N; j++ {
			M.Set(i, j, P[i].Get(j/P[i].N, j%P[i].N))
		}
	}

	return M
}

func Pi(A, G, B *matrix.Matrix) *matrix.Matrix {
	P := RowsToMatricies(G, A.M, B.N)
	P_prime := make([]*matrix.Matrix, G.M)

	for i := 0; i < G.M; i++ {
		P_prime[i] = A.Mul(P[i]).Mul(B)
	}

	return MatriciesToRows(P_prime)
}

func ExpandInvMat(seed []byte, q, d int) *matrix.Matrix {
	shake := sha3.NewShake256()
	shake.Write(seed)
	M := matrix.New(d, d, q)
	I := matrix.Identity(d, q)
INVERTABLE_LOOP:
	for i := 0; i < d; i++ {
		for j := 0; j < d; j++ {
			M.Set(i, j, expandFqs(shake, q))
		}
	}
	M_sf := SF(M)
	if M_sf == nil || !M_sf.Equals(I) {
		goto INVERTABLE_LOOP
	}

	return M
}

func Bitlen(x int) int {
	return int(math.Ceil(math.Log2(float64(x))))
}

func Bytelen(x int) int {
	return int(math.Ceil(math.Log2(float64(x)) / float64(8)))
}

func ParseHash(s, t, w int, d []byte) []byte {
	t_bitlen := Bitlen(t)
	t_bytelen := Bytelen(t)
	s_bitlen := Bitlen(s)
	xof := sha3.NewShake256()
	xof.Write(d)
	h := make([]byte, t)
	buf := make([]byte, 1)
	xof.Read(buf)

	for i := 0; i < w; i++ {
	LINE_6:
		f_h := 0
		for j := 0; j < t_bytelen; j++ {
			f_h += int(buf[0]) * int(math.Pow(float64(2), float64(8*j)))
			xof.Read(buf)
		}
		f_h %= int(math.Pow(float64(2), float64(t_bitlen)))
		if f_h >= t || h[f_h] > 0 {
			goto LINE_6
		}
	LINE_13:
		h[f_h] = buf[0]
		xof.Read(buf)
		h[f_h] = byte(int(h[f_h]) % int(math.Pow(float64(2), float64(s_bitlen))))
		if h[f_h] == 0 || int(h[f_h]) >= s {
			goto LINE_13
		}
	}
	return h
}

func swapElm(M *matrix.Matrix, i, j, k int) {
	tmp := M.Get(i, k)
	M.Set(i, k, M.Get(j, k))
	M.Set(j, k, tmp)
}

func multFixedConst(M *matrix.Matrix, i int, c *finiteField.Fq) {
	for j := 0; j < M.N; j++ {
		M.Set(i, j, M.Get(i, j).Mul(c))
	}
}

func swapRows(M *matrix.Matrix, i, j int) {
	for k := 0; k < M.N; k++ {
		swapElm(M, i, j, k)
	}
}

func constTimesEq1PlusEq2(M *matrix.Matrix, c *finiteField.Fq, i, j int) {
	for k := 0; k < M.N; k++ {
		M.Set(j, k, M.Get(j, k).Add(M.Get(i, k).Mul(c)))
	}
}

func zeroCol(M *matrix.Matrix, j int) bool {
	zero_col := true
	zero := finiteField.NewFieldElm(0, M.Q)

	for i := 0; i < M.M && zero_col; i++ {
		zero_col = M.Get(i, j).Equals(zero)
	}

	return zero_col
}

func SF(M *matrix.Matrix) *matrix.Matrix {
	sf := matrix.New(M.M, M.N, M.Q)
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			sf.Set(i, j, finiteField.NewFromBytes(M.Get(i, j).Bytes(), M.Q))
		}
	}
	zero := finiteField.NewFieldElm(0, M.Q)

	for i := 0; i < sf.M; i++ {
		// fmt.Printf("%v\n", sf)
		l := i
		if zeroCol(sf, l) {
			return nil
		}
		for k := i + 1; k < sf.M && sf.Get(i, l).Equals(zero); k++ {
			swapRows(sf, i, k)
		}
		multFixedConst(sf, i, sf.Get(i, l).Inv())
		for k := 0; k < sf.M; k++ {
			if k == i {
				continue
			}
			c := sf.Get(k, l).UnaryMinus()
			constTimesEq1PlusEq2(sf, c, i, k)
		}
	}

	return sf
}

func augment(M *matrix.Matrix) *matrix.Matrix {
	aug := matrix.New(M.M, 2*M.N, M.Q)
	for i := 0; i < M.M; i++ {
		aug.Set(i, i+M.N, finiteField.NewFieldElm(1, M.Q))
	}

	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			aug.Set(i, j, M.Get(i, j))
		}
	}
	return aug
}

func Inverse(M *matrix.Matrix) *matrix.Matrix {
	aug := augment(M)
	sf := SF(aug)
	inv := matrix.New(M.M, M.N, M.Q)
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			inv.Set(i, j, sf.Get(i, j+M.N))
		}
	}
	return inv
}

func Solve(G_prime *matrix.Matrix, a *finiteField.Fq, m, n int) (*matrix.Matrix, *matrix.Matrix) {
	P0_prime := RowsToMatricies(G_prime, m, n)[:2]
	// fmt.Printf("(m, n): (%v, %v)\n", m, n)
	P := make([]*matrix.Matrix, 2)
	P[0] = matrix.Identity(m, P0_prime[0].Q)
	P[1] = matrix.UpperShift(m, P0_prime[0].Q)

	rsys := matrix.New(m*m+n*n-1, m*m+n*n, P0_prime[0].Q)
	fill_rsys(rsys, P0_prime, P, a)
	// fmt.Printf("%v", strings.ReplaceAll(fmt.Sprintf("rsys: %v\n", rsys), "0", " "))
	err := solve_sub_matricies(rsys, m, n)
	if err != nil {
		// fmt.Printf("%v\n", err)
		// fmt.Printf("%v\n", rsys)
		return nil, nil
	}
	// fmt.Printf("%v", strings.ReplaceAll(fmt.Sprintf("rsys: %v\n", rsys), "0", " "))
	backprop_to_sf(rsys, m)
	// sf_on_submatrix(rsys, 0, 0, rsys.M, rsys.N)
	// fmt.Printf("%v", strings.ReplaceAll(fmt.Sprintf("rsys_rref: %v\n", rsys), "0", " "))

	values := make([]*finiteField.Fq, m*m+n*n)
	for i := 0; i < rsys.M; i++ {
		values[i] = rsys.Get(i, rsys.N-1)
	}
	values[len(values)-1] = a

	// fmt.Printf("values: %v\n", values)

	B_inv := array_to_matrix(m, values[:n*n], rsys.Q)
	A := array_to_matrix(m, values[n*n:], rsys.Q)

	return A, B_inv
}

func array_to_matrix(m int, values []*finiteField.Fq, q int) *matrix.Matrix {
	M := matrix.New(m, m, q)
	for i := 0; i < len(values); i++ {
		M.Set(i/m, i%m, values[i])
	}
	return M
}

func fill_rsys(rsys *matrix.Matrix, P0_prime []*matrix.Matrix, P []*matrix.Matrix, a *finiteField.Fq) {
	eqs1_A_coeff := P0_prime[0].Transpose().UnaryMinus()
	eqs2_A_coeff := P0_prime[1].Transpose().UnaryMinus()
	// fmt.Printf("P[1]: %v\n", P[1])
	// fmt.Printf("A_coefficients: %v\n", eqs1_A_coeff)

	// Fill in coefficients of A
	row := 0
	col := rsys.N / 2
	for col < (rsys.N - eqs1_A_coeff.N) {
		for i := 0; i < eqs1_A_coeff.M; i++ {
			for j := 0; j < eqs1_A_coeff.N; j++ {
				rsys.Set(row+i, col+j, eqs1_A_coeff.Get(i, j))
			}
		}
		row += eqs1_A_coeff.M
		col += eqs1_A_coeff.N
	}

	// Special case for last column of last matrix in eqs1
	for i := 0; i < eqs1_A_coeff.M; i++ {
		for j := 0; j < eqs1_A_coeff.N-1; j++ {
			rsys.Set(row+i, col+j, eqs1_A_coeff.Get(i, j))
		}
		rsys.Set(row+i, col+eqs1_A_coeff.N-1, eqs1_A_coeff.Get(i, eqs1_A_coeff.N-1).Mul(a).UnaryMinus())
	}

	row = P[0].M * P[0].N
	col = rsys.N / 2
	for col < (rsys.N - eqs2_A_coeff.N) {
		for i := 0; i < eqs2_A_coeff.M && (row+i) < rsys.M; i++ {
			for j := 0; j < eqs2_A_coeff.N && (col+j) < rsys.N; j++ {
				rsys.Set(row+i, col+j, eqs2_A_coeff.Get(i, j))
			}
			// Insert -eqs1_A_coeff after eqs2_A_coeff since the U_k matrix of B coefficients needs to be eliminated.
			for j := eqs2_A_coeff.N; j < eqs2_A_coeff.N+eqs1_A_coeff.N; j++ {
				if col+j == rsys.N-1 { // Special case for last column as this is different from the rest of the coefficients
					rsys.Set(row+i, col+j, eqs1_A_coeff.Get(i, j-eqs2_A_coeff.N).Mul(a))
				} else {
					rsys.Set(row+i, col+j, eqs1_A_coeff.Get(i, j-eqs2_A_coeff.N).UnaryMinus())
				}
			}
		}
		row += eqs2_A_coeff.M
		col += eqs2_A_coeff.N
	}

	// Special case for last column of last matrix in eqs2
	for i := 0; i < eqs2_A_coeff.M && (row+i) < rsys.M; i++ {
		for j := 0; j < eqs2_A_coeff.N-1 && (col+j) < rsys.N; j++ {
			rsys.Set(row+i, col+j, eqs2_A_coeff.Get(i, j))
		}
		rsys.Set(row+i, col+eqs2_A_coeff.N-1, eqs2_A_coeff.Get(i, eqs2_A_coeff.N-1).Mul(a).UnaryMinus())
	}

	// Filling in coefficients of B
	// Only needing the first part with B * I_m as the other part gets eliminated by inserting -eqs1_A_coeff after eqs2_A_coeff
	row = 0
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			rsys.Set(row, row, finiteField.NewFieldElm(1, rsys.Q))
			row++
		}
	}

	// col = P[1].N
	// for i := 1; i < P[1].M; i++ {
	// 	for j := 0; j < P[1].N; j++ {
	// 		rsys.Set(row, col, finiteField.NewFieldElm(1, rsys.Q))
	// 		row++
	// 		col++
	// 	}
	// }
}

func solve_sub_matricies(rsys *matrix.Matrix, m, n int) error {
	base_row := m * m
	base_col := n * n
	// fmt.Printf("(row, col, m, n): (%v, %v, %v, %v)\n", base_row, base_col, m, 2*n)
	err := sf_on_submatrix(rsys, base_row, base_col, m, 2*n)
	if err != nil {
		return err
	}

	row := base_row + m
	col := base_col + n
	for row < (rsys.M - 2*m) {
		for i := 0; i < m; i++ {
			for j := 0; j < 2*n; j++ {
				rsys.Set(row+i, col+j, rsys.Get(base_row+i, base_col+j))
			}
		}
		row += m
		col += n
	}
	// fmt.Printf("(row, col, m, n): (%v, %v, %v, %v)\n", row, col, m, 2*n)
	err = sf_on_submatrix(rsys, row, col, m, 2*n)
	if err != nil {
		return err
	}
	row += m
	col += n

	// fmt.Printf("(row, col, m, n): (%v, %v, %v, %v)\n", row, col, m-1, n)
	err = sf_on_submatrix(rsys, row, col, m-1, n)
	if err != nil {
		return err
	}

	return nil
}

func sf_on_submatrix(M *matrix.Matrix, row, col, m, n int) error {
	zero := finiteField.NewFieldElm(0, M.Q)

	failed := false
	for i := 0; i < m; i++ {
		// fmt.Printf("%v\n", M)
		l := i
		if zero_col_submatrix(M, l, row, col, m, zero) {
			return fmt.Errorf("includes zero col %v", l)
		}
		for k := 1 + 1; k < m && M.Get(row+i, col+l).Equals(zero); k++ {
			swap_rows_submatrix(M, i, k, row, col, n)
		}
		mult_fixed_const_submatrix(M, i, M.Get(row+i, col+l).Inv(), row, col, n)
		for k := 0; k < m; k++ {
			if k == i {
				continue
			}
			c := M.Get(row+k, col+l).UnaryMinus()
			const_times_eq1_plus_eq2_submatrix(M, c, i, k, row, col, n)
		}
	}
	if failed {
		// fmt.Printf("Failed\n")
		return errors.New("failed sf of submatrix")
	}
	return nil
}

func mult_fixed_const_submatrix(M *matrix.Matrix, i int, c *finiteField.Fq, row, col, n int) {
	for j := 0; j < n; j++ {
		M.Set(row+i, col+j, M.Get(row+i, col+j).Mul(c))
	}
}

func swap_rows_submatrix(M *matrix.Matrix, i, j, row, col, n int) {
	for k := 0; k < n; k++ {
		swapElm(M, row+i, row+j, col+k)
	}
}

func const_times_eq1_plus_eq2_submatrix(M *matrix.Matrix, c *finiteField.Fq, i, j, row, col, n int) {
	for k := 0; k < n; k++ {
		M.Set(row+j, col+k, M.Get(row+j, col+k).Add(M.Get(row+i, col+k).Mul(c)))
	}
}

func zero_col_submatrix(M *matrix.Matrix, j, row, col, m int, zero *finiteField.Fq) bool {
	zero_col := true

	for i := 0; i < m && zero_col; i++ {
		zero_col = M.Get(row+i, col+j).Equals(zero)
		// fmt.Printf("i: %v\n", i)
	}

	return zero_col
}

func backprop_to_sf(M *matrix.Matrix, m int) {
	zero := finiteField.NewFieldElm(0, M.Q)
	col := M.N - 2
	for row := M.M - 1; row >= m*m; row-- {
		for i := 0; i < row; i++ {
			if M.Get(i, col).Equals(zero) {
				continue
			}
			// fmt.Printf("row, col, i: %v, %v, %v\n", row, col, i)
			// fmt.Printf("M[i, col]: %v\n", M.Get(i, col))
			c := M.Get(i, col).UnaryMinus()
			// fmt.Printf("c: %v\n", c)
			constTimesEq1PlusEq2(M, c, row, i)
			// const_times_eq1_plus_eq2_submatrix(M, c, row, i, 0, col, M.N-col)
		}
		col--
	}
}

func SeedTree(seed, salt []byte, t int) ([][]byte, error) {
	treeHeight := int(math.Ceil(math.Log2(float64(t))))
	root := seedTree.New(seed, salt, 0, 0, nil, nil, nil)
	leafs := make([]*seedTree.SeedTreeNode, t)
	err := root.CreateSeedTree(treeHeight, &leafs)
	// fmt.Printf("seed tree: %v\n", root)
	if err != nil {
		return [][]byte{}, err
	}

	seeds := make([][]byte, t)
	for i := 0; i < t; i++ {
		seeds[i] = leafs[i].Seed()
	}

	return seeds, nil
}

func SeedTreeToPath(w, t int, digest, seed, salt []byte) []byte {
	l_path := (int(math.Pow(2, math.Ceil(math.Log2(float64(w))))) + w*(int(math.Ceil(math.Log2(float64(t))))-int(math.Ceil(math.Log2(float64(w))))-1)) * len(seed)
	path := make([]byte, l_path)
	treeHeight := int(math.Ceil(math.Log2(float64(t))))
	root := seedTree.New(seed, salt, 0, 0, nil, nil, nil)
	leafs := make([]*seedTree.SeedTreeNode, t)
	root.CreateSeedTree(treeHeight, &leafs)
	for i := 0; i < t; i++ {
		if digest[i] != byte(0) {
			leafs[i].RemoveSeedLabel()
		}
	}
	idx := 0
	root.SeedTreeToPath(&path, &idx)
	return path
}

func PathToSeedTree(digest, path, salt []byte, l_tree_seed int) [][]byte {
	t := len(digest)
	treeHeight := int(math.Ceil(math.Log2(float64(t))))
	leafs := make([]*seedTree.SeedTreeNode, t)
	root := seedTree.EmptyTree(treeHeight, salt, &leafs)
	for i := 0; i < t; i++ {
		if digest[i] != byte(0) {
			leafs[i].RemoveSeedLabel()
		}
	}
	seeds := make([][]byte, t)
	for i := 0; i < t; i++ {
		seeds[i] = make([]byte, l_tree_seed)
	}
	seeds_idx := 0
	path_idx := 0
	root.PathToSeedTree(&path, l_tree_seed, &seeds, &seeds_idx, &path_idx)
	return seeds
}

func ToBytes(a int32, l int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, a)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes()[:l], nil
}

func Invertable(M, I *matrix.Matrix) bool {
	M_sf := SF(M)
	if M_sf == nil {
		return false
	}
	return M_sf.Equals(I)
}
