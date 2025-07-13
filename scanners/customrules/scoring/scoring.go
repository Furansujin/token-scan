package scoring

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/s-Amine/token-scan/token"
)

// Weights defines scoring weights loaded from JSON configuration.
type Weights struct {
	NonRenouncedNotOpenSource int `json:"non_renounced_not_open_source"`
	HighTax                   int `json:"high_tax"`
	MintFunction              int `json:"mint_function"`
	TopHolderOver70           int `json:"top_holder_over_70"`
}

// defaultWeights provides fallbacks when configuration cannot be loaded.
var defaultWeights = Weights{
	NonRenouncedNotOpenSource: 20,
	HighTax:                   15,
	MintFunction:              25,
	TopHolderOver70:           30,
}

// loadWeights attempts to read weights from weights.json. When the file cannot
// be read or parsed, default values are used.
func loadWeights() Weights {
	path := filepath.Join("scanners", "customrules", "scoring", "weights.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return defaultWeights
	}
	// Remove comment lines starting with '//' before unmarshalling
	lines := bytes.Split(data, []byte("\n"))
	cleaned := make([]byte, 0, len(data))
	for _, l := range lines {
		t := strings.TrimSpace(string(l))
		if strings.HasPrefix(t, "//") || t == "" {
			continue
		}
		cleaned = append(cleaned, l...)
		cleaned = append(cleaned, '\n')
	}

	var w Weights
	if err := json.Unmarshal(cleaned, &w); err != nil {
		return defaultWeights
	}
	return w
}

// parseTax converts a string tax value to integer percentage.
func parseTax(tax string) int {
	if tax == "" {
		return 0
	}
	if i, err := strconv.Atoi(tax); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(tax, 64); err == nil {
		return int(f)
	}
	return 0
}

// Calculate fills RiskScore and RiskLevel fields of the provided TokenInfo
// based on configured weights.
func Calculate(info *token.TokenInfo) {
	if info == nil {
		return
	}
	w := loadWeights()
	score := 0

	if !info.ContractRenounced && !info.IsOpenSource {
		score += w.NonRenouncedNotOpenSource
	}

	buyTax := parseTax(info.BuyTax)
	sellTax := parseTax(info.SellTax)
	if buyTax > 15 || sellTax > 15 {
		score += w.HighTax
	}

	if info.IsMintable {
		score += w.MintFunction
	}

	if info.TopHolderPercent > 70 {
		score += w.TopHolderOver70
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
