package clone

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

// Scan compares the given token with a whitelist and returns true if it is a clone.
func Scan(name, symbol, address string) (bool, string) {
	data, err := ioutil.ReadFile(filepath.Join("data", "token_whitelist.json"))
	if err != nil {
		return false, ""
	}
	var m map[string]string
	_ = json.Unmarshal(data, &m)
	key := name
	if v, ok := m[key]; ok && v != address {
		return true, v
	}
	return false, ""
}
