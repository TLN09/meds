from sys import argv

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

def compress_and_write(M, q_length):
    with open('E_test.txt', 'wb') as f:
        f.write(compress(M, q_length))
    

def test_mul(q: int, a_m: int, a_n: int, b_m: int, b_n: int):
    q_length = len(bin(q)[2:])
    if q_length % 8 != 0:
        q_length += 8 - (q_length % 8)
    q_length = q_length // 8
    
    with open('A_test.txt', 'rb') as f:
        data = f.read()

    A = decompress(data, a_m, a_n, q, q_length)
    
    with open('B_test.txt', 'rb') as f:
        data = f.read()
    
    B = decompress(data, b_m, b_n, q, q_length)

    R = A * B
    print("A:\n", A, sep='')
    print("B:\n", B, sep='')
    compress_and_write(R, q_length)

        
if __name__ == '__main__':
    q, a_m, a_n, b_m, b_n = int(argv[1]), int(argv[2]), int(argv[3]), int(argv[4]), int(argv[5])
    test_mul(q, a_m, a_n, b_m, b_n)
