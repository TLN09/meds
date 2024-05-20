def solve_symb(P0prime, Amm):
    m = P0prime[0].nrows()
    n = P0prime[0].ncols()

    #print("P0prime[0].transpose():")
    print(f"{P0prime[0].transpose()}", end="\n\n")
    #print()
    #print("P0prime[1].transpose():")
    print(f"{P0prime[1].transpose()}", end="\n\n")
    #print()


    #print("-P0prime[0].transpose():")
    #print(-P0prime[0].transpose())
    #print()
    #print("-P0prime[1].transpose():")
    #print(-P0prime[1].transpose())
    #print()

    GFq = Amm.base_ring()

    Pj = [None] * 2

    Pj[0] = matrix(GFq, m, n, [[GFq(1) if i==j else GFq(0) for i in range(n)] for j in range(m)])
    #print("Pj[0]")
    #print(Pj[0], end="\n\n")
    Pj[1] = matrix(GFq, m, n, [[GFq(1) if i==j else GFq(0) for i in range(n)] for j in range(1,m)] + [[GFq(0)]*n])
    #print("Pj[1]")
    #print(Pj[1], end="\n\n")


    R = PolynomialRing(GFq, m*m + n*n,
    names = ','.join([f"b{i}_{j}" for i in range(n) for j in range(n)]) + "," \
            + ','.join([f"a{i}_{j}" for i in range(m) for j in range(m)]))
    #print("R")
    #print(R)
    #print()

    A     = matrix(R, m, var(','.join([f"a{i}_{j}" for i in range(m) for j in range(m)])))
    B_inv = matrix(R, n, var(','.join([f"b{i}_{j}" for i in range(n) for j in range(n)])))
    A[m-1,m-1] = Amm
    #print("A", A, sep='\n', end="\n\n")
    #print("B_inv", B_inv, sep='\n', end="\n\n")
    #print("P0prime[0]")
    #print(P0prime[0])
    #print()
    #print("A * P0prime[0]")
    #print(A*P0prime[0])
    #print()
    #print("Pj[0] * B_inv")
    #print(Pj[0] * B_inv)
    #print()
    eqs1 = Pj[0] * B_inv - A*P0prime[0]
    #print("eqs1")
    #print(eqs1.coefficients())
    #print()
    eqs2 = Pj[1] * B_inv - A*P0prime[1]
    #print("eqs2")
    #print(eqs2.coefficients())
    #print()

    eqs = eqs1.coefficients() + eqs2.coefficients()[:-1]
    #print("eqs")
    #print(eqs)
    #print(len(eqs))
    #print(eqs)

    #for eq in eqs:
    #    print(eq)
    #    print(eq.coefficients())
    #    print(-eq.constant_coefficient())
    #    print()

    #print("R.gens()")
    #print(R.gens())
    #print()

    rsys = matrix(GFq, [[eq.coefficient(v) for v in R.gens()[:-1]] + [-eq.constant_coefficient()] for eq in eqs])

    #print('rsys')
    #print(str(rsys).replace("0", " "))
    #print(rsys.nrows())
    #print()

    rsys_rref = rsys.rref()

    #print('rsys_rref')
    #print(str(rsys_rref).replace("0", " "))
    #print()

    if not all([rsys_rref[i][i] == 1 for i in range(rsys_rref.nrows())]):
        return None, None

    sol = rsys_rref.columns()[-1].list()

    A = matrix(GFq, m, sol[n*n:] + [Amm])
    B_inv = matrix(GFq, m, sol[:n*n])

    return A, B_inv


q = 4093

m = 30
n = 30
k = 30

GFq = GF(q)

I_k = matrix.identity(ring=GFq, n=k)

G0prime = I_k.augment(matrix(GFq, k, m*n-k, [GFq.random_element() for _ in range(k*(m*n-k))])) 
#print("G0prime")
#print(G0prime)

Amm = 0
while Amm == 0:
  Amm = GFq.random_element()

#print("G0prime.rows()")
#print(G0prime.rows())
check_A_i, check_B_i = solve_symb([matrix(GFq, m, n, G0prime.rows()[i]) for i in range(2)], Amm)

print(f"{Amm}", end="\n\n")
print(f"{check_A_i}", end="\n\n")
print(f"{check_B_i}", end="\n\n")
