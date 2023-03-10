# yt-simple-dl

You can download mp4 (and convert to mp3) of a Youtube video.

Installation: `go mod tidy`

Run in development mode: `go run -ldflags "-X main.MODE=DEV" .`

To build executable:
First install: `go install fyne.io/fyne/v2/cmd/fyne@latest`
* windows: `fyne package -os windows -icon icon.png` (OLD `go build -ldflags "-H windowsgui"`)
* mac: todo
* andorid: toimplement