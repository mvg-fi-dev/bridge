# MVG Bridge (Broker MVP)

This repository implements a **bridge-like** cross-chain swap experience using Mixin as the execution rail.

Important: this is **not** a trustless on-chain bridge. It is a **broker service** (custodial in the execution flow) with strict, explicit rules to minimize disputes and risk.

## What it does

- User sends funds on a source chain to a per-order deposit address
- Service waits for the deposit to be **credited in Mixin**
- Service performs swap via Mixin
- Service withdraws to the user’s target-chain address
- If the minimum-out condition is not met, or if the user pays late: **auto-refund**

## Key rules (MVP)

- **Exact-in + min_out**
- Compute final output at **Mixin credited time**
- If `final_out < min_out` ⇒ **auto-refund (return to original deposit address)**
- Refund network fee ⇒ **paid by user**
- Payment window default **15 minutes** (per-chain override)
- Late deposit ⇒ **auto-refund**
- No “force market fill” toggle

See full specification in `docs/`.
