//go:build darwin

package mpris

import (
	"github.com/kumneger0/yt-tracks/internal/types"
	"github.com/kumneger0/yt-tracks/internal/ui"
)

func GetDbusInstance() (*ui.Instance, *chan types.DBusMessage, error) {
	// TODO: implement mpris for darwin
	return nil, nil, nil
}
