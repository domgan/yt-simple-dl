package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func downloadLatestRelease(goos, goarch string) (string, error) {
	assetName, err := ytDlpReleaseAsset(goos, goarch)
	if err != nil {
		return "", err
	}

	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	url := ""
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			url = asset.BrowserDownloadURL
			break
		}
	}
	if url == "" {
		return "", fmt.Errorf("no yt-dlp asset %q found for %s/%s", assetName, goos, goarch)
	}

	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpPattern := "yt-dlp-*"
	if ytDlpAssetIsZip(assetName) {
		tmpPattern = "yt-dlp-*.zip"
	}
	f, err := os.CreateTemp("", tmpPattern)
	if err != nil {
		return "", err
	}
	tmpPath := f.Name()
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	if ytDlpAssetIsZip(assetName) {
		binPath, err := extractZipEntryByBasenames(tmpPath, []string{"yt-dlp", "yt-dlp.exe", "yt-dlp_linux_armv7l"})
		os.Remove(tmpPath)
		if err != nil {
			return "", err
		}
		log.Printf("yt-dlp path: %s", binPath)
		return binPath, nil
	}

	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return "", err
	}
	log.Printf("yt-dlp path: %s", tmpPath)
	return tmpPath, nil
}

func downloadLatestFfmpeg(goos, goarch string) (string, error) {
	plat, err := ffbinariesPlatform(goos, goarch)
	if err != nil {
		return "", err
	}

	resp, err := http.Get("https://ffbinaries.com/api/v1/version/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type FfbinariesResponse struct {
		Bin map[string]map[string]string `json:"bin"`
	}
	var ffbinariesResponse FfbinariesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ffbinariesResponse); err != nil {
		return "", err
	}

	platformBin := ffbinariesResponse.Bin[plat]
	if platformBin == nil {
		return "", fmt.Errorf("ffbinaries: unknown platform %q", plat)
	}
	url := platformBin["ffmpeg"]
	if url == "" {
		return "", fmt.Errorf("no ffmpeg URL for platform %q", plat)
	}

	resp, err = http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	zipFile, err := os.CreateTemp("", fmt.Sprintf("*-ffmpeg-%s.zip", plat))
	if err != nil {
		return "", err
	}
	zipPath := zipFile.Name()
	if _, err := io.Copy(zipFile, resp.Body); err != nil {
		zipFile.Close()
		os.Remove(zipPath)
		return "", err
	}
	if err := zipFile.Close(); err != nil {
		os.Remove(zipPath)
		return "", err
	}

	want := []string{"ffmpeg", "ffmpeg.exe"}
	path, err := extractZipEntryByBasenames(zipPath, want)
	os.Remove(zipPath)
	if err != nil {
		return "", err
	}
	if err := os.Chmod(path, 0755); err != nil {
		os.Remove(path)
		return "", err
	}
	log.Printf("ffmpeg path: %s", path)
	return path, nil
}

func extractZipEntryByBasenames(zipPath string, basenames []string) (string, error) {
	read, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer read.Close()

	want := make(map[string]struct{}, len(basenames))
	for _, b := range basenames {
		want[strings.ToLower(b)] = struct{}{}
	}

	for _, file := range read.File {
		if file.FileInfo().IsDir() {
			continue
		}
		clean := filepath.Clean(file.Name)
		if strings.Contains(clean, "..") {
			continue
		}
		base := strings.ToLower(filepath.Base(clean))
		if _, ok := want[base]; !ok {
			continue
		}
		out, err := os.CreateTemp("", fmt.Sprintf("*-%s", base))
		if err != nil {
			return "", err
		}
		outPath := out.Name()
		rc, err := file.Open()
		if err != nil {
			out.Close()
			os.Remove(outPath)
			return "", err
		}
		_, copyErr := io.Copy(out, rc)
		closeErr := rc.Close()
		if err := out.Close(); err != nil && copyErr == nil {
			copyErr = err
		}
		if copyErr != nil {
			os.Remove(outPath)
			return "", copyErr
		}
		if closeErr != nil {
			os.Remove(outPath)
			return "", closeErr
		}
		return outPath, nil
	}
	return "", fmt.Errorf("no matching entry in zip %s (wanted basenames %v)", zipPath, basenames)
}

func checkVersion(currentVersion string) (string, string, error) {
	resp, err := http.Get("https://api.github.com/repos/domgan/yt-simple-dl/releases/latest")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	cur := strings.TrimPrefix(currentVersion, "v")
	if latest == cur || release.TagName == currentVersion {
		return "", "", nil
	}

	for _, asset := range release.Assets {
		if releaseArtifactMatchesUpdate(asset.Name, runtime.GOOS, runtime.GOARCH) {
			return release.TagName, asset.BrowserDownloadURL, nil
		}
	}
	return "", "", nil
}
