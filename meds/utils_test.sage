def decompress(data: bytes, m: int, n: int, q: int, q_length: int):
    M = [[0 for _ in range(n)] for _ in range(m)]
    idx = 0
    for i in range(m):
        for j in range(n):
            M[i][j] = int.from_bytes(data[idx:idx+q_length], 'big')
            idx += q_length
    return Matrix(GF(q), M)

def compress(M, q_length):
    b = b''
    for row in M:
        for elm in row:
            b += int(elm).to_bytes(q_length, 'big')
    return b

def compress_and_write(M, q_length, submatrix):
    path = 'E_submatrix_test.txt' if submatrix else 'E_test.txt'
    with open(path, 'wb') as f:
        f.write(compress(M, q_length))

def test_sf(q, m, n, submatrix):
    q_length = 2

    if submatrix:
        with open('SF_submatrix_test.txt', 'rb') as f:
            data = f.read()
    else:
        with open('SF_test.txt', 'rb') as f:
            data = f.read()

    A = decompress(data, m, n, q, q_length)
    C = LinearCode(A)
    Cs, p  = C.standard_form()
    compress_and_write(Cs.systematic_generator_matrix(), q_length, submatrix)

if __name__ == '__main__':
    from sys import argv

    if sys.argv[1] == 'sf':
        q, m, n = int(argv[2]), int(argv[3]), int(argv[4])
        test_sf(q, m, n, False)
    elif sys.argv[1] == 'sf_submatrix':
        q, m, n = int(argv[2]), int(argv[3]), int(argv[4])
        test_sf(q, m, n, True)