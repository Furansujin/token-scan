# Wash Trading Detection

This module checks transfer history for repetitive loops between a small set of addresses.
It returns a `WashTradeScore` between 0 and 1. Values above `0.2` indicate suspicious activity.
The heuristic relies on public transaction data and may miss behaviour using new anonymous wallets.
