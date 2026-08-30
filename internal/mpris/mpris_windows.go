//go:build windows

package mpris

import (
	"github.com/kumneger0/yt-tracks/internal/types"
	"github.com/kumneger0/yt-tracks/internal/ui"
)

func GetDbusInstance() (*ui.Instance, *chan types.DBusMessage, error) {
	// TODO: implement mpris for windows
	return nil, nil, nil
}
