package backend

import "embed"

//go:embed all:binaries/wheel
var embeddedWheelFS embed.FS
