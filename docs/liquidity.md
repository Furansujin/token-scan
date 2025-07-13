# Liquidity Lock Verification

Attempts to identify whether LP tokens are locked in known Solana lockers.
If the LP is held by the deployer or unlock time is less than 30 days, the token is flagged as risky.
Information can be incomplete when lockers are not publicly queryable.
