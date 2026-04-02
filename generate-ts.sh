#!/usr/bin/env bash
set -euo pipefail
swagger-typescript-api generate --path docs/swagger.yaml --output frontend/src/lib/types --name api.ts
