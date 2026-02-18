# State Machine — MVG Bridge (Broker MVP)

## States

- `quote_created`
- `awaiting_deposit`
- `deposit_tx_detected`
- `deposit_pending_mixin`
- `deposit_credited`
- `executing_swap`
- `withdrawing`
- `completed`

Refund path:
- `refunding`
- `refunded`

Manual:
- `failed_manual_review`

## Transition rules (core)

1) awaiting_deposit → deposit_tx_detected
- when chain watcher detects deposit tx for this order

2) deposit_tx_detected → refunding
- if detected late: `deposit_tx_detected_at > created_at + pay_window_seconds`

3) deposit_tx_detected → deposit_pending_mixin
- when mixin shows pending deposit corresponding to tx/order

4) deposit_pending_mixin → deposit_credited
- when mixin credits balance

5) deposit_credited → refunding
- if computed `final_out < min_out`

6) deposit_credited → executing_swap
- else

7) executing_swap → withdrawing
- swap ok

8) withdrawing → completed
- withdraw submitted (txid available)

## Notes

- If any step fails transiently, worker retries must be idempotent.
- Completion = withdraw submitted (MVP). Chain confirmation can be tracked asynchronously.

