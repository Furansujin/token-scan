package score

import "github.com/s-Amine/token-scan/token"

// Compute calculates the trust score and sets the risk level.
func Compute(info *token.TokenInfo) int {
	score := 0
	if info.IsHoneypot {
		score += 30
	}
	if !info.LiquidityLocked {
		score += 20
	}
	if info.TopHolderPercent > 50 {
		score += 15
	}
	if info.IsClone {
		score += 10
	}
	if info.WashTradeScore > 0.2 {
		score += 10
	}
	if info.BotActivity {
		score += 10
	}
	if !info.SocialProof {
		score += 5
	}
	info.TrustScore = score
	switch {
	case score <= 30:
		info.RiskLevel = "SAFE"
	case score <= 70:
		info.RiskLevel = "WARNING"
	default:
		info.RiskLevel = "DANGEROUS"
	}
	return score
}
