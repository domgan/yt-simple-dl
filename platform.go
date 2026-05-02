package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ytDlpReleaseAsset(goos, goarch string) (string, error) {
	switch goos {
	case "darwin":
		return "yt-dlp_macos", nil
	case "windows":
		if goarch == "arm64" {
			return "yt-dlp_arm64.exe", nil
		}
		if goarch == "386" {
			return "yt-dlp_x86.exe", nil
		}
		return "yt-dlp.exe", nil
	case "linux":
		switch goarch {
		case "amd64":
			return "yt-dlp_linux", nil
		case "arm64":
			return "yt-dlp_linux_aarch64", nil
		case "arm":
			return "yt-dlp_linux_armv7l.zip", nil
		default:
			return "", fmt.Errorf("unsupported Linux GOARCH %q for yt-dlp", goarch)
		}
	default:
		return "", fmt.Errorf("unsupported GOOS %q", goos)
	}
}

func ytDlpAssetIsZip(assetName string) bool {
	return strings.HasSuffix(strings.ToLower(assetName), ".zip")
}

func ffbinariesPlatform(goos, goarch string) (string, error) {
	switch goos {
	case "windows":
		if goarch == "amd64" {
			return "windows-64", nil
		}
	case "darwin":
		if goarch == "amd64" || goarch == "arm64" {
			return "osx-64", nil // no darwin/arm64 blob; prefer ffmpeg on PATH on Apple Silicon
		}
	case "linux":
		switch goarch {
		case "amd64":
			return "linux-64", nil
		case "arm64":
			return "linux-arm64", nil
		case "arm":
			return "linux-armhf", nil
		}
	}
	return "", fmt.Errorf("no ffbinaries build for %s/%s", goos, goarch)
}

func goReleaserArchiveSubstring(goos, goarch string) string {
	return fmt.Sprintf("_%s_%s", goos, goarch)
}

func releaseArtifactMatchesUpdate(assetName string, goos, goarch string) bool {
	name := strings.ToLower(assetName)
	sub := strings.ToLower(goReleaserArchiveSubstring(goos, goarch))
	if strings.Contains(name, sub) && (strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz")) {
		return true
	}
	if goos == "windows" && goarch == "amd64" {
		return name == "yt-simple-dl.exe" || strings.HasSuffix(name, "_windows_amd64.exe")
	}
	return false
}

func ensureDownloadDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	downloads := filepath.Join(home, "Downloads", "yt-simple-dl")
	if err := os.MkdirAll(downloads, 0755); err != nil {
		cache, cerr := os.UserCacheDir()
		if cerr != nil {
			return "", err
		}
		fallback := filepath.Join(cache, "yt-simple-dl", "downloads")
		if err := os.MkdirAll(fallback, 0755); err != nil {
			return "", err
		}
		return fallback, nil
	}
	return downloads, nil
}
