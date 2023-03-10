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

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func downloadLatestRelease(OS string) (string, error) {
	// Get the latest release info from GitHub API
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the release info to get the download URL for the executable
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	var osName string
	if OS == "mac" {
		osName = "yt-dlp"
	} else if OS == "win" {
		osName = "yt-dlp.exe"
	}
	var url string
	for _, asset := range release.Assets {
		if asset.Name == osName {
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
	f, err := os.CreateTemp("", fmt.Sprintf("*-%s", osName))
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
	log.Printf("yt-dlp path: %s", f.Name())
	return f.Name(), nil
}

func downloadLatestFfmpeg(OS string) (string, error) { // todo usuwanie ffmpeg zip'a
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

	var osName string
	if OS == "mac" {
		osName = "osx-64"
	} else if OS == "win" {
		osName = "windows-64"
	}
	url := ffbinariesResponse.Bin[osName]["ffmpeg"]
	if url == "" {
		return "", fmt.Errorf("no executable found for %s", OS)
	}

	// Download the executable to a temporary file
	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.CreateTemp("", fmt.Sprintf("*-ffmpeg-%s.zip", osName))
	if err != nil {
		return "", err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
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
	log.Printf("ffmpeg path: %s", f.Name())
	f.Close()
	os.Remove(f.Name())
	return path, nil
}

func unzip(source string, OS string) (string, error) {
	log.Printf("source: %s", source)
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
	var osName string
	if OS == "mac" {
		osName = "ffmpeg"
	} else if OS == "win" {
		osName = "ffmpeg.exe"
	}
	f, err := os.CreateTemp("", fmt.Sprintf("*-%s", osName))
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

func checkVersion(currentVersion string) (string, string, error) {
	// Get the latest release from GitHub
	resp, err := http.Get("https://api.github.com/repos/domgan/yt-simple-dl/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	// Check if the latest release is newer than the current version
	if release.TagName != currentVersion {
		for _, asset := range release.Assets {
			if asset.Name == "yt-simple-dl.exe" {
				return release.TagName, asset.BrowserDownloadURL, nil
			}
		}
	}
	return "", "", nil
}
