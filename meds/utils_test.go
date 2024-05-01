package meds

import (
	"fmt"
	"math"
	"math/rand"
	"meds/finiteField"
	"meds/matrix"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestCompressG(t *testing.T) {
	k := 14
	m := 14
	n := 14
	G := matrix.New(k, m*n, q)
	E := matrix.New(k, m*n, q)
	for i := 0; i < G.M; i++ {
		for j := 0; j < G.N; j++ {
			n := rand.Intn(int(q))
			G.Set(i, j, finiteField.NewFieldElm(n, q))
			if i == 1 && j >= G.N-k {
				E.Set(i, j, G.Get(i, j))

			} else if i > 1 {
				E.Set(i, j, G.Get(i, j))
			}
		}
	}

	for i := 0; i < G.M; i++ {
		for j := 0; j < G.M; j++ {
			if i == j {
				G.Get(i, j).Set(1)
				E.Get(i, j).Set(1)
			} else {
				G.Get(i, j).Set(0)
				E.Get(i, j).Set(0)
			}
		}
	}
	for i := 1; i < m; i++ {
		G.Get(0, i*(n+1)).Set(1)
		E.Get(0, i*(n+1)).Set(1)
	}
	for i := 1; i < m-1; i++ {
		G.Get(0, i*(n+1)+1).Set(1)
		E.Get(0, i*(n+1)+1).Set(1)
	}

	R := DecompressG(CompressG(G), q, m, n, k)
	if !R.Equals(E) {
		t.Errorf("Compressed then decompressed is not equal to itself\nE:%v\nR:%v", E, R)
		return
	}
}

func TestSF(t *testing.T) {
	for k := 2; k < 21; k++ {
		m := k
		n := k
		G := matrix.New(k, m*n, q)
		for i := 0; i < G.M; i++ {
			for j := 0; j < G.N; j++ {
				n := rand.Intn(int(q))
				G.Set(i, j, finiteField.NewFieldElm(n, q))
			}
		}
		t.Logf("G: %v", G)

		// Write A and B to a file
		err := os.WriteFile("SF_test.txt", G.Compress(), 0666)
		if err != nil {
			t.Errorf("Unable to create file SF_test.txt: %v", err)
			return
		}

		// Compute using SageMath the correct matrix product and save to a file
		cmd := exec.Command("sage", "utils_test.sage", "sf", strconv.Itoa(int(G.Q)), strconv.Itoa(G.M), strconv.Itoa(G.N))
		output, err := cmd.Output()
		if err != nil {
			t.Errorf("Unable to create the E matrix: %v", err)
			return
		}

		var b strings.Builder
		b.Write(output)
		t.Logf("sagemath output: \n%v", b.String())
		// Read E result and parse into a matrix
		file_content, err := os.ReadFile("E_test.txt")
		if err != nil {
			t.Errorf("Unable to read E matrix file: %v", err)
			return
		}

		// Check E == AB
		E := matrix.Decompress(file_content, G.M, G.N, G.Q)

		R := SF(G)
		if !R.Equals(E) {
			t.Errorf("\nSF(G): %v\nE: %v\n", R, E)
		}
	}
}

func TestSolve(t *testing.T) {
	// q = 7 test case
	k := 4
	m := k
	n := k
	p0 := [][]int{
		{1, 1, 4, 0},
		{0, 3, 0, 4},
		{0, 1, 0, 6},
		{0, 2, 4, 6},
	}
	p1 := [][]int{
		{0, 6, 2, 4},
		{1, 1, 6, 0},
		{0, 0, 3, 2},
		{0, 3, 5, 1},
	}
	e_a := [][]int{
		{6, 1, 1, 3},
		{4, 3, 5, 1},
		{0, 6, 3, 3},
		{5, 3, 1, 2},
	}
	e_b_inv := [][]int{
		{4, 1, 5, 3},
		{6, 6, 2, 4},
		{4, 2, 3, 0},
		{5, 3, 1, 1},
	}

	P := make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 7)
	P[1] = matrix.New(m, n, 7)
	E_A := matrix.New(m, n, 7)
	E_B_inv := matrix.New(m, n, 7)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 7))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 7))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 7))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 7))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G := MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a := finiteField.NewFieldElm(2, 7)
	t.Logf("(k, m, n): (%v, %v, %v)\n", k, m, n)
	t.Logf("q: %v\n", 7)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv := Solve(G, a, m, n)
	if A == nil || B_inv == nil || !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 11 test

	p0 = [][]int{
		{1, 4, 4, 2},
		{0, 7, 7, 5},
		{0, 3, 6, 10},
		{0, 9, 0, 0},
	}
	p1 = [][]int{
		{0, 9, 6, 10},
		{1, 4, 8, 6},
		{0, 8, 10, 7},
		{0, 8, 9, 9},
	}
	e_a = [][]int{
		{2, 10, 8, 0},
		{3, 1, 1, 3},
		{3, 0, 2, 6},
		{9, 8, 6, 9},
	}
	e_b_inv = [][]int{
		{8, 5, 1, 2},
		{6, 7, 6, 9},
		{1, 0, 6, 0},
		{6, 0, 7, 6},
	}

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 11)
	P[1] = matrix.New(m, n, 11)
	E_A = matrix.New(m, n, 11)
	E_B_inv = matrix.New(m, n, 11)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 11))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 11))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 11))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 11))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a = finiteField.NewFieldElm(-2, 11)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 13 test
	p0 = [][]int{
		{1, 1, 1, 5},
		{0, 12, 2, 5},
		{0, 4, 4, 0},
		{0, 9, 12, 5},
	}
	p1 = [][]int{
		{0, 6, 8, 0},
		{1, 12, 5, 8},
		{0, 9, 5, 9},
		{0, 5, 1, 2},
	}
	e_a = [][]int{
		{5, 12, 2, 7},
		{11, 7, 9, 7},
		{0, 1, 6, 11},
		{6, 8, 7, 4},
	}
	e_b_inv = [][]int{
		{2, 1, 4, 11},
		{10, 7, 12, 11},
		{10, 1, 2, 6},
		{2, 0, 8, 7},
	}

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 13)
	P[1] = matrix.New(m, n, 13)
	E_A = matrix.New(m, n, 13)
	E_B_inv = matrix.New(m, n, 13)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 13))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 13))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 13))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 13))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a = finiteField.NewFieldElm(4, 13)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 101 test
	p0 = [][]int{
		{1, 79, 35, 20},
		{0, 31, 99, 76},
		{0, 77, 82, 77},
		{0, 90, 87, 50},
	}
	p1 = [][]int{
		{0, 78, 62, 84},
		{1, 69, 62, 19},
		{0, 0, 74, 51},
		{0, 12, 72, 100},
	}
	e_a = [][]int{
		{47, 61, 57, 86},
		{87, 27, 28, 30},
		{18, 85, 38, 21},
		{74, 28, 96, 33},
	}
	e_b_inv = [][]int{
		{97, 31, 35, 3},
		{63, 31, 19, 3},
		{100, 14, 67, 88},
		{44, 53, 45, 99},
	}

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 101)
	P[1] = matrix.New(m, n, 101)
	E_A = matrix.New(m, n, 101)
	E_B_inv = matrix.New(m, n, 101)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 101))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 101))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 101))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 101))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a = finiteField.NewFieldElm(33, 101)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 1009 test
	p0 = [][]int{
		{1, 153, 223, 642},
		{0, 943, 172, 293},
		{0, 991, 438, 411},
		{0, 527, 879, 927},
	}
	p1 = [][]int{
		{0, 149, 194, 624},
		{1, 654, 438, 968},
		{0, 587, 224, 79},
		{0, 196, 354, 819},
	}
	e_a = [][]int{
		{826, 431, 64, 368},
		{663, 225, 683, 94},
		{331, 914, 293, 0},
		{385, 442, 77, 278},
	}
	e_b_inv = [][]int{
		{471, 585, 1001, 967},
		{540, 7, 768, 888},
		{685, 162, 892, 637},
		{309, 950, 786, 346},
	}

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 1009)
	P[1] = matrix.New(m, n, 1009)
	E_A = matrix.New(m, n, 1009)
	E_B_inv = matrix.New(m, n, 1009)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 1009))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 1009))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 1009))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 1009))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a = finiteField.NewFieldElm(278, 1009)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 3359 test
	p0 = [][]int{
		{1, 2703, 2753, 2312},
		{0, 397, 1110, 2603},
		{0, 2062, 2831, 2029},
		{0, 1458, 1841, 403},
	}
	p1 = [][]int{
		{0, 2049, 1418, 2771},
		{1, 1800, 1861, 2709},
		{0, 3020, 1748, 2414},
		{0, 1376, 710, 2791},
	}
	e_a = [][]int{
		{1771, 3214, 64, 1970},
		{754, 3101, 2072, 2862},
		{1405, 2604, 962, 940},
		{1593, 366, 3211, 153},
	}
	e_b_inv = [][]int{
		{840, 2115, 3038, 1652},
		{2400, 232, 2390, 21},
		{1054, 342, 387, 1064},
		{22, 3073, 1213, 355},
	}

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 3359)
	P[1] = matrix.New(m, n, 3359)
	E_A = matrix.New(m, n, 3359)
	E_B_inv = matrix.New(m, n, 3359)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 3359))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 3359))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 3359))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 3359))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	a = finiteField.NewFieldElm(153, 3359)
	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 4091 test
	p0 = [][]int{
		{1, 1704, 1886, 1044},
		{0, 4031, 3534, 232},
		{0, 1172, 620, 295},
		{0, 1871, 870, 346},
	}
	p1 = [][]int{
		{0, 3288, 3541, 1851},
		{1, 3332, 547, 1254},
		{0, 649, 2656, 797},
		{0, 488, 3932, 3930},
	}
	e_a = [][]int{
		{2465, 530, 2191, 2795},
		{1819, 3662, 634, 3174},
		{485, 473, 739, 3320},
		{694, 3465, 1507, 826},
	}
	e_b_inv = [][]int{
		{2887, 1721, 1770, 2966},
		{99, 3961, 240, 288},
		{276, 2957, 3710, 1119},
		{3932, 3443, 2510, 176},
	}
	a = finiteField.NewFieldElm(826, 4091)

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 4091)
	P[1] = matrix.New(m, n, 4091)
	E_A = matrix.New(m, n, 4091)
	E_B_inv = matrix.New(m, n, 4091)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 4091))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 4091))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 4091))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 4091))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}

	// q = 4093 test
	p0 = [][]int{
		{1, 2424, 503, 1617},
		{0, 2969, 408, 1700},
		{0, 4033, 2346, 4024},
		{0, 2573, 3468, 3058},
	}
	p1 = [][]int{
		{0, 1190, 1263, 3475},
		{1, 1577, 2826, 1777},
		{0, 1290, 580, 3886},
		{0, 2643, 3358, 3967},
	}
	e_a = [][]int{
		{2282, 3060, 2207, 1989},
		{3307, 3469, 1807, 314},
		{3345, 962, 2727, 1805},
		{3062, 1562, 2641, 3471},
	}
	e_b_inv = [][]int{
		{3239, 3251, 2483, 2671},
		{1529, 3699, 2369, 1647},
		{3125, 1427, 2103, 3693},
		{2627, 3965, 1391, 3816},
	}
	a = finiteField.NewFieldElm(-622, 4093)

	P = make([]*matrix.Matrix, 2)

	P[0] = matrix.New(m, n, 4093)
	P[1] = matrix.New(m, n, 4093)
	E_A = matrix.New(m, n, 4093)
	E_B_inv = matrix.New(m, n, 4093)
	for i := 0; i < P[0].M; i++ {
		for j := 0; j < P[0].N; j++ {
			P[0].Set(i, j, finiteField.NewFieldElm(p0[i][j], 4093))
			P[1].Set(i, j, finiteField.NewFieldElm(p1[i][j], 4093))
			E_A.Set(i, j, finiteField.NewFieldElm(e_a[i][j], 4093))
			E_B_inv.Set(i, j, finiteField.NewFieldElm(e_b_inv[i][j], 4093))
		}
	}
	P[0] = P[0].Transpose()
	P[1] = P[1].Transpose()

	G = MatriciesToRows(P)

	// for i := 0; i < G.M; i++ {
	// 	for j := 0; j < G.N; j++ {
	// 		n := rand.Intn(int(q))
	// 		G.Set(i, j, finiteField.NewFieldElm(n, q))
	// 	}
	// }

	t.Logf("a[m-1, m-1]: %v\n", a)
	// for a.Equals(finiteField.NewFieldElm(0, q)) {
	// 	a = finiteField.NewFieldElm(rand.Intn(q), q)
	// }
	A, B_inv = Solve(G, a, m, n)
	if !A.Equals(E_A) || !B_inv.Equals(E_B_inv) {
		t.Errorf("A: %v\nB_inv: %v", A, B_inv)
	}
}

func TestInverse(test *testing.T) {

}
func TestParseHash(test *testing.T) {
	sk, _ := KeyGen()
	msg := []byte("This is a message")
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
	seeds, err := SeedTree(rho, alpha, t)
	if err != nil {
		test.Errorf("error: %v\n", err)
		return
	}
	G_tilde := make([]*matrix.Matrix, t)
	A_tilde := make([]*matrix.Matrix, t)
	B_tilde := make([]*matrix.Matrix, t)

	for i := 0; i < t; i++ {
		sigma_prime := make([]byte, l_salt+l_tree_seed+4)
		sigma_A_tilde := make([]byte, l_pub_seed)
		sigma_B_tilde := make([]byte, l_pub_seed)
		x, err := ToBytes(int32(math.Pow(2, math.Ceil(math.Log2(float64(t))))), 4)
		if err != nil {
			test.Errorf("error: %v\n", err)
			return
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
	LINE_16:
		xof = sha3.NewShake256()
		xof.Write(sigma_prime)
		xof.Read(sigma_A_tilde)
		xof.Read(sigma_B_tilde)
		xof.Read(sigma_prime)
		A_tilde[i] = ExpandInvMat(sigma_A_tilde, q, m)
		B_tilde[i] = ExpandInvMat(sigma_B_tilde, q, n)
		G_tilde[i] = Pi(A_tilde[i], G_0, B_tilde[i])
		G_tilde[i] = SF(G_tilde[i])
		if G_tilde[i] == nil {
			goto LINE_16
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

	count := 0
	for i := 0; i < t; i++ {
		if h[i] > 0 {
			count++
		}
	}

	if count != w {
		test.Errorf("Not enough non-zero values in h\nh:%v\n", h)
	}
}

func TestSeedTree(test *testing.T) {
	t := 5
	seed := []byte("seedseedseedseed")
	salt := []byte("saltsaltsaltsaltsaltsaltsaltsalt")
	seeds, err := SeedTree(seed, salt, t)

	test.Logf("Seed: %v\n", seed)
	test.Logf("Salt: %v\n", salt)
	test.Logf("T: %v\n", t)
	test.Logf("Seeds: %v\n", seeds)
	if err != nil {
		test.Errorf("error: %v\n", err)
	}
}

func TestSeedTreeToPath(test *testing.T) {
	seed := []byte("seedseedseedseed")
	salt := []byte("saltsaltsaltsaltsaltsaltsaltsalt")
	h := ParseHash(s, t, w, []byte("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh"))
	test.Logf("h: %v\n", h)
	path := SeedTreeToPath(w, t, h, seed, salt)
	test.Logf("len(path): %v\n", len(path))
	test.Errorf("path: %v\n", path)
}

func TestPathToSeedTree(test *testing.T) {
	seed := []byte("seedseedseedseed")
	salt := []byte("saltsaltsaltsaltsaltsaltsaltsalt")
	h := []byte{0, 0, 0, 1, 0}
	t := 5
	w := 1
	path := SeedTreeToPath(w, t, h, seed, salt)
	// test.Logf("path: %v\n", path)
	seeds := PathToSeedTree(h, path, salt, len(seed))
	test.Errorf("seeds: %v\n", seeds)
}

func TestBase(test *testing.T) {
	delta := Randombytes(l_sec_seed)
	sigma_G_0 := make([]byte, l_pub_seed)
	sigma_A := make([]byte, l_pub_seed)
	sigma_B := make([]byte, l_pub_seed)
	sigma_Ap := make([]byte, l_pub_seed)
	sigma_Bp := make([]byte, l_pub_seed)
	xof := sha3.NewShake256()
	xof.Write(delta)
	xof.Read(sigma_G_0)
	xof.Read(sigma_A)
	xof.Read(sigma_B)
	xof.Read(sigma_Ap)
	xof.Read(sigma_Bp)

	G_0 := ExpandSystMat(sigma_G_0, q, k, m, n)

	A := ExpandInvMat(sigma_A, q, m)
	B := ExpandInvMat(sigma_B, q, n)

	G_1 := Pi(A, G_0, B)
	G_1 = SF(G_1)

	Ap := ExpandInvMat(sigma_Ap, q, m)
	Bp := ExpandInvMat(sigma_Bp, q, n)

	G_1p := Pi(Ap, G_0, Bp)
	G_1p = SF(G_1p)

	G_1t := Pi(Ap.Mul(Inverse(A)), G_1, Inverse(B).Mul(Bp))
	G_1t = SF(G_1t)

	fmt.Printf("%v\n\n\n", G_1p)
	fmt.Printf("%v\n", G_1t)
}
