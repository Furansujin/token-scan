# Source Check

This module retrieves verified Solidity source code from Etherscan and searches for sensitive functions.

## Usage

```
./token-scan -mode sourcecheck -token <address>
```

Environment variable `ETHERSCAN_API_KEY` can be used for higher rate limits.

The following keywords are searched: `mint`, `setFee`, `blacklist`, `setMaxTxAmount`, `transferOwnership`.
If the contract is not verified on Etherscan the module returns an error.
