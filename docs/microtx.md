# Micro Transactions

Counts very small buy operations (less than 0.005 SOL) to detect bot activity.
More than twenty such transactions within two minutes trigger `BotActivity`.
These heuristics are indicative only and may produce false positives.
