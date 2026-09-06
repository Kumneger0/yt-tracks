package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kumneger0/yt-tracks/internal/types"
)

type dummyModel struct {
	receivedMsg tea.Msg
}

func (d *dummyModel) Init() tea.Cmd { return nil }
func (d *dummyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.receivedMsg = msg
	return d, nil
}
func (d *dummyModel) View() string { return "dummy" }

func TestManager_MessageRoutingInForegroundState(t *testing.T) {
	foreground := &dummyModel{}
	background := &dummyModel{}

	mgr := Manager{
		State:      Foreground,
		Foreground: foreground,
		Background: background,
	}

	addMsg := types.AddToPlaylistMsg{PlaylistID: testPlaylistID, TrackID: "T1"}
	_, _ = mgr.Update(addMsg)

	if background.receivedMsg != addMsg {
		t.Fatalf("expected Background model to receive AddToPlaylistMsg, got %v", background.receivedMsg)
	}

	respMsg := types.AddToPlaylistResponseMsg{Success: true}
	_, _ = mgr.Update(respMsg)

	if foreground.receivedMsg != respMsg {
		t.Fatalf("expected Foreground model to receive AddToPlaylistResponseMsg, got %v", foreground.receivedMsg)
	}
	if background.receivedMsg != respMsg {
		t.Fatalf("expected Background model to receive AddToPlaylistResponseMsg, got %v", background.receivedMsg)
	}
}

func TestManager_KeyMsgIsolationInForegroundState(t *testing.T) {
	foreground := &dummyModel{}
	background := &dummyModel{}

	mgr := Manager{
		State:      Foreground,
		Foreground: foreground,
		Background: background,
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	_, _ = mgr.Update(keyMsg)

	k, ok := foreground.receivedMsg.(tea.KeyMsg)
	if !ok || k.Type != keyMsg.Type {
		t.Fatalf("expected Foreground model to receive KeyMsg, got %v", foreground.receivedMsg)
	}
	if background.receivedMsg != nil {
		t.Fatalf("expected Background model NOT to receive KeyMsg while in Foreground state, got %v", background.receivedMsg)
	}
}
