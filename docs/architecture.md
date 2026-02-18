# Architecture — MVG Bridge (Broker MVP)

## 1. Components

- **API** (HTTP)
  - Create order
  - Get order status
  - Admin endpoints (limits, toggles)

- **Chain Watchers** (workers)
  - Detect deposits on supported chains
  - Attach txid + detected_at to order

- **Mixin Ingest** (worker)
  - Track pending deposits until credited
  - Track internal snapshots if needed

- **Executor** (worker)
  - On `deposit_credited`: compute `final_out`
  - If below `min_out`: refund
  - Else: swap → withdraw

- **Ledger / Reconciliation**
  - Store all external ids (txid, snapshot id, withdraw id)
  - Daily reconciliation report

## 2. Event Flow (happy path)

1) Order created → allocate deposit address
2) Chain watcher detects deposit tx → `deposit_tx_detected`
3) Mixin ingest sees pending deposit → `deposit_pending_mixin`
4) Mixin ingest sees credited → `deposit_credited`
5) Executor computes `final_out`:
   - if ok: swap + withdraw
   - else: refund
6) Status updated with txids

## 3. Idempotency

- Handlers must be idempotent:
  - chain tx events keyed by (chain, txid)
  - mixin deposit events keyed by deposit id
  - swap/withdraw actions keyed by order id

Use optimistic locking on `orders.version` to avoid double execution.

## 4. Observability

- Structured logs with `order_public_id` and `trace_id`
- Metrics:
  - orders created
  - deposit detected latency
  - pending→credited latency
  - completion time distribution
  - refund rate and reasons

## 5. Security

- Keep deploy keys private and separate
- Separate hot-wallet operations from API surface where possible
- Rate limit public endpoints

