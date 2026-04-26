# Directory Structure

```text
.
+-- backend/
|   +-- api/
|   |   +-- typespec/          # TypeSpec source of truth for API contract
|   |   +-- openapi/           # Generated OpenAPI documents
|   +-- cmd/
|   |   +-- api/               # API server entrypoint
|   +-- internal/
|       +-- auth/              # Discord OAuth2 and JWT
|       +-- config/            # Environment configuration
|       +-- generated/         # oapi-codegen generated Go code
|       +-- presence/          # Presence domain logic and handlers
|       +-- server/            # HTTP server wiring for generated API interface
|       +-- store/
|           +-- redis/         # Redis client and persistence
+-- frontend/
|   +-- app/                   # Next.js App Router pages/layouts
|   +-- components/            # Reusable UI components
|   +-- lib/                   # API client and browser utilities
|   +-- public/                # PWA assets
+-- deployments/               # Deploy config such as Docker/Fly/Render
+-- docs/                      # Project notes and design docs
+-- CODEX_PROMPT.md            # Implementation prompt for Codex
+-- README.md                  # MVP requirements
+-- go.mod                     # Go module definition
```
