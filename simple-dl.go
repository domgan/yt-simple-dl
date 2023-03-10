package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

var MODE string

func pwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func downloadVideo(link string, audio bool) error {
	var OS string
	switch runtime.GOOS {
	case "darwin":
		OS = "mac"
	case "windows":
		OS = "win"
	case "android":
		return fmt.Errorf("android is not supported yet")
	default:
		return fmt.Errorf("unknown OS")
	}

	path, err := downloadLatestRelease(OS)
	if err != nil {
		return err
	}
	defer os.Remove(path)

	args := []string{link, "-f", "mp4", "-o", pwd() + "/%(title)s.%(ext)s"}
	if audio {
		ffmpegPath, err := downloadLatestFfmpeg(OS)
		if err != nil {
			return err
		}
		defer os.Remove(ffmpegPath)

		args = append(args, "-x", "--audio-format", "mp3", "--ffmpeg-location", ffmpegPath)
	}
	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if MODE == "DEV" {
		// Create pipes to capture output from stdout and stderr
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Download finished")
	return nil
}
