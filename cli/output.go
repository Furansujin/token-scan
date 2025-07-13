package cli

import (
    "encoding/json"
    "fmt"

    "github.com/fatih/color"
    "github.com/s-Amine/token-scan/token"
)

// PrintPretty prints TokenInfo in a readable colored form.
func PrintPretty(info *token.TokenInfo) {
    if info == nil {
        return
    }

    green := color.New(color.FgGreen).SprintFunc()
    yellow := color.New(color.FgYellow).SprintFunc()
    red := color.New(color.FgRed).SprintFunc()

    fmt.Printf("Token: %s (%s)\n", info.TokenName, info.TokenSymbol)
    fmt.Printf("Risk Score: %d/100 - ", info.RiskScore)
    switch info.RiskLevel {
    case "SAFE":
        fmt.Println(green(info.RiskLevel))
    case "WARNING":
        fmt.Println(yellow(info.RiskLevel))
    default:
        fmt.Println(red(info.RiskLevel))
    }

    if len(info.DetectedPatterns) > 0 {
        fmt.Printf("Suspicious Patterns: %v\n", info.DetectedPatterns)
    }

    if len(info.TopHolders) > 0 {
        fmt.Println("Top Holders:")
        i := 0
        for addr, pct := range info.TopHolders {
            if i >= 5 {
                break
            }
            fmt.Printf("  %s: %.2f%%\n", addr, pct)
            i++
        }
    }
}

// PrintJSON prints struct as JSON.
func PrintJSON(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(jsonData))
}
