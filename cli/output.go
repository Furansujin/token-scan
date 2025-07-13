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
	if len(info.TopHolders) > 0 {
		fmt.Println("Top Holders:")
		for addr, pct := range info.TopHolders {
			fmt.Printf("  %s : %.2f%%\n", addr, pct)
		}
	}
	if info.WashTradeScore > 0.1 {
		fmt.Printf("Wash Trading Score: %.2f\n", info.WashTradeScore)
	}
	if info.LiquidityLocked {
		fmt.Printf("LP Status: Locked until %v\n", time.Unix(info.LiquidityUnlockTime, 0).Format(time.RFC3339))
	} else {
		fmt.Println("LP Status: Not Locked")
	}
	if info.IsHoneypot {
		color.New(color.FgRed).Println("Honeypot detected!")
	}
	if info.BotActivity {
		fmt.Printf("Bot Activity: Detected (%d micro tx)\n", info.MicroTxCount)
	}
	if info.IsClone {
		fmt.Printf("Clone of: %s\n", info.OriginalTokenAddress)
	}
	if info.AnnouncedBy != "" {
		fmt.Printf("Announced by: %s\n", info.AnnouncedBy)
	} else if !info.SocialProof {
		fmt.Println("Announced by: none")
	}
	if len(info.DetectedPatterns) > 0 {
		fmt.Printf("Suspicious Patterns: %v\n", info.DetectedPatterns)
	}
	if len(info.MatchedScamContracts) > 0 {
		fmt.Printf("Matched Scam Contracts: %v\n", info.MatchedScamContracts)
	}
	if len(info.HasSensitiveFunctions) > 0 {
		fmt.Printf("Sensitive Functions: %v\n", info.HasSensitiveFunctions)
	}
}
