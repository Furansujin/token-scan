package holders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// HolderResult represents the percentage of holdings by address.
type HolderResult struct {
	Address string
	Percent float64
}

// Scan queries Etherscan for token holder information and returns the top
// holders sorted by balance percentage.
func Scan(address string, apiKey string) ([]HolderResult, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=token&action=tokenholderlist&contractaddress=%s&page=1&offset=10", address)
	if apiKey != "" {
		url += "&apikey=" + apiKey
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result []struct {
			HolderAddress string  `json:"HolderAddress"`
			HolderPct     float64 `json:"HolderPct"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	results := make([]HolderResult, len(resp.Result))
	for i, r := range resp.Result {
		results[i] = HolderResult{Address: r.HolderAddress, Percent: r.HolderPct}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Percent > results[j].Percent })
	return results, nil
}
