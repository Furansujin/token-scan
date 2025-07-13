# Source Code Check

This module queries the Etherscan API to retrieve verified Solidity source code for a contract. It looks for sensitive functions such as `mint`, `setFee`, `blacklist`, `setMaxTxAmount` and `transferOwnership`.

## CLI Example

```
./token-scan -mode sourcecheck -token <address>
```

## API

- [Etherscan getsourcecode](https://api.etherscan.io/api?module=contract&action=getsourcecode&address=...)

## Limitations

- Requires an Etherscan API key for heavy usage.
- Fails if the contract is not verified on Etherscan.
