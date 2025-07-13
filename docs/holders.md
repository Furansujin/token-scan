# Holders Analysis

Queries Etherscan for the list of token holders and reports the percentage distribution.
If a single holder owns more than 50%% of the returned supply, a warning is raised.

## Usage
```
./token-scan -mode multiscan -token <address> --pretty
```

This module requires `ETHERSCAN_API_KEY` for reliable results.
