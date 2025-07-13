package sourcecheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Scan retrieves verified source code from Etherscan and searches for sensitive
// functions. It returns a list of detected function names.
func Scan(address string, apiKey string) ([]string, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=contract&action=getsourcecode&address=%s", address)
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
			SourceCode string `json:"SourceCode"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	if len(resp.Result) == 0 || resp.Result[0].SourceCode == "" {
		return nil, fmt.Errorf("contract not verified")
	}

	code := strings.ToLower(resp.Result[0].SourceCode)
	keywords := []string{"mint", "setfee", "blacklist", "setmaxtxamount", "transferownership"}
	detected := []string{}
	for _, k := range keywords {
		if strings.Contains(code, k) {
			detected = append(detected, k)
		}
	}
	return detected, nil
}
