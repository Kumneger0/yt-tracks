//go:build windows && amd64

package backend

import "embed"

//go:embed all:binaries/windows_amd64
var embeddedUVFS embed.FS
