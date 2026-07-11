# Shopify Plugin

**Path**: `apps/wms/plugins/shopify/`
**Version**: 1.0.0
**Type**: Integration Plugin

## Purpose
Imports orders from Shopify into the I56 Warehouse Management System (WMS).

## Features
- List imported Shopify orders via HTTP API
- Import orders from Shopify Admin API (webhook & polling)
- Registers navigation menu item under "Integrations"

## API Endpoints
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/plugins/shopify/orders` | List imported Shopify orders |
| POST | `/api/plugins/shopify/import` | Trigger order import from Shopify |

## Configuration
```json
{
  "store": "my-store.myshopify.com",
  "api_key": "shpat_xxx",
  "api_secret": "shpss_xxx",
  "webhook_secret": "whsec_xxx"
}
```

## Dependencies
- Framework plugin system (`github.com/i56/framework/plugin`)
- Shopify Admin API (external)
