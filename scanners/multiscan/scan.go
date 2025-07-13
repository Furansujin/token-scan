package multiscan

import (
	"github.com/s-Amine/token-scan/scanners/customrules/fingerprint"
	"github.com/s-Amine/token-scan/scanners/customrules/heuristics"
	"github.com/s-Amine/token-scan/scanners/customrules/holders"
	"github.com/s-Amine/token-scan/scanners/customrules/scoring"
	"github.com/s-Amine/token-scan/scanners/customrules/sourcecheck"
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
	holdersChan := make(chan []holders.HolderResult)
	sourceChan := make(chan []string)
	fingerprintChan := make(chan []string)

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

	// Holders analysis
	go func() {
		res, _ := holders.Scan(tokenHash, "")
		holdersChan <- res
	}()

	// Source code check
	go func() {
		funcs, _ := sourcecheck.Scan(tokenHash, "")
		sourceChan <- funcs
	}()

	// Fingerprint
	go func() {
		matches, _ := fingerprint.Scan(tokenHash)
		fingerprintChan <- matches
	}()

	// Variables to store scan results
	var goPlusTokenInfo *token.TokenInfo
	var honeypotTokenInfo *token.TokenInfo
	var quickIntelTokenInfo *token.TokenInfo
	var holderResults []holders.HolderResult
	var sourceFuncs []string
	var fingerprintMatches []string

	// Receive scan results from channels
	for i := 0; i < 6; i++ {
		select {
		case goPlusTokenInfo = <-goPlusScanResultChan:
		case honeypotTokenInfo = <-isHoneypotScanResultChan:
		case quickIntelTokenInfo = <-quickIntelScanResultChan:
		case holderResults = <-holdersChan:
		case sourceFuncs = <-sourceChan:
		case fingerprintMatches = <-fingerprintChan:
		}
	}

	// Unify scan results into one TokenInfo
	unifiedInfo := token.UnifyTokenInfo(goPlusTokenInfo, honeypotTokenInfo, quickIntelTokenInfo)
	if len(holderResults) > 0 {
		unifiedInfo.TopHolderPercent = holderResults[0].Percent
		unifiedInfo.TopHolders = map[string]float64{}
		for _, h := range holderResults {
			unifiedInfo.TopHolders[h.Address] = h.Percent
		}
	}
	unifiedInfo.HasSensitiveFunctions = sourceFuncs
	unifiedInfo.MatchedScamContracts = fingerprintMatches
	unifiedInfo.DetectedPatterns = heuristics.Detect(unifiedInfo)
	scoring.Calculate(unifiedInfo)

	return unifiedInfo
}
