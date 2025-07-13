package holders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type holder struct {
	Address string `json:"HolderAddress"`
	Balance string `json:"TokenHolderQuantity"`
}

type response struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []holder `json:"result"`
}

// Scan retrieves top token holders percentage and returns the highest percentage and a map of top holders.
func Scan(address string) (float64, map[string]float64, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	url := fmt.Sprintf("https://api.etherscan.io/api?module=token&action=tokenholderlist&contractaddress=%s&page=1&offset=10", address)
	if apiKey != "" {
		url += "&apikey=" + apiKey
	}
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	var r response
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, nil, err
	}
	if r.Status != "1" {
		return 0, nil, fmt.Errorf("holders fetch failed")
	}
	total := 0.0
	for _, h := range r.Result {
		b, _ := strconv.ParseFloat(h.Balance, 64)
		total += b
	}
	topMap := make(map[string]float64)
	var maxPercent float64
	for i, h := range r.Result {
		b, _ := strconv.ParseFloat(h.Balance, 64)
		percent := 0.0
		if total > 0 {
			percent = (b / total) * 100
		}
		if i < 5 {
			topMap[h.Address] = percent
		}
		if percent > maxPercent {
			maxPercent = percent
		}
	}
	return maxPercent, topMap, nil
}
