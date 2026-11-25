//go:build linux && amd64

package backend

import "embed"

//go:embed all:binaries/linux_amd64
var embeddedUVFS embed.FS
