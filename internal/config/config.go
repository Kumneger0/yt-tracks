package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	WINDOWS = "windows"
	DARWIN  = "darwin"
)

type YtDlpArgs struct {
	CookiesFromBrowser *string `json:"cookies-from-browser"`
	Cookies            *string `json:"cookies"`
}

type Config struct {
	DebugDir      *string    `json:"debug-dir"`
	CacheDisabled bool       `json:"disable-cache"`
	CacheDir      *string    `json:"cache-dir"`
	YtDlpArgs     *YtDlpArgs `json:"yt-dlp-args"`
	SkipOnNoMatch bool       `json:"skip-on-no-match"`
}

var userConfigDir = os.UserConfigDir
var userCacheDir = os.UserCacheDir
var userHomeDir = os.UserHomeDir

func GetConfigDir(goos string) string {
	configDir, err := userConfigDir()
	if err != nil {
		if goos == WINDOWS {
			return filepath.Join(os.Getenv("APPDATA"), "yt-tracks")
		}
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig != "" {
			return filepath.Join(xdgConfig, "yt-tracks")
		}
		if goos == DARWIN {
			return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "yt-tracks")
		}
		return filepath.Join(os.Getenv("HOME"), ".config", "yt-tracks")
	}
	return filepath.Join(configDir, "yt-tracks")
}

func GetStateDir(goos string) string {
	if goos == WINDOWS {
		return filepath.Join(os.Getenv("APPDATA"), "yt-tracks")
	}
	if goos == DARWIN {
		return filepath.Join(os.Getenv("HOME"), "Library", "State", "yt-tracks")
	}
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		homeDir, _ := userHomeDir()
		stateDir = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(stateDir, "yt-tracks")
}

func GetCacheDir(goos string) string {
	cacheDir, err := userCacheDir()
	if err != nil {
		homeDir, _ := userHomeDir()
		if goos == WINDOWS {
			return filepath.Join(os.Getenv("LOCALAPPDATA"), "yt-tracks")
		}
		if goos == DARWIN {
			return filepath.Join(os.Getenv("HOME"), "Library", "Caches", "yt-tracks")
		}
		return filepath.Join(homeDir, ".cache", "yt-tracks")
	}
	return filepath.Join(cacheDir, "yt-tracks")
}

func GetDefaultConfig(goos string) *Config {
	defaultDebugDir := filepath.Join(GetStateDir(goos), "logs")
	defaultCacheDir := GetCacheDir(goos)
	return &Config{
		DebugDir:      &defaultDebugDir,
		CacheDisabled: true,
		CacheDir:      &defaultCacheDir,
		YtDlpArgs:     &YtDlpArgs{},
		SkipOnNoMatch: true,
	}
}

func GetUserConfig(goos string) *Config {
	configPath := filepath.Join(GetConfigDir(goos), "config.json")
	fileStat, err := os.Stat(configPath)
	if err != nil {
		return GetDefaultConfig(goos)
	}
	if fileStat.IsDir() {
		slog.Error("User config is a directory", "path", configPath)
		return GetDefaultConfig(goos)
	}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("Failed to read user config", "err", err)
		return GetDefaultConfig(goos)
	}

	config := GetDefaultConfig(goos)
	err = json.Unmarshal(configFile, config)
	if err != nil {
		slog.Error("Failed to unmarshal user config", "err", err)
		return GetDefaultConfig(goos)
	}
	return config
}

var AppConfig Config

func SetConfig(config *Config) {
	AppConfig = *config
}

func GetConfig() *Config {
	return &AppConfig
}
