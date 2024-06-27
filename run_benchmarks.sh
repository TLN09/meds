OPS="KeyGen Sign Verify"
for op in $OPS
do
    go test ./meds -timeout=30m -benchmem -count=10 -run=^$ -bench ^Benchmark${op} meds > ${op}_benchmark.txt
done