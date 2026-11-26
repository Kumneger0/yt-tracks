# `yt-tracks` Reference: Prerequisites, Commands & Keybindings

---

## 1. Prerequisites

Before running `yt-tracks`, ensure the following system tools are installed:

- **`yt-dlp`**: Required for extracting and fetching YouTube Music audio stream URLs. Standard `yt-dlp` configuration files and plugins are loaded natively; if the initial call fails, `yt-tracks` retries with `--ignore-config` (which drops config-file flags such as `--proxy` — use `HTTP_PROXY`/`HTTPS_PROXY` environment variables if proxy routing must also apply to fallback).
- **`ffmpeg`**: Required for decoding and playing audio streams.

Install both using your OS package manager:

- **Ubuntu / Debian**: `sudo apt install ffmpeg yt-dlp`
- **Arch Linux**: `sudo pacman -S ffmpeg yt-dlp`
- **Fedora**: `sudo dnf install ffmpeg yt-dlp`
- **macOS**: `brew install ffmpeg yt-dlp`
- **Windows**: `winget install ffmpeg yt-dlp` or `scoop install ffmpeg yt-dlp`

---

## 2. CLI Commands & Subcommands

| Command | Purpose |
| --- | --- |
| `yt-tracks` | Launch the interactive terminal user interface. |
| `yt-tracks setup` | Download CPython and initialize the local backend runtime in advance (also available via `--setup`). |
| `yt-tracks extract-cookies` | Automatically extract YouTube Music cookies from supported browsers (`chrome`, `firefox`, `safari`), verify via API, and store locally. |
| `yt-tracks version` | Print application version. |
| `yt-tracks log` | View debug logs. |

### Active Flags & Environment Variables

| Option | Purpose |
| --- | --- |
| `--setup` | Download and initialize the local backend runtime without starting the TUI. |
| `-d`, `--debug-dir <path>` | Override the directory where `yt-tracks.log`, `yt-dlp-error.log`, and `ffstderr.log` are written. |
| `YT_TRACKS_YTDLP_VERBOSE=1` | Enable `yt-dlp -v` verbose output in `yt-dlp-error.log` for troubleshooting stream extraction. |

---

## 3. Interactive TUI Keybindings

### Navigation & Focus

| Key | Context | Action |
| --- | --- | --- |
| `Tab` | Global | Cycle focus forward (Sidebar &rarr; Main View &rarr; Queue &rarr; Player) |
| `Shift + Tab` | Global | Cycle focus backward |
| `Ctrl + k` | Global | Jump directly to Search Bar focus |
| `j` / `Down` | List / Lyrics | Navigate list item down / Scroll lyrics down |
| `k` / `Up` | List / Lyrics | Navigate list item up / Scroll lyrics up |
| `Enter` | Main View | Play selected track or open album / playlist / section |
| `Esc` | Search / Home / Modals | Unfocus search bar (returns focus to sidebar), return to main Home sections view, or close active modal |

### Queue, Playlist & Track Management

| Key | Context | Action |
| --- | --- | --- |
| `a` | Main View / Related | Add selected track to queue (works in Playlists, Albums, Search, Home sections, and Related list) |
| `r` | Queue List | Remove highlighted track from the queue |
| `Ctrl + q` | Queue Column | Toggle right column view mode (Queue List vs Related Tracks) |
| `Ctrl + t` | Global | Open Create New Playlist modal (`Tab`/`↑`/`↓` switch fields, `Space`/`←`/`→` select privacy, `Enter` submit) |
| `Ctrl + p` | Track Focused / Selected | Open Manage Playlist modal (`j`/`k` or `↑`/`↓` select playlist, `Enter`/`a` add, `r` remove) |
| `l` | Global | Toggle Like (`♥`) or Unlike (`♡`) for current track |

### Playback & Media Controls

| Key | Context | Action |
| --- | --- | --- |
| `Space` | Player Focused | Play / Pause current track |
| `n` | Player Focused | Skip to next track (Queue &rarr; Context &rarr; History) |
| `b` | Player Focused | Skip to previous track from history |
| `Ctrl + l` | Global | Toggle Synced Lyrics View overlay (supports mouse wheel & `j`/`k` scrolling) |
| `q` / `Ctrl + c` | Non-Input | Quit `yt-tracks` cleanly |

---

## 4. Configuration, Cache & State Files

Credentials extracted via `yt-tracks extract-cookies` (`browser.json`, plus Netscape `cookie.txt` generated for `yt-dlp`) and optional `config.json` (`{ "debug-dir": "/path/to/logs" }`) are stored in your config directory:
- **Linux**: `~/.config/yt-tracks/` (`$XDG_CONFIG_HOME/yt-tracks/`)
- **macOS**: `~/Library/Application Support/yt-tracks/`
- **Windows**: `%APPDATA%\yt-tracks\`

Logs (`yt-tracks.log`, `yt-dlp-error.log`, `ffstderr.log`) are stored in your state directory (or inspect via `yt-tracks log`):
- **Linux**: `~/.local/state/yt-tracks/logs/` (`$XDG_STATE_HOME/yt-tracks/logs/`)
- **macOS**: `~/Library/State/yt-tracks/logs/`
- **Windows**: `%APPDATA%\yt-tracks\logs\`

The embedded `uv` binary and backend CPython virtualenv are cached in your user cache directory:
- **Linux**: `~/.cache/yt-tracks/` (`$XDG_CACHE_HOME/yt-tracks/`)
- **macOS**: `~/Library/Caches/yt-tracks/`
- **Windows**: `%LOCALAPPDATA%\yt-tracks\`
