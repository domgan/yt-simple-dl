# yt-simple-dl

You can download mp4 (and convert to mp3) of a Youtube video.

## Installation 
```
go mod tidy
```

## Run in development mode
```
go run .
```

## To build executable on windows
* To embed icon and version tag, first install: `go install github.com/tc-hib/go-winres@latest` \
* Initialize with: `go-winres init`
* Make changes to `winres/winres.json`
* Run `go-winres make` to create `.syso` which will be picked up automatically by the build command
``` bash
go build -ldflags "-H windowsgui -X main.VERSION=<tag>"
```