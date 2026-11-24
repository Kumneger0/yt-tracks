#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run -ldflags "-X main.version=dev" main.go completion "$sh" >"completions/yt-tracks.$sh"
done
