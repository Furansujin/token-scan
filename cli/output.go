package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/s-Amine/token-scan/token"
)

// Print renders the result based on options.
func Print(data interface{}, opts Options) {
	if info, ok := data.(*token.TokenInfo); ok {
		if opts.Pretty {
			printPretty(info)
			return
		}
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}

func printPretty(info *token.TokenInfo) {
	fmt.Printf("Token: %s (%s)\n", info.TokenName, info.TokenSymbol)
	levelColor := color.New(color.FgGreen)
	switch info.RiskLevel {
	case "WARNING":
		levelColor = color.New(color.FgYellow)
	case "DANGEROUS":
		levelColor = color.New(color.FgRed)
	}
	levelColor.Printf("Risk Score: %d/100 - %s\n", info.TrustScore, info.RiskLevel)

	fmt.Println("Top Holders:")
	if len(info.TopHolders) > 0 {
		for addr, pct := range info.TopHolders {
			fmt.Printf("  %s : %.2f%%\n", addr, pct)
		}
	} else {
		fmt.Println("  none")
	}

	fmt.Printf("Wash Trading Score: %.2f\n", info.WashTradeScore)
	if info.LiquidityLocked {
		fmt.Printf("LP Status: Locked until %v\n", time.Unix(info.LiquidityUnlockTime, 0).Format(time.RFC3339))
	} else {
		fmt.Println("LP Status: Not Locked")
	}
	if info.IsHoneypot {
		color.New(color.FgRed).Println("Honeypot detected!")
	} else {
		fmt.Println("Honeypot: none")
	}

	fmt.Print("Bot Activity: ")
	if info.BotActivity {
		fmt.Printf("Detected (%d micro tx)\n", info.MicroTxCount)
	} else {
		fmt.Println("none")
	}

	if info.IsClone {
		fmt.Printf("Clone of: %s\n", info.OriginalTokenAddress)
	} else {
		fmt.Println("Clone: none")
	}

	if info.AnnouncedBy != "" {
		fmt.Printf("Announced by: %s\n", info.AnnouncedBy)
	} else {
		fmt.Println("Announced by: none")
	}

	if len(info.DetectedPatterns) > 0 {
		fmt.Printf("Suspicious Patterns: %v\n", info.DetectedPatterns)
	} else {
		fmt.Println("Suspicious Patterns: none")
	}

	if len(info.MatchedScamContracts) > 0 {
		fmt.Printf("Matched Scam Contracts: %v\n", info.MatchedScamContracts)
	} else {
		fmt.Println("Matched Scam Contracts: none")
	}

	if len(info.HasSensitiveFunctions) > 0 {
		fmt.Printf("Sensitive Functions: %v\n", info.HasSensitiveFunctions)
	} else {
		fmt.Println("Sensitive Functions: none")
	}
}
