# Honeypot Simulation

Uses the Solana `simulateTransaction` RPC call to test buying and selling the token.
If selling fails in the simulation the token is marked as a honeypot.
The method needs a recent blockhash and may fail on limited RPC nodes.
