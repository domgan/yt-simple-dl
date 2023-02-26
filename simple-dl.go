package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func pwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func downloadVideo(link string, audio bool) error {
	path, err := downloadLatestRelease("macos")
	if err != nil {
		return err
	}
	defer os.Remove(path)

	args := []string{link, "-o", pwd() + "/%(title)s"}
	if audio {
		ffmpegPath, err := downloadLatestFfmpeg("macos")
		if err != nil {
			return err
		}
		defer os.Remove(ffmpegPath)

		args = append(args, "-x", "--audio-format", "mp3", "--ffmpeg-location", ffmpegPath)
	}
	cmd := exec.Command(path, args...)

	// Create pipes to capture output from stdout and stderr
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	// Print the captured output
	log.Println("Command stdout:")
	log.Println(stdout.String())
	log.Println("Command stderr:")
	log.Println(stderr.String())

	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Download finished")
	return nil
}
