package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
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

	args := []string{link, "-o", pwd() + "/%(title)s"}
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

	// todo add flag to not do below when it's passed for windows build
	// // Create pipes to capture output from stdout and stderr
	// var stdout bytes.Buffer
	// var stderr bytes.Buffer
	// cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	// cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	// // Print the captured output
	// log.Println("Command stdout:")
	// log.Println(stdout.String())
	// log.Println("Command stderr:")
	// log.Println(stderr.String())

	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Download finished")
	return nil
}
