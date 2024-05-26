OPS="Keygen Sign Verify"
for op in $OPS
do
    benchstat ${op}*_benchmark.txt >> benchmark_res.txt
done