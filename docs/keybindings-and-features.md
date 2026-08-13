# `ytmusic-tui` Reference: Prerequisites, Commands & Keybindings

---

## 1. Prerequisites

Before running `ytmusic-tui`, ensure the following system tools are installed:

- **`yt-dlp`**: Required for extracting and fetching YouTube Music audio stream URLs.
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
| `ytmusic-tui` | Launch the interactive terminal user interface. |
| `ytmusic-tui extract-cookies` | Automatically extract YouTube Music cookies from supported browsers (`chrome`, `firefox`, `safari`), verify via API, and store locally. |
| `ytmusic-tui version` | Print application version. |
| `ytmusic-tui log` | View debug logs. |

---

## 3. Interactive TUI Keybindings

### Navigation & Focus

| Key | Context | Action |
| --- | --- | --- |
| `Tab` | Global | Cycle focus forward (Sidebar &rarr; Main View &rarr; Queue/Player) |
| `Shift + Tab` | Global | Cycle focus backward |
| `Ctrl + k` | Global | Jump directly to Search Bar focus |
| `j` / `Down` | List / Lyrics | Navigate list item down / Scroll lyrics down |
| `k` / `Up` | List / Lyrics | Navigate list item up / Scroll lyrics up |
| `Enter` | Main View | Play selected track or open album / playlist / section |
| `Esc` | Search / Home | Unfocus search bar (returns focus to sidebar) or return to main Home sections view |

### Queue & Track Management

| Key | Context | Action |
| --- | --- | --- |
| `a` | Main View / Related | Add selected track to queue (works in Playlists, Albums, Search, Home sections, and Related list) |
| `r` | Queue / Related List | Remove highlighted track from queue or related list |
| `Ctrl + q` | Queue Column | Toggle right column view mode (Queue List vs Related Tracks) |
| `l` | Global | Toggle Like (`♥`) or Unlike (`♡`) for current track |

### Playback & Media Controls

| Key | Context | Action |
| --- | --- | --- |
| `Space` | Player Focused | Play / Pause current track |
| `n` | Player Focused | Skip to next track (Queue &rarr; Context &rarr; History) |
| `b` | Player Focused | Skip to previous track from history |
| `Ctrl + l` | Global | Toggle Synced Lyrics View overlay |
| `q` / `Ctrl + c` | Non-Input | Quit `ytmusic-tui` cleanly |

---

## 4. Configuration & State Files

Credentials and configuration are stored locally:
- **Linux**: `~/.config/ytmusic-tui/browser.json`
- **macOS**: `~/Library/Application Support/ytmusic-tui/browser.json`
- **Windows**: `%APPDATA%\ytmusic-tui\browser.json`
