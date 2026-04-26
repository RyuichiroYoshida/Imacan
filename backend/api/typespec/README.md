# ProjectImacan API Contract

This directory contains the TypeSpec source of truth for the ProjectImacan API.

## Files

- `main.tsp`: MVP API contract for auth and presence.
- `tspconfig.yaml`: TypeSpec compiler configuration.
- `package.json`: Local TypeSpec tooling scripts.

## Generate OpenAPI

```bash
npm install
npm run build
```

The OpenAPI document is generated at:

```text
backend/api/openapi/project-imacan.openapi.yaml
```

The generated OpenAPI file is intended to be the input for `oapi-codegen`.
