package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadLatestRelease(OS string) (string, error) {
	// Get the latest release info from GitHub API
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the release info to get the download URL for the executable
	var release struct {
		Assets []struct {
			Name               string
			BrowserDownloadURL string `json:"browser_download_url"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	var url string
	for _, asset := range release.Assets {
		if asset.Name == "yt-dlp" {
			url = asset.BrowserDownloadURL
			break
		}
	}
	if url == "" {
		return "", fmt.Errorf("no executable found for %s", OS)
	}

	// Download the executable to a temporary file
	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.CreateTemp("", fmt.Sprintf("yt-dlp-temp-%s", OS))
	if err != nil {
		return "", err
	}
	err = os.Chmod(f.Name(), 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	log.Println(fmt.Sprintf("yt-dlp path: %s", f.Name()))
	return f.Name(), nil
}

func downloadLatestFfmpeg(OS string) (string, error) {
	// Get the latest release info from GitHub API
	resp, err := http.Get("https://ffbinaries.com/api/v1/version/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the release info to get the download URL for the executable
	type FfbinariesResponse struct {
		Version   string                       `json:"version"`
		Permalink string                       `json:"permalink"`
		Bin       map[string]map[string]string `json:"bin"`
	}
	var ffbinariesResponse FfbinariesResponse
	err = json.NewDecoder(resp.Body).Decode(&ffbinariesResponse)
	if err != nil {
		return "", err
	}

	url := ffbinariesResponse.Bin["osx-64"]["ffmpeg"]
	if url == "" {
		return "", fmt.Errorf("no executable found for %s", OS)
	}

	// Download the executable to a temporary file
	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.CreateTemp("", fmt.Sprintf("ffmpeg-temp-%s*.zip", OS))
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	path, err := unzip(f.Name(), OS)
	if err != nil {
		return "", err
	}
	err = os.Chmod(path, 0755)
	if err != nil {
		return "", err
	}
	log.Println(fmt.Sprintf("ffmpeg path: %s", f.Name()))
	return path, nil
}

func unzip(source string, OS string) (string, error) {
	read, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer read.Close()
	file := read.File[0]
	open, err := file.Open()
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", fmt.Sprintf("ffmpeg-temp-%s", OS))
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.ReadFrom(open)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
