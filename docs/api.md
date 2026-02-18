# API — MVG Bridge (Broker MVP)

This document specifies the external HTTP API (public) and a minimal admin API.

## Conventions

- All numeric amounts are **strings**.
- `public_id` is the stable identifier exposed to users.

## 1) Create Order

`POST /v1/orders`

Request
```json
{
  "source_chain": "ETH",
  "source_asset": "USDT",
  "amount_in": "100",
  "target_chain": "TRON",
  "target_asset": "USDT",
  "target_address": "T..."
}
```

Response
```json
{
  "public_id": "BRG_9K3F...",
  "status": "awaiting_deposit",
  "pay_window_seconds": 900,
  "deposit": {
    "chain": "ETH",
    "asset": "USDT",
    "address": "0x...",
    "amount": "100"
  },
  "quote": {
    "estimated_out": "99.7",
    "min_out": "99.2",
    "quote_expiry_seconds": 60
  },
  "terms": {
    "late_deposit": "auto_refund",
    "below_min_out": "auto_refund",
    "refund_fee": "paid_by_user",
    "refund_to": "original_address"
  }
}
```

## 2) Get Order

`GET /v1/orders/{public_id}`

Response includes:
- status
- detected txid + timestamps
- credited timestamp + credited amount
- swap/withdraw txids
- next user instructions

## 3) Admin (minimal)

- `POST /admin/chains/{chain}/toggle` enable/disable
- `POST /admin/limits` update per-tier limits

Auth: static token in header for MVP.

