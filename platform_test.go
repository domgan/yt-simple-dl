package main

import "testing"

func TestYtDlpReleaseAsset(t *testing.T) {
	tests := []struct {
		goos, goarch string
		want         string
		wantErr      bool
	}{
		{"darwin", "arm64", "yt-dlp_macos", false},
		{"darwin", "amd64", "yt-dlp_macos", false},
		{"windows", "amd64", "yt-dlp.exe", false},
		{"windows", "arm64", "yt-dlp_arm64.exe", false},
		{"windows", "386", "yt-dlp_x86.exe", false},
		{"linux", "amd64", "yt-dlp_linux", false},
		{"linux", "arm64", "yt-dlp_linux_aarch64", false},
		{"linux", "arm", "yt-dlp_linux_armv7l.zip", false},
		{"freebsd", "amd64", "", true},
	}
	for _, tt := range tests {
		got, err := ytDlpReleaseAsset(tt.goos, tt.goarch)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ytDlpReleaseAsset(%q,%q) want error", tt.goos, tt.goarch)
			}
			continue
		}
		if err != nil {
			t.Errorf("ytDlpReleaseAsset(%q,%q): %v", tt.goos, tt.goarch, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ytDlpReleaseAsset(%q,%q) = %q, want %q", tt.goos, tt.goarch, got, tt.want)
		}
	}
}

func TestYtDlpAssetIsZip(t *testing.T) {
	if !ytDlpAssetIsZip("yt-dlp_linux_armv7l.zip") {
		t.Fatal("expected zip")
	}
	if ytDlpAssetIsZip("yt-dlp_linux") {
		t.Fatal("expected not zip")
	}
}

func TestFfbinariesPlatform(t *testing.T) {
	tests := []struct {
		goos, goarch, want string
	}{
		{"linux", "amd64", "linux-64"},
		{"linux", "arm64", "linux-arm64"},
		{"linux", "arm", "linux-armhf"},
		{"darwin", "arm64", "osx-64"},
		{"darwin", "amd64", "osx-64"},
		{"windows", "amd64", "windows-64"},
	}
	for _, tt := range tests {
		got, err := ffbinariesPlatform(tt.goos, tt.goarch)
		if err != nil {
			t.Fatalf("ffbinariesPlatform(%q,%q): %v", tt.goos, tt.goarch, err)
		}
		if got != tt.want {
			t.Errorf("ffbinariesPlatform(%q,%q) = %q, want %q", tt.goos, tt.goarch, got, tt.want)
		}
	}
}

func TestGoReleaserArchiveSubstring(t *testing.T) {
	if goReleaserArchiveSubstring("darwin", "arm64") != "_darwin_arm64" {
		t.Fatal(goReleaserArchiveSubstring("darwin", "arm64"))
	}
}

func TestReleaseArtifactMatchesUpdate(t *testing.T) {
	if !releaseArtifactMatchesUpdate("yt-simple-dl_v1.2.3_darwin_arm64.tar.gz", "darwin", "arm64") {
		t.Fatal("darwin arm64 archive")
	}
	if !releaseArtifactMatchesUpdate("yt-simple-dl_1.2.3_linux_amd64.tar.gz", "linux", "amd64") {
		t.Fatal("linux amd64 archive")
	}
	if !releaseArtifactMatchesUpdate("yt-simple-dl.exe", "windows", "amd64") {
		t.Fatal("legacy windows exe")
	}
	if releaseArtifactMatchesUpdate("yt-simple-dl_1.2.3_linux_amd64.tar.gz", "darwin", "arm64") {
		t.Fatal("should not match wrong OS")
	}
}
