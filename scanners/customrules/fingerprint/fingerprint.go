package fingerprint

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/s-Amine/token-scan/token"
)

var blacklistPath = filepath.Join("data", "blacklist_hashes.json")

func loadBlacklist() []string {
	data, err := ioutil.ReadFile(blacklistPath)
	if err != nil {
		return nil
	}
	var list []string
	_ = json.Unmarshal(data, &list)
	return list
}

// Scan compares contract bytecode hash with blacklist.
func Scan(address string, info *token.TokenInfo) {
	rpc := os.Getenv("RPC_ENDPOINT")
	if rpc == "" {
		rpc = "https://rpc.ankr.com/eth"
	}
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getCode",
		"params":  []string{address, "latest"},
		"id":      1,
	}
	data, _ := json.Marshal(payload)
	resp, err := http.Post(rpc, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}
	hash := sha256.Sum256([]byte(result.Result))
	h := hex.EncodeToString(hash[:])
	bl := loadBlacklist()
	for _, b := range bl {
		if h == b {
			info.MatchedScamContracts = append(info.MatchedScamContracts, b)
			break
		}
	}
}
