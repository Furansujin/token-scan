package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/s-Amine/token-scan/cli"
	"github.com/s-Amine/token-scan/scanners/customrules/fingerprint"
	"github.com/s-Amine/token-scan/scanners/customrules/holders"
	"github.com/s-Amine/token-scan/scanners/customrules/sourcecheck"
	"github.com/s-Amine/token-scan/scanners/goplus"
	"github.com/s-Amine/token-scan/scanners/ishoneypot"
	"github.com/s-Amine/token-scan/scanners/multiscan"
	"github.com/s-Amine/token-scan/scanners/quickintel"
	"github.com/s-Amine/token-scan/token"
)

func main() {
	// Define command-line flags
	mode := flag.String("mode", "", "Mode of operation: multiscan, goplus, ishoneypot, quickIntel, sourcecheck, holders or fingerprint")
	tokenHash := flag.String("token", "", "Token hash to scan")
	pretty := flag.Bool("pretty", false, "Pretty colored output")
	scoreOnly := flag.Bool("score-only", false, "Output only the risk score")
	etherscanKey := flag.String("etherscan-key", "", "Etherscan API key")
	flag.Parse()

	if *mode == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *tokenHash == "" {
		fmt.Println("Error: Token hash is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var result interface{}
	var err error

	switch *mode {
	case "multiscan":
		result = multiscan.Scan(*tokenHash)
	case "goplus":
		result, err = goplus.Scan(*tokenHash)
	case "ishoneypot":
		result, err = ishoneypot.Scan(*tokenHash)
	case "quickIntel":
		result, err = quickintel.Scan(*tokenHash)
	case "sourcecheck":
		funcs, err2 := sourcecheck.Scan(*tokenHash, *etherscanKey)
		err = err2
		result = funcs
	case "holders":
		result, err = holders.Scan(*tokenHash, *etherscanKey)
	case "fingerprint":
		result, err = fingerprint.Scan(*tokenHash)
	default:
		fmt.Println("Error: Invalid mode specified")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error occurred during %s scan: %v\n", *mode, err)
		os.Exit(1)
	}

	if *scoreOnly {
		if info, ok := result.(*token.TokenInfo); ok {
			fmt.Println(info.RiskScore)
		} else {
			cli.PrintJSON(result)
		}
		return
	}

	if *pretty {
		if info, ok := result.(*token.TokenInfo); ok {
			cli.PrintPretty(info)
			return
		}
	}

	cli.PrintJSON(result)
}
