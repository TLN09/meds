package matrix

import (
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

func (M *Matrix) Compress() []byte {
	length := int(math.Ceil(math.Log2(float64(M.Q))/float64(8))) * M.M * M.N
	b := make([]byte, length)
	f_byte := 0
	f_bit := 0
	q_bitlen := M.Get(0, 0).BitLen()
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			c := 0
			v := M.Get(i, j).Value()
			for c < q_bitlen {
				c_prime := min(8-f_bit, q_bitlen-c)
				b[f_byte] += byte(v % int(math.Pow(float64(2), float64(c_prime))) * int(math.Pow(float64(2), float64(f_bit))))
				v = int(math.Floor(float64(v) / math.Pow(float64(2), float64(c_prime))))
				c += c_prime
				f_bit += c_prime
				if f_bit == 8 {
					f_bit = 0
					f_byte++
				}
			}
		}
	}

	return b
}

func Decompress(b []byte, m int, n int, q int) *Matrix {
	M := New(m, n, q)
	f_byte := 0
	f_bit := 0
	q_bitlen := M.Get(0, 0).BitLen()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			c := 0
			v := 0
			for c < q_bitlen {
				c_prime := min(8-f_bit, q_bitlen-c)
				v += int(math.Floor(float64(b[f_byte])/math.Pow(float64(2), float64(f_bit)))) % int(math.Pow(2, float64(c_prime))) * int(math.Pow(2, float64(c)))
				c += c_prime
				f_bit += c_prime
				if f_bit == 8 {
					f_bit = 0
					f_byte++
				}
			}

			M.Set(i, j, finiteField.NewFieldElm(v, M.Q))
		}
	}

	return M
}

func (M *Matrix) UnaryMinus() *Matrix {
	for i := 0; i < M.M; i++ {
		for j := 0; j < M.N; j++ {
			M.Set(i, j, M.Get(i, j).UnaryMinus())
		}
	}

	return M
}
