OPS="Keygen Sign Verify"
for op in $OPS
do
    go test -timeout=30m -benchmem -count=10 -run=^$ -bench ^Benchmark${op} meds > ${op}_benchmark.txt
done