//go:build !(linux && amd64) && !(darwin && arm64) && !(windows && amd64)

package backend

import "embed"

var embeddedUVFS embed.FS
