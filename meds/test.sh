PARAMETERSETS="167717"
OPS="Verify"
for p in $PARAMETERSETS
do
    for op in $OPS
    do
        go test -benchmem -count=20 -run=^$ -bench ^Benchmark${op}${p}$ meds
    done
done