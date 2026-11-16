#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run -ldflags "-X main.version=dev" . man | gzip -c -9 >manpages/yt-tracks.1.gz
