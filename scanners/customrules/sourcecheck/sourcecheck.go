package sourcecheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// SourceResponse represents the structure returned by Etherscan API.
type SourceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		SourceCode string `json:"SourceCode"`
	} `json:"result"`
}

var keywords = []string{"mint", "setFee", "blacklist", "setMaxTxAmount", "transferOwnership"}

// Scan retrieves the contract source code from Etherscan and detects sensitive functions.
func Scan(address string) ([]string, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s", address)
	if apiKey != "" {
		url += "&apikey=" + apiKey
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r SourceResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	if r.Status != "1" || len(r.Result) == 0 {
		return nil, fmt.Errorf("contract not verified")
	}
	code := strings.ToLower(r.Result[0].SourceCode)
	var detected []string
	for _, k := range keywords {
		if strings.Contains(code, strings.ToLower(k)) {
			detected = append(detected, k)
		}
	}
	return detected, nil
}
