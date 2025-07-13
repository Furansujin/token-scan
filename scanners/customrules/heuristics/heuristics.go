package heuristics

import "github.com/s-Amine/token-scan/token"

// Detect inspects TokenInfo and returns a list of patterns that might indicate
// suspicious behaviour.
func Detect(info *token.TokenInfo) []string {
	if info == nil {
		return nil
	}

	patterns := []string{}

	if info.IsHoneypot && !info.IsOpenSource && info.IsMintable {
		patterns = append(patterns, "Possible rug pull")
	}

	if info.TopHolderPercent > 50 && !info.ContractRenounced {
		patterns = append(patterns, "Centralized ownership")
	}

	if info.IsBlacklisted && info.TransferPausable {
		patterns = append(patterns, "Blacklisting and pausable transfers")
	}

	return patterns
}
