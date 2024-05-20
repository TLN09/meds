package matrix

import (
	"fmt"
	"math"
	"meds/finiteField"
	"strings"
)

type Matrix struct {
	M      int
	N      int
	Q      int
	matrix [][]*finiteField.Fq
}

// Get returns the element at position (i, j) in the matrix
// Precondition: i > 0 && j > 0 && i < A.M && j < A.N
// Returns $a_{ij}$
func (A *Matrix) Get(i, j int) *finiteField.Fq {
	return A.matrix[i][j]
}

func (A *Matrix) Set(i, j int, elm *finiteField.Fq) {
	A.matrix[i][j] = elm
}

func (A *Matrix) Submatrix(startRow, endRow, startCol, endCol int) *Matrix {
	M := New(endRow-startRow, endCol-startCol, A.Q)

	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			M.Set(i, j, A.Get(i+startRow, j+startCol))
		}
	}

	return M
}

func (A *Matrix) String() string {
	var str strings.Builder
	str.WriteString("[\n")
	for i := 0; i < A.M; i++ {
		str.WriteString("[")
		for j := 0; j < A.N; j++ {
			str.WriteString(A.Get(i, j).String())
			if j != A.N-1 {
				str.WriteString(", ")
			}
		}
		str.WriteString("]\n")
	}
	str.WriteString("]")
	return str.String()
}

// New initializes a new Matrix
// Precondition: m > 0 and n > 0
// Returns: $M_{mn}$ initialized to all zeroes
func New(m int, n int, q int) *Matrix {
	matrix := make([][]*finiteField.Fq, m)

	for i := range matrix {
		matrix[i] = make([]*finiteField.Fq, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			matrix[i][j] = finiteField.NewFieldElm(0, q)
		}
	}

	return &Matrix{m, n, q, matrix}
}

func New_with_default(m, n, q, val int) *Matrix {
	matrix := make([][]*finiteField.Fq, m)

	for i := range matrix {
		matrix[i] = make([]*finiteField.Fq, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			matrix[i][j] = finiteField.NewFieldElm(val, q)
		}
	}

	return &Matrix{m, n, q, matrix}
}

// Identity creates an identity matrix of the specified size
// Returns: $I_n$
func Identity(n int, q int) *Matrix {
	I := New(n, n, q)

	for i := 0; i < n; i++ {
		I.Get(i, i).Set(1)
	}

	return I
}

// UpperShift creates an Upper shift matrix of size d (identity matrix shifted one column to the right)
func UpperShift(d int, q int) *Matrix {
	U := New(d, d, q)

	for i := 1; i < d; i++ {
		U.Get(i-1, i).Set(1)
	}

	return U
}

// Equals is the equality operation on matricies
// Precondition: Matricies are of the same dimentions
// Returns: $A = B$
func (A *Matrix) Equals(B *Matrix) bool {
	equal := B != nil

	for i := 0; equal && i < A.M; i++ {
		for j := 0; equal && j < A.N; j++ {
			equal = A.Get(i, j).Equals(B.Get(i, j))
		}
	}

	return equal
}

// Add is the addition operation on matricies.
// Precondition: Matricies are of the same dimentions
// Returns: $A + B$
func (A *Matrix) Add(B *Matrix) *Matrix {
	R := New(A.M, A.N, A.Q)

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			R.Set(i, j, A.Get(i, j).Add(B.Get(i, j)))
		}
	}

	return R
}

// Sub is the subtraction operation on matricies
// Precondition: Matricis are of the same dimentions
// Returns: $A - B$
func (A *Matrix) Sub(B *Matrix) *Matrix {
	R := New(A.M, A.N, A.Q)

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			R.Set(i, j, A.Get(i, j).Sub(B.Get(i, j)))
		}
	}

	return R
}

// Scalar_mul is the scalar multiplication operation on matricies
// Returns: c A
func (A *Matrix) Scalar_mul(c *finiteField.Fq) *Matrix {
	R := New(A.M, A.N, A.Q)

	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			R.Set(i, j, A.Get(i, j).Mul(c))
		}
	}

	return R
}

// Mul is the multiplication operation on matricies
// Precondition: A.N == B.M
// Returns: $A \cdot B$
func (A *Matrix) Mul(B *Matrix) *Matrix {
	R := New(A.M, B.N, A.Q)

	for i := 0; i < R.M; i++ {
		for j := 0; j < R.N; j++ {
			elm := finiteField.NewFieldElm(0, R.Q)
			for k := 0; k < A.N; k++ {
				elm = elm.Add(A.Get(i, k).Mul(B.Get(k, j)))
			}
			R.Set(i, j, elm)
		}
	}

	return R
}

// Transpose is the transpose operation on a Matrix
// Returns: $A^T$
func (A *Matrix) Transpose() *Matrix {
	R := New(A.N, A.M, A.Q)

	for i := 0; i < A.N; i++ {
		for j := 0; j < A.M; j++ {
			R.Set(i, j, A.Get(j, i))
		}
	}

	return R
}

// Kroenecker_product calculates the Kroenecker Product of two matricies
// Returns: $A \otimes B$
func (A *Matrix) Kroenecker_product(B *Matrix) *Matrix {
	R := New(A.M*B.M, A.N*B.N, A.Q)

	all_a_ij_times_B := make([]Matrix, A.M*A.N)
	idx := 0
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			all_a_ij_times_B[idx] = *B.Scalar_mul(A.Get(i, j))
			idx++
		}
	}

	row_off := 0
	col_off := 0
	for idx, M := range all_a_ij_times_B {
		row_off = (idx / A.N) * B.M
		col_off = (idx % A.N) * B.N
		for i := 0; i < M.M; i++ {
			for j := 0; j < M.N; j++ {
				R.Set(row_off+i, col_off+j, M.Get(i, j))
			}
		}
	}

	return R
}

func (A *Matrix) Compress() []byte {
	q_length := 16
	// Alignment to an even number of bytes
	if q_length%8 != 0 {
		q_length += 8 - (q_length % 8)
	}
	q_length /= 8
	b := make([]byte, A.N*A.M*q_length)

	idx := 0
	for i := 0; i < A.M; i++ {
		for j := 0; j < A.N; j++ {
			b_ij := A.Get(i, j).Bytes()
			b[idx] = b_ij[0]
			b[idx+1] = b_ij[1]
			idx += 2
		}
	}

	return b
}

func Decompress(b []byte, m int, n int, q int) *Matrix {
	M := New(m, n, q)
	q_length := M.Get(0, 0).BitLen()
	// Alignment to an even number of bytes
	if q_length%8 != 0 {
		q_length += 8 - (q_length % 8)
	}
	q_length /= 8
	idx := 0

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			M.Set(i, j, finiteField.NewFromBytes(b[idx:idx+q_length], M.Q))
			idx += q_length
		}
	}

	return M
}

func (M *Matrix) nonFiniteFieldREF() [][]float64 {
	m := make([][]float64, M.M)
	// Copy values in matrix
	for i := 0; i < M.M; i++ {
		m[i] = make([]float64, M.N)
	}
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			m[i][j] = float64(M.Get(i, j).Value())
		}
	}

	// Make REF
	for i := 0; i < M.M; i++ {
		// fmt.Printf("%v\n", sf)
		l := i
		zero_col := true

		for j := 0; j < M.M && zero_col; j++ {
			zero_col = m[i][j] == 0
		}

		if zero_col {
			return nil
		}
		for k := i + 1; k < M.M && m[i][l] == 0; k++ {
			for x := 0; x < len(m[0]); x++ {
				m[i][x], m[k][x] = m[k][x], m[i][x]
			}
		}

		c := 1 / m[i][l]
		for j := 0; j < M.N; j++ {
			m[i][j] *= c
		}

		for k := i + 1; k < M.M; k++ {
			c := -m[k][l]
			for j := 0; j < M.N; j++ {
				m[k][j] += m[i][j] * c
			}
		}
	}

	return m
}

func (M *Matrix) Det() int {
	ref := M.nonFiniteFieldREF()
	if ref == nil {
		return 1
	}
	det := float64(1)
	for i := 0; i < M.M; i++ {
		det *= ref[i][i]
	}
	return int(math.Round(det))
}

func (M *Matrix) Invertable() bool {
	return M.Det() != 0
}

func (M *Matrix) UnaryMinus() *Matrix {
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			M.Set(i, j, M.Get(i, j).UnaryMinus())
		}
	}

	return M
}

func (M *Matrix) Flatten() []*finiteField.Fq {
	r := make([]*finiteField.Fq, M.M*M.N)
	idx := 0

	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			r[idx] = M.Get(i, j)
			idx++
		}
	}

	return r
}

// Precondition: B has the correct shape to calculate A * B symbolicly
func (A *Matrix) SymMul(B [][]string) *[][]string {
	m := A.M
	n := len(B[0])
	R := make([][]string, m)
	for i := 0; i < m; i++ {
		R[i] = make([]string, n)
		for j := 0; j < n; j++ {
			elm := []string{}
			for k := 0; k < A.N; k++ {
				if !A.Get(i, k).Equals(finiteField.NewFieldElm(0, A.Q)) {
					elm = append(elm, fmt.Sprintf("%v * %v", A.Get(i, k), B[k][j]))
				}
			}
			R[i][j] = strings.Join(elm, " + ")
		}
	}

	return &R
}

func SymMul(A *[][]string, B *Matrix) *[][]string {
	m := len(*A)
	n := B.N
	R := make([][]string, m)
	for i := 0; i < m; i++ {
		R[i] = make([]string, n)
		for j := 0; j < n; j++ {
			elm := []string{}
			for k := 0; k < B.M; k++ {
				if !B.Get(k, j).Equals(finiteField.NewFieldElm(0, B.Q)) {
					elm = append(elm, fmt.Sprintf("%v * %v", B.Get(k, j), (*A)[i][k]))
				}
			}
			if len(elm) == 0 {
				R[i][j] = "0"
			} else {
				R[i][j] = strings.Join(elm, " + ")
			}
		}
	}

	return &R
}

func SymSub(A *[][]string, B *[][]string) *[][]string {
	m := len(*A)
	n := len((*A)[0])
	R := make([][]string, m)

	for i := 0; i < m; i++ {
		R[i] = make([]string, n)
		for j := 0; j < n; j++ {
			R[i][j] = fmt.Sprintf("%v#-#%v", (*A)[i][j], strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll((*B)[i][j], "-", "$"), "+", "-"), "$", "+"))
		}
	}

	return &R
}