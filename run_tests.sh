OPS="KeyGen Sign Verify"
rm *_test.txt
for op in $OPS
do
    go test ./meds -timeout=30m -run ^Test${op} meds > ${op}_test.txt
done

cat *_test.txt
