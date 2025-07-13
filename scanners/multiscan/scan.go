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

	// Additional analyses run concurrently
	type holdersRes struct {
		pct float64
		top map[string]float64
	}
	type liquidityRes struct {
		locked bool
		unlock int64
	}
	type honeypotRes struct {
		isH   bool
		errSt string
	}
	type microtxRes struct {
		count int
		bot   bool
	}
	type cloneRes struct {
		isClone bool
		orig    string
	}
	type socialRes struct {
		proof bool
		by    string
	}

	holdersChan := make(chan holdersRes)
	sourceCheckChan := make(chan []string)
	washTradeChan := make(chan float64)
	liquidityChan := make(chan liquidityRes)
	honeypotChan := make(chan honeypotRes)
	microtxChan := make(chan microtxRes)
	cloneChan := make(chan cloneRes)
	socialChan := make(chan socialRes)

	go func() {
		funcs, err := sourcecheck.Scan(tokenHash)
		if err == nil {
			sourceCheckChan <- funcs
		} else {
			sourceCheckChan <- nil
		}
	}()

	go func() {
		pct, top, err := holders.Scan(tokenHash)
		if err == nil {
			holdersChan <- holdersRes{pct: pct, top: top}
		} else {
			holdersChan <- holdersRes{}
		}
	}()

	go func() {
		s, _ := washtrade.Scan(tokenHash)
		washTradeChan <- s
	}()

	go func() {
		locked, unlock, _ := liquidity.Scan(tokenHash)
		liquidityChan <- liquidityRes{locked: locked, unlock: unlock}
	}()

	go func() {
		isH, errStr := honeypot.Scan(tokenHash)
		honeypotChan <- honeypotRes{isH: isH, errSt: errStr}
	}()

	go func() {
		count, bot := microtx.Scan(tokenHash)
		microtxChan <- microtxRes{count: count, bot: bot}
	}()

	go func() {
		if unifiedInfo.TokenName != "" {
			cl, orig := clone.Scan(unifiedInfo.TokenName, unifiedInfo.TokenSymbol, tokenHash)
			cloneChan <- cloneRes{isClone: cl, orig: orig}
		} else {
			cloneChan <- cloneRes{}
		}
	}()

	go func() {
		proof, by := social.Scan(tokenHash)
		socialChan <- socialRes{proof: proof, by: by}
	}()

	// Receive results
	for i := 0; i < 8; i++ {
		select {
		case funcs := <-sourceCheckChan:
			if funcs != nil {
				unifiedInfo.HasSensitiveFunctions = funcs
			}
		case h := <-holdersChan:
			unifiedInfo.TopHolderPercent = h.pct
			unifiedInfo.TopHolders = h.top
		case s := <-washTradeChan:
			unifiedInfo.WashTradeScore = s
		case l := <-liquidityChan:
			unifiedInfo.LiquidityLocked = l.locked
			unifiedInfo.LiquidityUnlockTime = l.unlock
		case h := <-honeypotChan:
			unifiedInfo.IsHoneypot = h.isH
			unifiedInfo.HoneypotError = h.errSt
		case m := <-microtxChan:
			unifiedInfo.MicroTxCount = m.count
			unifiedInfo.BotActivity = m.bot
		case c := <-cloneChan:
			unifiedInfo.IsClone = c.isClone
			unifiedInfo.OriginalTokenAddress = c.orig
		case s := <-socialChan:
			unifiedInfo.SocialProof = s.proof
			unifiedInfo.AnnouncedBy = s.by
		}
	}

	fingerprint.Scan(tokenHash, unifiedInfo)

	heuristics.Analyze(unifiedInfo)

	// Compute risk score
	score.Compute(unifiedInfo)

	return unifiedInfo
}
