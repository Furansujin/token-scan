package multiscan

import (
	"github.com/s-Amine/token-scan/scanners/customrules/clone"
	"github.com/s-Amine/token-scan/scanners/customrules/fingerprint"
	"github.com/s-Amine/token-scan/scanners/customrules/heuristics"
	"github.com/s-Amine/token-scan/scanners/customrules/holders"
	"github.com/s-Amine/token-scan/scanners/customrules/honeypot"
	"github.com/s-Amine/token-scan/scanners/customrules/liquidity"
	"github.com/s-Amine/token-scan/scanners/customrules/microtx"
	"github.com/s-Amine/token-scan/scanners/customrules/score"
	"github.com/s-Amine/token-scan/scanners/customrules/social"
	"github.com/s-Amine/token-scan/scanners/customrules/sourcecheck"
	"github.com/s-Amine/token-scan/scanners/customrules/washtrade"
	"github.com/s-Amine/token-scan/scanners/goplus"
	"github.com/s-Amine/token-scan/scanners/ishoneypot"
	"github.com/s-Amine/token-scan/scanners/quickintel"
	"github.com/s-Amine/token-scan/token"
)

// MultiScan performs multiple scans using different scanners and unifies the results into one TokenInfo.
func Scan(tokenHash string) *token.TokenInfo {
	// Channels to receive scan results from different scanners
	goPlusScanResultChan := make(chan *token.TokenInfo)
	isHoneypotScanResultChan := make(chan *token.TokenInfo)
	quickIntelScanResultChan := make(chan *token.TokenInfo)

	// Perform GoPlus scan concurrently
	go func() {
		goPlusScanResult, _ := goplus.Scan(tokenHash)
		goPlusTokenInfo := token.InitTokenInfoFromGoPlus(goPlusScanResult)
		goPlusScanResultChan <- goPlusTokenInfo
	}()

	// Perform isHoneypot scan concurrently
	go func() {
		isHoneypotScanResult, _ := ishoneypot.Scan(tokenHash)
		honeypotTokenInfo := token.InitTokenInfoFromHoneypotResponse(isHoneypotScanResult)
		isHoneypotScanResultChan <- honeypotTokenInfo
	}()

	// Perform QuickIntel scan concurrently
	go func() {
		quickIntelScanResult, _ := quickintel.Scan(tokenHash)
		quickIntelTokenInfo := token.InitTokenInfoFromQuickIntelResponse(quickIntelScanResult)
		quickIntelScanResultChan <- quickIntelTokenInfo
	}()

	// Variables to store scan results
	var goPlusTokenInfo *token.TokenInfo
	var honeypotTokenInfo *token.TokenInfo
	var quickIntelTokenInfo *token.TokenInfo

	// Receive scan results from channels
	for i := 0; i < 3; i++ {
		select {
		case goPlusTokenInfo = <-goPlusScanResultChan:
		case honeypotTokenInfo = <-isHoneypotScanResultChan:
		case quickIntelTokenInfo = <-quickIntelScanResultChan:
		}
	}

	// Unify scan results into one TokenInfo
	unifiedInfo := token.UnifyTokenInfo(goPlusTokenInfo, honeypotTokenInfo, quickIntelTokenInfo)

	// Additional analyses
	if funcs, err := sourcecheck.Scan(tokenHash); err == nil {
		unifiedInfo.HasSensitiveFunctions = funcs
	}

	if pct, top, err := holders.Scan(tokenHash); err == nil {
		unifiedInfo.TopHolderPercent = pct
		unifiedInfo.TopHolders = top
	}

	if score, err := washtrade.Scan(tokenHash); err == nil {
		unifiedInfo.WashTradeScore = score
	}

	if locked, unlock, err := liquidity.Scan(tokenHash); err == nil {
		unifiedInfo.LiquidityLocked = locked
		unifiedInfo.LiquidityUnlockTime = unlock
	}

	if isH, errStr := honeypot.Scan(tokenHash); errStr != "" || isH {
		unifiedInfo.IsHoneypot = isH
		unifiedInfo.HoneypotError = errStr
	}

	if count, bot := microtx.Scan(tokenHash); bot {
		unifiedInfo.MicroTxCount = count
		unifiedInfo.BotActivity = bot
	}

	if unifiedInfo.TokenName != "" {
		if cl, orig := clone.Scan(unifiedInfo.TokenName, unifiedInfo.TokenSymbol, tokenHash); cl {
			unifiedInfo.IsClone = cl
			unifiedInfo.OriginalTokenAddress = orig
		}
	}

	if proof, by := social.Scan(tokenHash); proof {
		unifiedInfo.SocialProof = proof
		unifiedInfo.AnnouncedBy = by
	}

	fingerprint.Scan(tokenHash, unifiedInfo)

	heuristics.Analyze(unifiedInfo)

	// Compute risk score
	score.Compute(unifiedInfo)

	return unifiedInfo
}
