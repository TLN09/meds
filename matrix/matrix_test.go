package matrix

import (
	"math/rand"
	"meds/finiteField"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

const q int = 4093

func TestCompress(t *testing.T) {
	A := New(2, 3, q)
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Set(i, j, finiteField.NewFieldElm(n, q))
		}
	}

	R := Decompress(A.Compress(), A.M, A.N, A.Q)
	if !R.Equals(A) {
		t.Errorf("Compressed: %v\nDecompressed: %v\nA:            %v", A.Compress(), R.matrix, A.matrix)
	}
}

func TestIdentity(t *testing.T) {
	A := New(2, 3, q)
	E := true
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Set(i, j, finiteField.NewFieldElm(n, q))
		}
	}

	result := A.Equals(A)
	if result != E {
		t.Errorf("Matrix is not equal to itself")
	}
}

func TestUpperShiftCorrect(t *testing.T) {
	A := UpperShift(6, q)
	E := Matrix{
		6,
		6,
		q,
		[][]*finiteField.Fq{
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(1, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q)},
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(1, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q)},
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(1, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q)},
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(1, q), finiteField.NewFieldElm(0, q)},
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(1, q)},
			{finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q), finiteField.NewFieldElm(0, q)},
		},
	}
	if !E.Equals(A) {
		t.Errorf("U_6: %v", A)
	}
}

func TestEquals(t *testing.T) {
	A := New(2, 3, q)
	B := New(2, 3, q)
	E := true

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Set(i, j, finiteField.NewFieldElm(n, q))
			B.Set(i, j, finiteField.NewFieldElm(n, q))
		}
	}

	result := A.Equals(B)
	if result != E {
		t.Errorf("\nResult:   %v\nE: %v\nValues:\n\tA: %v\n\tB: %v\n", result, E, A.matrix, B.matrix)
	}

	B.Set(0, 1, B.Get(0, 1).Add(finiteField.NewFieldElm(1, B.Q)))
	E = false

	result = A.Equals(B)
	if result != E {
		t.Errorf("\nResult:   %v\nE: %v\nValues:\n\tA: %v\n\tB: %v\n", result, E, A.matrix, B.matrix)
	}
}

func TestAdd(t *testing.T) {
	A := New(2, 3, q)
	B := New(2, 3, q)
	E := New(2, 3, q)
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Set(i, j, finiteField.NewFieldElm(n, q))
			B.Set(i, j, finiteField.NewFieldElm(0, q))
			E.Set(i, j, finiteField.NewFieldElm(n, q))
		}
	}

	result := A.Add(B)
	if !result.Equals(E) {
		t.Errorf("\nR: %v\nE: %v\n", result.matrix, E.matrix)
	}

	for i := 0; i < B.M; i++ {
		for j := 0; j < B.N; j++ {
			n := rand.Intn(int(q))
			B.Set(i, j, finiteField.NewFieldElm(n, q))
			E.Set(i, j, A.Get(i, j).Add(finiteField.NewFieldElm(n, q)))
		}
	}

	result = A.Add(B)
	if !result.Equals(E) {
		t.Errorf("\nR: %v\nE: %v\n", result.matrix, E.matrix)
	}
}

func TestSub(t *testing.T) {
	A := New(2, 3, q)
	B := New(2, 3, q)
	E := New(2, 3, q)
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Set(i, j, finiteField.NewFieldElm(n, q))
			B.Set(i, j, finiteField.NewFieldElm(0, q))
			E.Set(i, j, finiteField.NewFieldElm(n, q))
		}
	}

	result := A.Sub(B)
	if !result.Equals(E) {
		t.Errorf("\nResult:   %v\nE: %v\n", result.matrix, E.matrix)
	}

	for i := 0; i < B.M; i++ {
		for j := 0; j < B.N; j++ {
			n := rand.Intn(int(q))
			B.Set(i, j, finiteField.NewFieldElm(n, q))
			E.Set(i, j, A.Get(i, j).Sub(finiteField.NewFieldElm(n, q)))
		}
	}

	result = A.Sub(B)
	if !result.Equals(E) {
		t.Errorf("\nResult:   %v\nE: %v\n", result.matrix, E.matrix)
	}
}

func TestScalar_mul(t *testing.T) {
	m := 2
	n := 5
	A := New(m, n, q)

	for k := 0; k < 100; k++ {
		scalar := finiteField.NewFieldElm(rand.Intn(int(q)), q)
		E := New(m, n, q)

		for i := 0; i < A.M; i++ {
			for j := 0; j < A.N; j++ {
				n := rand.Intn(int(q))
				A.Get(i, j).Set(n)
				E.Set(i, j, finiteField.NewFieldElm(n, q).Mul(scalar))
			}
		}

		result := A.Scalar_mul(scalar)
		if !result.Equals(E) {
			t.Errorf("R: %v\nE: %v", result.String(), E.String())
			return
		}
	}
}

func TestMul(t *testing.T) {
	A := New(5, 5, q)
	I := Identity(5, q)
	E := New(5, 5, q)

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Get(i, j).Set(n)
			E.Get(i, j).Set(n)
		}
	}

	result := A.Mul(I)
	if !result.Equals(E) {
		t.Errorf("\nResult:   %v\nE: %v\n", result.matrix, E.matrix)
	}

	A = New(2, 3, q)
	B := New(3, 2, q)
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Get(i, j).Set(n)
		}
	}
	for i := 0; i < B.M; i++ {
		for j := 0; j < B.N; j++ {
			n := rand.Intn(int(q))
			B.Get(i, j).Set(n)
		}
	}

	// Write A and B to a file
	t.Logf("A: %v", A.matrix)
	err := os.WriteFile("A_test.txt", A.Compress(), 0666)
	if err != nil {
		t.Errorf("Unable to create file A: %v", err)
		return
	}

	t.Logf("B: %v", B.matrix)
	err = os.WriteFile("B_test.txt", B.Compress(), 0666)
	if err != nil {
		t.Errorf("Unable to create file B: %v", err)
		return
	}

	// Compute using SageMath the correct matrix product and save to a file
	cmd := exec.Command("sage", "testMul.sage", strconv.Itoa(int(A.Q)), strconv.Itoa(A.M), strconv.Itoa(A.N), strconv.Itoa(B.M), strconv.Itoa(B.N))
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
	E = Decompress(file_content, A.M, B.N, A.Q)

	result = A.Mul(B)
	if result.M != E.M || result.N != E.N || !result.Equals(E) {
		t.Errorf("\nR: %v\nE: %v\n", result.matrix, E.matrix)
	}
}

func TestTranspose(t *testing.T) {
	m := 5
	n := 2
	A := New(m, n, q)
	E := New(n, m, q)

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Get(i, j).Set(n)
			E.Get(j, i).Set(n)
		}
	}

	result := A.Transpose()
	if !result.Equals(E) {
		t.Errorf("\nResult:   %v\nE: %v\n", result.matrix, E.matrix)
	}

	E = New(m, n, q)
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			n := rand.Intn(int(q))
			A.Get(i, j).Set(n)
			E.Get(i, j).Set(n)
		}
	}

	result = A.Transpose().Transpose()
	if !result.Equals(E) {
		t.Errorf("\nResult:   %v\nE: %v\n", result.matrix, E.matrix)
	}
}
