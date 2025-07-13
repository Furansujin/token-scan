package main

import (
	"fmt"

	"github.com/s-Amine/token-scan/cli"
	"github.com/s-Amine/token-scan/scanners/customrules/sourcecheck"
	"github.com/s-Amine/token-scan/scanners/goplus"
	"github.com/s-Amine/token-scan/scanners/ishoneypot"
	"github.com/s-Amine/token-scan/scanners/multiscan"
	"github.com/s-Amine/token-scan/scanners/quickintel"
	"github.com/s-Amine/token-scan/token"
)

func main() {
	opts := cli.Parse()
	if opts.Token == "" {
		fmt.Println("token address required")
		return
	}

	var result interface{}
	var err error

	switch opts.Mode {
	case "multiscan":
		result = multiscan.Scan(opts.Token)
	case "goplus":
		result, err = goplus.Scan(opts.Token)
	case "ishoneypot":
		result, err = ishoneypot.Scan(opts.Token)
	case "quickIntel":
		result, err = quickintel.Scan(opts.Token)
	case "sourcecheck":
		result, err = sourcecheck.Scan(opts.Token)
	default:
		fmt.Println("invalid mode")
		return
	}

	if err != nil {
		fmt.Printf("scan error: %v\n", err)
		return
	}
	if opts.ScoreOnly {
		if info, ok := result.(*token.TokenInfo); ok {
			fmt.Println(info.TrustScore)
			return
		}
	}

	cli.Print(result, opts)
}
