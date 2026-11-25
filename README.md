# yt-tracks

listen to youtube music in your terminal

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Docs](https://img.shields.io/badge/docs-yt--tracks.kumneger.dev-blue.svg)](https://yt-tracks.kumneger.dev)

## ✨ Preview

<div align="center">

### 🎤 Synced Lyrics
![Synced Lyrics](https://github.com/user-attachments/assets/146b3669-4396-44ca-a113-54ea674930b2)

<br />

### ➕ Playlist Management
![Create Playlist](https://github.com/user-attachments/assets/210b72ac-c755-443d-a3d7-76add0a03a06)

</div>

## 📚 Docs

check the [docs](https://yt-tracks.kumneger.dev) for the full lowdown on how to use this

## 🚀 Features

- [x] Home page & personalized recommendations
- [x] Search tracks & playlists
- [x] Queue management & playback control
- [x] Synced lyrics support
- [x] Library support & like/unlike songs
- [x] Playlist management
- [ ] More browsing features

## 🖥️ Local Backend

yt-tracks embeds its YouTube Music Connect RPC backend wheel (`yt-tracks-backend`)
and `uv` inside the Go binary via `//go:embed`. The backend runs on `127.0.0.1`
only and requires no user-installed Python — on first launch it downloads a
CPython runtime and the Python dependencies into your user cache directory
(network required once). You can also run `yt-tracks setup` to initialize it
in advance with live progress.

## 🤝 Contributing

wanna help? check [contributing.md](CONTRIBUTING.md) for the details. found a bug or got ideas? [open an issue](https://github.com/Kumneger0/yt-tracks/issues)

## 📄 License

mit license see [license](LICENSE) file for details

