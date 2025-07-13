# Contract Fingerprint

This check fetches the on-chain bytecode of a contract via `eth_getCode`, hashes it with SHA256 and compares the hash against a local blacklist.

## CLI Example

```
./token-scan -mode fingerprint -token <address>
```

## Limitations

Requires access to an Ethereum RPC node through the `ETH_RPC_URL` environment variable. The blacklist is stored in `data/blacklist_hashes.json` and must be updated manually.
