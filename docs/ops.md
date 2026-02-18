# Ops — MVG Bridge (Broker MVP)

## Runbooks (MVP)

### 1) Chain congestion / fee spikes
- Disable affected chain via config toggle
- Increase `min_amount_in` and/or widen min_out buffer

### 2) High refund rate
- Check min_out buffer too tight
- Check price source latency
- Check deposit credit latency (Mixin pending)

### 3) Stuck orders
- Identify stage:
  - tx detected but no pending mixin
  - pending but not credited
  - swap succeeded but withdraw pending
- Use admin tooling to requeue executor or mark manual review

## Metrics to track
- p50/p95 time to detect tx
- p50/p95 time pending→credited
- p50/p95 time to withdraw submit
- refund rate by reason (late, below min_out)
- gross margin per order

## Security
- Keep deploy key read-only
- Rotate API admin token
- Principle of least privilege for any cloud/provider keys

