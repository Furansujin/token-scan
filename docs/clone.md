# Clone Detection

Token names and symbols are compared with a whitelist stored in `data/token_whitelist.json`.
If a known name is used with a different address, the token is flagged as a clone.
Generic names can lead to false alerts.
