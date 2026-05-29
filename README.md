# Dari Coffee CLI

Public CLI for ordering Dari Coffee in San Francisco.

## Install

Install the latest native macOS/Linux binary:

```bash
curl -fsSL https://raw.githubusercontent.com/mupt-ai/dari-coffee-cli/main/install.sh | bash
dari-coffee --help
```

## Commands

```bash
dari-coffee version
dari-coffee menu
dari-coffee menu --json
dari-coffee service
dari-coffee check-address --json "1 Market St, San Francisco, CA 94105"
dari-coffee order --name "Ada Lovelace" --email "ada@example.com" --company "Analytical Engines Inc." --address "1 Market St, San Francisco, CA 94105" --items-json '[{"shop_slug":"starbucks","drink_slug":"americano","size":"grande","quantity":1}]'
```

`menu --json` is the source of truth for shop, drink, and size identifiers.
Each order can contain multiple drinks, but all drinks must come from one shop.

`order` prints a Stripe Checkout URL. Checkout authorizes payment only; Dari
captures only after accepting the order.

## Agent Guide

Agents should read [llms.txt](llms.txt) before ordering for a user.

The public CLI defaults to:

```text
https://coffee.dari.dev
```

## Development

The public API contract lives in [api/openapi.public.yaml](api/openapi.public.yaml).
Generated Go types live in `internal/openapi/`.

After changing the public API contract, run:

```bash
go generate ./internal/openapi
go test ./...
```
