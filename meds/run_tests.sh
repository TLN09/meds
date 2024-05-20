PARAMETERSETS="9923 13220 41711 69497 134180 167717"
OPS="Keygen Sign Verify"
for op in $OPS
do
    for p in $PARAMETERSETS
    do
        go test -benchmem -count=128 -run=^$ -bench ^Benchmark${op}${p}$ meds > ${op}${p}_benchmark.txt
    done
done