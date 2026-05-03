[![CI](https://github.com/domgan/yt-simple-dl/actions/workflows/ci.yml/badge.svg)](https://github.com/domgan/yt-simple-dl/actions/workflows/ci.yml) [![Release](https://github.com/domgan/yt-simple-dl/actions/workflows/release.yml/badge.svg)](https://github.com/domgan/yt-simple-dl/actions/workflows/release.yml)

# yt-simple-dl

Desktop GUI for downloading YouTube videos as MP4 (optional MP3 via ffmpeg). Uses [yt-dlp](https://github.com/yt-dlp/yt-dlp) at runtime (downloaded automatically; no separate install).

## Supported platforms

- **Windows** amd64
- **macOS** arm64 and amd64
- **Linux** amd64 (glibc), arm64, armv7 (zip asset)

Development builds require **Go 1.25+** and **CGO** (Fyne). Release archives are built on GitHub Actions with native toolchains.

## Behavior

- **Output folder:** `~/Downloads/yt-simple-dl` (falls back to the OS user cache dir if `Downloads` cannot be created).
- **First run / each session:** yt-dlp is fetched from GitHub releases into a temp file and removed after the command finishes. For MP3 conversion the app prefers **`ffmpeg` on your PATH**; if missing it downloads a build via [ffbinaries](https://ffbinaries.com/) (Intel macOS binary on Apple Silicon may require Rosetta—installing ffmpeg yourself avoids that).
- **Updates:** If a newer [GitHub release](https://github.com/domgan/yt-simple-dl/releases) exists, the app can show a dialog with a download link.

## Development

```bash
go mod tidy
go run .
```

Linux dev deps (Debian/Ubuntu example):

```bash
sudo apt-get install build-essential pkg-config libgl1-mesa-dev libglu1-mesa-dev \
  libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev
```

## Local release-style binary

```bash
go build -trimpath -ldflags "-s -w -X main.VERSION=$(git describe --tags --always)" .
```

## CI / packaging

- **`.github/workflows/ci.yml`** — `go vet`, `go test`, `go build`, and `goreleaser check`.
- **`.github/workflows/release.yml`** — on tag `v*`, builds native binaries (Fyne + CGO) for linux/amd64, darwin/arm64, darwin/amd64, windows/amd64 and uploads archives to GitHub Releases (names match `releaseArtifactMatchesUpdate` in code, e.g. `yt-simple-dl_v1.0.0_linux_amd64.tar.gz`).
- **`.goreleaser.yaml`** — documents release naming and supports `goreleaser build --clean --snapshot --single-target` on your machine.

### Windows executable with icon (optional)

- Install: `go install github.com/tc-hib/go-winres@latest`
- Initialize: `go-winres init`, edit `winres/winres.json`
- Build: `go-winres make` then  
  `go build -ldflags "-H windowsgui -X main.VERSION=<tag>"`

## Legal

Respect YouTube’s Terms of Service and copyright when downloading content.
