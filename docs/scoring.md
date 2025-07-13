# Global Scoring

The scoring module aggregates data from all checks to provide a risk score between 0 and 100. Depending on this score the risk level is **SAFE**, **WARNING** or **DANGEROUS**.

| Condition | Points |
|-----------|-------|
| Non renounced and not open-source | 20 |
| Taxes > 15% | 15 |
| Mint function detected | 25 |
| 1 wallet > 70% supply | 30 |

A total score above 80 sets the level to **DANGEROUS**.

## Example

```
./token-scan -mode multiscan -token <address> --pretty
```

## Limitations

Weights are configurable in `scanners/customrules/scoring/weights.json`.
