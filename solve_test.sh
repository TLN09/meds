go test ./meds -timeout=30m -benchmem -count=10 -run=^$ -bench ^BenchmarkSolve meds > non_optimized_solve_benchmark.txt
