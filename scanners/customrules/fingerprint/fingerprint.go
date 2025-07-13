package fingerprint

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Entry in blacklist
type blacklistEntry struct {
	Hash        string `json:"hash"`
	Description string `json:"description"`
}

// Scan retrieves the bytecode of a contract using the ETH RPC endpoint and
// compares the sha256 hash against a local blacklist. It returns matching
// descriptions.
func Scan(address string) ([]string, error) {
	rpc := os.Getenv("ETH_RPC_URL")
	if rpc == "" {
		return nil, fmt.Errorf("ETH_RPC_URL not set")
	}

	payload := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getCode","params":["%s","latest"],"id":1}`, address))
	res, err := http.Post(rpc, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	codeHash := sha256.Sum256([]byte(resp.Result))
	hashStr := hex.EncodeToString(codeHash[:])

	// Load blacklist
	path := filepath.Join("data", "blacklist_hashes.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []blacklistEntry
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	matches := []string{}
	for _, e := range list {
		if strings.EqualFold(hashStr, e.Hash) {
			matches = append(matches, e.Description)
		}
	}

	return matches, nil
}
