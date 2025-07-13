package heuristics

import "github.com/s-Amine/token-scan/token"

// Analyze applies heuristic rules on the token info and fills DetectedPatterns.
func Analyze(info *token.TokenInfo) {
	var patterns []string
	if info.IsHoneypot && info.IsMintable && !info.IsOpenSource {
		patterns = append(patterns, "Possible rug pull")
	}
	if info.TopHolderPercent > 70 {
		patterns = append(patterns, "Centralized ownership")
	}
	info.DetectedPatterns = patterns
}
