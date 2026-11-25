//go:build darwin && arm64

package backend

import "embed"

//go:embed all:binaries/darwin_arm64
var embeddedUVFS embed.FS
