# Risk Scoring

This module assigns a risk score between 0 and 100 based on several indicators.
The weights are stored in `scanners/customrules/scoring/weights.json` and can be tweaked.

## Rules

- No renounce & not open source: +20
- Taxes over 15%%: +15
- Mint function detected: +25
- Top holder > 70%%: +30

Risk level is derived from the score:

- 0-30 SAFE
- 31-70 WARNING
- 71-100 DANGEROUS
