package main

import (
	"flag"
	"fmt"
	"meds/meds"
	"os"
)

func main() {
	parameterSet := flag.Int("meds", 9923, "MEDS parameter set (9923, 13220, 41711, 69497, 134180, 167717)")
	op := flag.String("op", "", "The MEDS operation to do (Keygen, Sign, Verify)")
	msg_file := flag.String("msg", "example.txt", "The message to sign")
	signed_file := flag.String("signed", "example.txt.signed", "Path to the message to verify")

	flag.Parse()

	switch *parameterSet {
	case 9923:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	case 13220:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	case 41711:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	case 69497:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	case 134180:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	case 167717:
		fmt.Printf("Chosen Parameterset: %v\n", *parameterSet)
	default:
		fmt.Printf("Invalid parameter set\n")
		flag.Usage()
		return
	}
	meds.ParameterSetup(*parameterSet)
	switch *op {
	case "Keygen":
		pk, sk := meds.KeyGen()

		err := os.WriteFile("meds_key", sk, 0666)
		if err != nil {
			fmt.Printf("Error writing to private key file\n")
			return
		}
		err = os.WriteFile("meds_key.pub", pk, 0666)
		if err != nil {
			fmt.Printf("Error writing to public key file\n")
			return
		}
		fmt.Printf("Keys saved to files\n")
	case "Sign":
		sk, err := os.ReadFile("meds_key")
		if err != nil {
			fmt.Printf("Error reading private key\n")
			return
		}
		msg, err := os.ReadFile(*msg_file)
		if err != nil {
			fmt.Printf("Error reading message file\n")
			return
		}
		signed_msg, err := meds.Sign(sk, msg)
		if err != nil {
			fmt.Printf("Error signing message. %v\n", err)
			return
		}
		err = os.WriteFile(fmt.Sprintf("%v.signed", (*msg_file)), signed_msg, 0666)
		if err != nil {
			fmt.Printf("Error writing signed message. %v\n", err)
			return
		}
		fmt.Printf("Message signed successfully\n")
	case "Verify":
		pk, err := os.ReadFile("meds_key.pub")
		if err != nil {
			fmt.Printf("Error reading public key file\n")
			return
		}
		msg_signed, err := os.ReadFile(*signed_file)
		if err != nil {
			fmt.Printf("Error reading signed file\n")
			return
		}
		msg := meds.Verify(pk, msg_signed)
		if msg == nil {
			fmt.Printf("Invalid Signature\n")
			return
		}
		fmt.Printf("Valid Signature\n")
	default:
		fmt.Printf("Invalid Operation\n")
		flag.Usage()
	}
}
