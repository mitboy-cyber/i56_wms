# FedEx Plugin

**Path**: `apps/wms/plugins/fedex/`
**Version**: 1.0.0
**Type**: Integration Plugin

## Purpose
Provides FedEx tracking and shipping integration for the I56 WMS.

## Features
- Track packages by tracking number via HTTP API
- Create FedEx shipments
- Registers navigation menu item under "Integrations"

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/plugins/fedex/track/:tracking` | Track a package by tracking number |
| POST | `/api/plugins/fedex/ship` | Create a new FedEx shipment |

## Configuration
```json
{
  "account_number": "123456789",
  "api_key": "l1234567890abcdef",
  "api_secret": "abcdef1234567890",
  "meter_number": "12345678"
}
```

## Dependencies
- Framework plugin system (`github.com/i56/framework/plugin`)
- FedEx REST API (external)
