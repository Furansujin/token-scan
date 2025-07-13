# Heuristic Analysis

This module combines several conditions collected during a scan to detect suspicious patterns, such as potential rug pulls.

## CLI Example

```
./token-scan -mode multiscan -token <address> --pretty
```

## Patterns
- **Possible rug pull**: contract not open-source, honeypot and mintable.
- **Centralized ownership**: a holder owns more than 50%% and contract not renounced.
- **Blacklisting and pausable transfers**: token can blacklist and pause transfers.

## Limitations

These rules are heuristic and might generate false positives.
