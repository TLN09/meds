PARAMETERSETS="9923 13220 41711 69497 134180 167717"
OPS="Keygen Sign Verify"
for p in $PARAMETERSETS
do
    for op in $OPS
    do
        benchstat ${op}${p}_benchmark.txt >> benchmark_res.txt
    done
done