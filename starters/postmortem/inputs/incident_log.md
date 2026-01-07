# Incident Timeline: API 500 Outage

## Timeline (UTC)
- **10:00**: Deployment of `checkout-service:v2.4.0` begins.
- **10:03**: Monitoring alerts "High Error Rate" on `/checkout` endpoint.
- **10:05**: On-call engineer (Alice) acknowledges alert.
- **10:07**: Alice attempts to rollback to `v2.3.0`.
- **10:10**: Rollback failed due to database schema incompatibility.
- **10:15**: Alice enables "Maintenance Mode" for checkout.
- **10:30**: Root cause identified: `v2.4.0` expected a column `user_id` which was postponed.
- **10:45**: Database migration applied manually.
- **10:50**: Maintenance mode disabled. Traffic normal.

## Logs Excerpt
```
[ERROR] 10:02:45 app.checkout: Query failed: column "user_id" does not exist
[ERROR] 10:02:46 app.checkout: Query failed: column "user_id" does not exist
...
```
