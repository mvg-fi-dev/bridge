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

5) deposit_credited → executing_swap
- submit swap execution (ExinSwap transfer with memo)

6) executing_swap → withdrawing
- swap ok (ExinSwap pays out target asset to our bot; reconcile by server memo TRACE)

7) executing_swap → refunding
- swap failed or refunded (ExinSwap refunds input asset to our bot; reconcile by server memo TRACE)

8) withdrawing → completed
- withdraw submitted (txid available)

9) refunding → refunded
- refund transfer submitted back to original payer (Mixin internal transfer)

## Notes

- If any step fails transiently, worker retries must be idempotent.
- Completion = withdraw submitted (MVP). Chain confirmation can be tracked asynchronously.

