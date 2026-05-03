# AGENTS.md

> **Living document** — update when tooling, CI, or conventions change. Trim stale lines; see Changelog.

## Commands

- Dev run: `go run .` (needs CGO + OS GUI deps; see README)
- Mod tidy: `go mod tidy`
- Vet / test / build: `CGO_ENABLED=1 go vet ./...`, `CGO_ENABLED=1 go test ./...`, `CGO_ENABLED=1 go build -trimpath .`
- Release-style binary (local): `go build -trimpath -ldflags "-s -w -X main.VERSION=$(git describe --tags --always)" .`
- GoReleaser config check: `goreleaser check` (install from [goreleaser.com](https://goreleaser.com/) or use CI)
- Single-target snapshot build: `goreleaser build --clean --snapshot --single-target`

Linux desktop deps (Debian/Ubuntu) match [.github/workflows/ci.yml](.github/workflows/ci.yml) (`apt-get install` line for GL/X11).

### gopls / VS Code: “build constraints exclude all Go files” (fyne, go-gl)

- Set **`CGO_ENABLED=1`** for the Go tools (workspace [.vscode/settings.json](.vscode/settings.json) does this).
- Do **not** leave **`GOOS` / `GOARCH`** set for cross-compilation while editing (e.g. `windows` / `amd64`). Run `unset GOOS GOARCH`, reload the window, **Go: Restart Language Server**.
- Confirm with `go env GOOS GOARCH CGO_ENABLED` — `GOOS`/`GOARCH` should match this machine or be empty defaults.

## Stack

- Go **1.25+**, **Fyne v2** GUI (`fyne.io/fyne/v2`), single `main` package
- **CGO required** for Fyne; do not assume pure-Go cross-compile from one host for all targets
- Official release archives: push tag `v*` → [.github/workflows/release.yml](.github/workflows/release.yml) (native runners). Names must stay consistent with `releaseArtifactMatchesUpdate` in [platform.go](platform.go) if you change archive patterns.

## Code style

- Match existing files: standard library first, then third-party; early returns; minimal comments (prefer clear names).
- `main.VERSION` is set at link time for releases (`-X main.VERSION=...`); default in source is `DEV` in [gui.go](gui.go).

## Testing

- Tests: `go test ./...` — table-driven helpers in [platform_test.go](platform_test.go). Run with `CGO_ENABLED=1` so the package builds like CI.

## Boundaries

- **Always:** run `go vet ./...` and `go test ./...` with `CGO_ENABLED=1` before committing GUI/runtime changes.
- **Ask first:** changing release archive naming, adding dependencies, or altering CI/release matrices (breaks update checker or users).
- **Never:** commit secrets; weaken zip path handling in `extractZipEntryByBasenames`; drop Windows `setHideWindow` behavior without intentional UX change.

## Changelog

- 2026-05-03: GitHub Actions upgraded (checkout v6, setup-go v6, artifacts v7/v8, goreleaser-action v7, action-gh-release v3).
- 2026-05-03: AGENTS troubleshooting for gopls + fyne/CGO; workspace `CGO_ENABLED` for Go tools.
- 2026-05-03: `go` directive 1.25; doc refresh.
- 2026-05-03: Initial AGENTS.md (post cross-platform / CI work).
