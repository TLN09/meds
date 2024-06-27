OPS="KeyGen Sign Verify"
for op in $OPS
do
    go test -timeout=30m -run ^Test${op} meds > ${op}_test.txt
done
