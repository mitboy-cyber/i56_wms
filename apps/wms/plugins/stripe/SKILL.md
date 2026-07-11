# Stripe Plugin

**Path**: `apps/wms/plugins/stripe/`
**Version**: 1.0.0
**Type**: Integration Plugin

## Purpose
Handles payment processing via Stripe for the I56 WMS.

## Features
- Process payment charges via Stripe Charges API
- List transaction history
- Registers navigation menu item under "Integrations"

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/plugins/stripe/charge` | Process a payment charge |
| GET | `/api/plugins/stripe/transactions` | List transaction history |

## Configuration
```json
{
  "api_key": "sk_test_xxx",
  "webhook_secret": "whsec_xxx",
  "currency": "usd"
}
```

## Dependencies
- Framework plugin system (`github.com/i56/framework/plugin`)
- Stripe API (external)
