# Local testing (Route quote/swap)

## ⚠️ Secrets

Do **not** put real mnemonics/keys in git.

- Use `.env` (ignored by git) or your shell env vars.
- If a mnemonic ever appeared in chat/logs, treat it as compromised.

## 1) Route quote test

From repo root:

```bash
export MIXIN_KEYSTORE_PATH=/path/to/keystore-xxxx.json
export TEST_MIXIN_ACCOUNT_ID=<execution-account-id>
export TEST_MIXIN_MNEMONIC="<mnemonic>"
export TEST_INPUT_MINT=<from.assetId>
export TEST_OUTPUT_MINT=<to.assetId>
export TEST_AMOUNT=1000000

go run ./cmd/route-quote-test
```

Expected output:
- prints route bot public key (base64url)
- prints `quote outAmount=... payload_len=...`

## 2) Troubleshooting

- `route quote status=401/403`:
  - account id does not match mnemonic-derived account
  - MR signing mismatch
  - route service may have additional anti-abuse headers

- `invalid mnemonic`:
  - use the correct Messenger account mnemonic (not a web3 wallet mnemonic)

