# ExinSwap V2 (reference)

Source: https://github.com/ExinOne/exinswap-doc (README)

## Bot trading via memo

- ExinSwap bot (Mixin ID): 7000102352
- ExinSwap bot (User ID): 29f23576-4651-47ff-8c16-6c8a5d76985e

### User transfer memo format (BASE64)

`ACTION$FIELD1$FIELD2$FIELD3$FIELD4`

Trade action (ACTION=0):
- FIELD1: target asset UUID
- FIELD2: min_out (optional; 0 means no refund on slippage)
- FIELD3: latest execution unix time seconds (optional)
- FIELD4: route string (optional)

This matches our broker design nicely:
- We can set FIELD2 = min_out
- Set FIELD3 = now + TTL
- Set FIELD4 if we want a specific route; otherwise omit.

### API

Endpoint: https://app.exinswap.com/api/v2
- GET /pairs
- GET /assets

