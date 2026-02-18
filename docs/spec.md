# Specification — MVG Bridge (Broker MVP)

> Battle-tested assumptions:
> - This is a **broker** service, not a trustless atomic bridge.
> - Reduce dispute surface area by enforcing **strict deterministic rules**.

---

## 1. Goals

- Provide a bridge-like UX for cross-chain value movement using Mixin swap + withdrawals.
- Support external (non-Mixin Messenger) users via on-chain deposits.
- Keep the system automatable and operationally sane.

## 2. Non-goals (MVP)

- Trustless atomic cross-chain bridging
- Multi-party custody / MPC / proof of reserves
- Advanced AML/KYC (only basic screening + limits)
- Supporting every chain on day 1

## 3. Product rules (MVP — fixed)

### 3.1 Pricing and execution

- Orders are **Exact-in**: user pays a fixed `amount_in`.
- On order creation, show:
  - `estimated_out` (estimate)
  - `min_out` (guarantee threshold)
  - `pay_window_seconds` (default 900)

**Final execution time:**
- `final_out` is computed at `deposit_credited_at` (when funds are credited to Mixin balance).

**Execution decision:**
- If `final_out >= min_out` ⇒ execute swap + withdraw
- If `final_out < min_out` ⇒ auto-refund (no user override)

### 3.2 Payment window and lateness

- Default payment window: **15 minutes**, with per-chain override.

**Lateness decision point (chosen):**
- Lateness is decided using `deposit_tx_detected_at` (the first time our system detects the deposit tx), not chain block timestamp.

**Rules:**
- If `deposit_tx_detected_at > created_at + pay_window_seconds` ⇒ **late deposit** ⇒ auto-refund
- If detected within window but Mixin credit is slow ⇒ treat as on-time (not user fault)

### 3.3 Refund policy

- Refund destination: **original deposit from-address** (source tx `from`).
- Refund network fee: **paid by user**.
- No custom refund address.

### 3.4 Deposit confirmation

- Chain tx detected is informational.
- Only `deposit_credited` (credited to Mixin) counts as “paid”.

States:
- `deposit_tx_detected` → `deposit_pending_mixin` → `deposit_credited`

### 3.5 Time estimate (UI)

- Show “estimated confirmation time” as informational only.
- Do not promise a fixed SLA.

## 4. Chain/asset support strategy

- Use tiering:
  - Tier-1: 2–4 chains + stablecoins
  - Expand by config toggles and limits

## 5. Abuse & risk controls (minimal)

- New-user limits (per order/day)
- Rate limiting (IP/device/address)
- Basic sanctions/blacklist screening on withdraw address
- Emergency kill-switch per chain/asset

---

## 6. Open decisions (to be confirmed)

- Source of `deposit_tx_detected_at` (system clock is acceptable for MVP)
- Completion definition: withdraw tx submitted vs confirmed
- How to charge refund fee: deduct from amount vs separate fee

