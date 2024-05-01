package main

import (
	"flag"
	"fmt"
)

func main() {
	secLvl := flag.Int("seclvl", 1, "NIST security level (1, 3, 5)")
	op := flag.String("op", "None", "The MEDS operation to do (Keygen, Sign, Verify)")

	flag.Parse()
	fmt.Printf("seclvl: %v\nop: %v\n", *secLvl, *op)

	switch *secLvl {
	case 1:

	}
}
