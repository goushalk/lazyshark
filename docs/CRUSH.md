# CRUSH.md

## Go Commands

- **Build:** `go build ./cmd/packetbreather++`
- **Run:** `go run ./cmd/packetbreather++`
- **Test:** `go test ./...`
- **Test (single package):** `go test ./internal/analyzer`
- **Lint/Vet:** `go vet ./...`

## Code Style Guidelines (Go)

- **Formatting:** Use `go fmt ./...` for consistent code formatting.
- **Imports:** Group standard library imports separately from third-party imports.
- **Naming Conventions:**
    - Package names: short, all lowercase, no underscores (e.g., `analyzer`, `tui`).
    - Variable names: `camelCase`. Acronyms in names should be all caps (e.g., `HTTPClient`).
    - Function names: `CamelCase` for exported functions, `camelCase` for unexported functions.
- **Error Handling:** Return errors explicitly as the last return value. Check errors immediately. Do not ignore errors.
- **Types:** Use specific types over `interface{}` when possible. Define custom types for clarity and type safety.
- **Dependencies:** Manage dependencies with `go mod`. Ensure all external modules are properly declared in `go.mod` and kept up-to-date. Key dependencies include `github.com/google/gopacket` and `github.com/charmbracelet/bubbletea`.
