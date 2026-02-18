# Notes: Mixin Swap API (reverse-engineered from Android app)

This is **not** yet implemented in the Go service. It documents what the official Android client appears to call.

## Evidence

From `MixinNetwork/android-app`:
- `RouteService.kt` defines:
  - `GET web3/quote`
  - `POST web3/swap`
  - other endpoints: `web3/tokens`, `web3/swap/orders`, ...
- Base URL for RouteService: `Constants.RouteConfig.ROUTE_BOT_URL = https://api.route.mixin.one`

## Quote

Request (query params):
- `inputMint`
- `outputMint`
- `amount`
- `source` (default `web3`)

Response (`QuoteResult`):
- `inputMint`, `inAmount`
- `outputMint`, `outAmount`
- `slippage` (int)
- `source`
- `payload` (string)

The returned `payload` is then sent back in the swap request.

## Swap

Request body (`SwapRequest`):
- `payer` (string)
- `inputMint` (string)
- `inputAmount` (string)
- `outputMint` (string)
- `payload` (string)
- `source` (string)
- `withdrawalDestination` (nullable string)
- `referral` (nullable string)
- `walletId` (nullable string)

Response type: `SwapResponse` (not inspected yet).

## Authentication / headers

Android app adds custom headers for RouteService calls:
- `MR-ACCESS-TIMESTAMP`
- `MR-ACCESS-SIGN`
- plus `Mixin-Device-Id`, `X-Request-Id`, etc.

`MR-ACCESS-SIGN` is an HMAC constructed from a shared key derived via ECDH between:
- user’s Ed25519 keypair converted to Curve25519
- route bot public key (fetched from sessions)

This suggests Route API is intended for **Messenger clients**, not bots.

## Next step

We need to decide how the broker service (bot backend) should authenticate to swap endpoints:
- Is there a bot/server auth mode (JWT) for `api.route.mixin.one`?
- Or should we execute swaps by interacting with a swap bot inside Mixin (transfer-in + transfer-out)?

