package scoring

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/s-Amine/token-scan/token"
)

var weightFile = filepath.Join("scanners", "customrules", "scoring", "weights.json")

func loadWeights() map[string]int {
	data, err := ioutil.ReadFile(weightFile)
	if err != nil {
		return map[string]int{
			"NoRenounceNoSource": 20,
			"HighTax":            15,
			"MintFunction":       25,
			"TopHolder70":        30,
		}
	}
	var m map[string]int
	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]int{
			"NoRenounceNoSource": 20,
			"HighTax":            15,
			"MintFunction":       25,
			"TopHolder70":        30,
		}
	}
	return m
}

// Calculate computes the risk score and risk level.
func Calculate(info *token.TokenInfo) {
	weights := loadWeights()
	score := 0
	if info.CanTakeBackOwnership && !info.IsOpenSource {
		score += weights["NoRenounceNoSource"]
	}
	buyTax, _ := strconv.ParseFloat(info.BuyTax, 64)
	sellTax, _ := strconv.ParseFloat(info.SellTax, 64)
	if buyTax > 15 || sellTax > 15 {
		score += weights["HighTax"]
	}
	if info.IsMintable {
		score += weights["MintFunction"]
	} else {
		for _, f := range info.HasSensitiveFunctions {
			if f == "mint" {
				score += weights["MintFunction"]
				break
			}
		}
	}
	if info.TopHolderPercent > 70 {
		score += weights["TopHolder70"]
	}
	info.RiskScore = score
	switch {
	case score <= 30:
		info.RiskLevel = "SAFE"
	case score <= 70:
		info.RiskLevel = "WARNING"
	default:
		info.RiskLevel = "DANGEROUS"
	}
}
