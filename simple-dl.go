package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func downloadVideo(link string, audio bool) error {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goos == "android" {
		return fmt.Errorf("android is not supported yet")
	}

	path, err := downloadLatestRelease(goos, goarch)
	if err != nil {
		return err
	}
	defer os.Remove(path)

	dir, err := ensureDownloadDir()
	if err != nil {
		return fmt.Errorf("download folder: %w", err)
	}
	outPattern := filepath.Join(dir, "%(title)s.%(ext)s")

	args := []string{link, "--no-playlist", "-f", "mp4", "-o", outPattern}
	if audio {
		ffmpegPath, cleanup, err := resolveFFmpeg(goos, goarch)
		if err != nil {
			return err
		}
		defer cleanup()

		ffmpegDir := filepath.Dir(ffmpegPath)
		args = append(args, "-x", "--audio-format", "mp3", "--ffmpeg-location", ffmpegDir)
	}
	cmd := exec.Command(path, args...)
	setHideWindow(cmd)

	if VERSION == "DEV" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Download finished")
	return nil
}

func resolveFFmpeg(goos, goarch string) (path string, cleanup func(), err error) {
	if p, err := exec.LookPath("ffmpeg"); err == nil {
		return p, func() {}, nil
	}
	tmp, err := downloadLatestFfmpeg(goos, goarch)
	if err != nil {
		return "", nil, err
	}
	return tmp, func() { os.Remove(tmp) }, nil
}
