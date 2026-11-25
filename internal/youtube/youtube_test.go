package youtube

import (
	"context"
	"slices"
	"testing"
	"time"

	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/types"
)

func TestStreamAndDuration_IsFresh(t *testing.T) {
	const testAudioURL = "http://example.com/audio"

	tests := []struct {
		name     string
		input    *types.StreamAndDuration
		maxAge   time.Duration
		expected bool
	}{
		{
			name:     "nil object",
			input:    nil,
			maxAge:   StreamURLMaxAge,
			expected: false,
		},
		{
			name: "empty URL",
			input: &types.StreamAndDuration{
				URL:       "",
				Duration:  "120",
				FetchedAt: time.Now(),
			},
			maxAge:   StreamURLMaxAge,
			expected: false,
		},
		{
			name: "zero FetchedAt",
			input: &types.StreamAndDuration{
				URL:       testAudioURL,
				Duration:  "120",
				FetchedAt: time.Time{},
			},
			maxAge:   StreamURLMaxAge,
			expected: false,
		},
		{
			name: "expired FetchedAt",
			input: &types.StreamAndDuration{
				URL:       testAudioURL,
				Duration:  "120",
				FetchedAt: time.Now().Add(-25 * time.Minute),
			},
			maxAge:   20 * time.Minute,
			expected: false,
		},
		{
			name: "fresh FetchedAt",
			input: &types.StreamAndDuration{
				URL:       testAudioURL,
				Duration:  "120",
				FetchedAt: time.Now().Add(-5 * time.Minute),
			},
			maxAge:   20 * time.Minute,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.IsFresh(tc.maxAge)
			if got != tc.expected {
				t.Errorf("IsFresh() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestSearchAndDownloadMusic_ClearsStreamAndDuration(t *testing.T) {
	track := &types.PlaylistTrackObject{
		Track: &musicpb.Song{VideoId: "test-video-id"},
		StreamAndDuration: &types.StreamAndDuration{
			URL:       "http://example.com/stream",
			Duration:  "180",
			FetchedAt: time.Now(),
		},
	}

	coreDeps := &CoreDepsPath{
		FFmpeg: "nonexistent-ffmpeg",
		YtDlp:  "nonexistent-ytdlp",
	}

	cmd := SearchAndDownloadMusic(context.Background(), track, coreDeps)
	if cmd == nil {
		t.Fatal("expected non-nil tea.Cmd")
	}

	_ = cmd()

	if track.StreamAndDuration != nil {
		t.Errorf("expected track.StreamAndDuration to be cleared to nil after playback call, got %+v", track.StreamAndDuration)
	}
}

func TestBuildYtDlpArgs(t *testing.T) {
	// Case 1: no cookie, non-verbose default
	args := buildYtDlpArgs("vid123", "", false)
	expectedTarget := "https://www.youtube.com/watch?v=vid123"
	if !slices.Contains(args, expectedTarget) {
		t.Errorf("expected args to contain target URL %s, got %v", expectedTarget, args)
	}
	if slices.Contains(args, "-v") {
		t.Errorf("expected default args to not contain -v, got %v", args)
	}
	if slices.Contains(args, "--cookies") {
		t.Errorf("expected args to not contain --cookies, got %v", args)
	}
	if slices.Contains(args, "--config-locations") {
		t.Errorf("expected args to not contain --config-locations, got %v", args)
	}

	// Case 2: with cookie and explicit verbose mode
	args = buildYtDlpArgs("vid123", "/path/to/cookie.txt", true)
	if !slices.Contains(args, "-v") {
		t.Errorf("expected verbose args to contain -v, got %v", args)
	}
	cookieIdx := slices.Index(args, "--cookies")
	if cookieIdx == -1 || cookieIdx+1 >= len(args) || args[cookieIdx+1] != "/path/to/cookie.txt" {
		t.Errorf("expected --cookies flag with cookie path, got %v", args)
	}
}

func TestBuildYtDlpFallbackArgs(t *testing.T) {
	args := buildYtDlpFallbackArgs("vid123", "/path/to/cookie.txt")
	if !slices.Contains(args, "--ignore-config") {
		t.Errorf("expected fallback args to contain --ignore-config, got %v", args)
	}
	if !slices.Contains(args, "https://www.youtube.com/watch?v=vid123") {
		t.Errorf("expected fallback args to contain target URL, got %v", args)
	}
	cookieIdx := slices.Index(args, "--cookies")
	if cookieIdx == -1 || cookieIdx+1 >= len(args) || args[cookieIdx+1] != "/path/to/cookie.txt" {
		t.Errorf("expected fallback args to contain --cookies flag, got %v", args)
	}
}
