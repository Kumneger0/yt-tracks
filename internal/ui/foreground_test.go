package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/types"
)

const (
	testTrackID    = "track123"
	testTrackTitle = "Blinding Lights"
	testPlaylistID = "PL1"
)

func TestForegroundModel_CreatePlaylistModal(t *testing.T) {
	foreground := NewForegroundModel()

	model, _ := foreground.Update(types.OpenModalMsg{ModalType: types.ModalTypeCreatePlaylist})
	foreground = model.(*ForegroundModel)

	if foreground.ActiveModal != types.ModalTypeCreatePlaylist {
		t.Fatalf("expected ActiveModal to be ModalTypeCreatePlaylist, got %v", foreground.ActiveModal)
	}
	if foreground.FocusIndex != FieldTitle {
		t.Fatalf("expected initial FocusIndex to be FieldTitle, got %v", foreground.FocusIndex)
	}

	view := foreground.View()
	if view == "" {
		t.Fatalf("expected non-empty view when modal is active")
	}

	model, _ = foreground.Update(tea.KeyMsg{Type: tea.KeyTab})
	foreground = model.(*ForegroundModel)
	if foreground.FocusIndex != FieldDescription {
		t.Fatalf("expected FocusIndex to be FieldDescription after Tab, got %v", foreground.FocusIndex)
	}

	model, _ = foreground.Update(tea.KeyMsg{Type: tea.KeyTab})
	foreground = model.(*ForegroundModel)
	if foreground.FocusIndex != FieldPrivacy {
		t.Fatalf("expected FocusIndex to be FieldPrivacy after second Tab, got %v", foreground.FocusIndex)
	}

	model, _ = foreground.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	foreground = model.(*ForegroundModel)
	if foreground.PrivacyOptions[foreground.PrivacyIndex] != "PUBLIC" {
		t.Fatalf("expected Privacy status to be PUBLIC after Space, got %s", foreground.PrivacyOptions[foreground.PrivacyIndex])
	}

	model, cmd := foreground.Update(tea.KeyMsg{Type: tea.KeyEsc})
	foreground = model.(*ForegroundModel)
	if foreground.ActiveModal != types.ModalTypeNone {
		t.Fatalf("expected ActiveModal to be ModalTypeNone after Esc, got %v", foreground.ActiveModal)
	}
	if cmd == nil {
		t.Fatalf("expected CloseModalMsg command on Esc")
	}
}

func TestForegroundModel_AddToPlaylistModal(t *testing.T) {
	foreground := NewForegroundModel()

	model, _ := foreground.Update(types.OpenAddToPlaylistLoadingMsg{
		TrackID:    testTrackID,
		TrackTitle: testTrackTitle,
	})
	foreground = model.(*ForegroundModel)

	if foreground.ActiveModal != types.ModalTypePlaylistManagement || !foreground.IsLoading {
		t.Fatalf("expected ModalTypePlaylistManagement with IsLoading=true")
	}

	loadingView := foreground.View()
	if loadingView == "" {
		t.Fatalf("expected non-empty loading view")
	}

	pls := []*musicpb.Playlist{
		{PlaylistId: testPlaylistID, Title: "Chill Vibes", Count: 10},
		{PlaylistId: "PL2", Title: "Workout Hits", Count: 25},
	}

	model, _ = foreground.Update(types.OpenAddToPlaylistModalMsg{
		TrackID:    testTrackID,
		TrackTitle: testTrackTitle,
		Playlists:  pls,
	})
	foreground = model.(*ForegroundModel)

	if foreground.IsLoading {
		t.Fatalf("expected IsLoading=false after playlists arrive")
	}

	model, _ = foreground.Update(tea.KeyMsg{Type: tea.KeyDown})
	foreground = model.(*ForegroundModel)
	if foreground.PlaylistSelectIndex != 1 {
		t.Fatalf("expected PlaylistSelectIndex to be 1 after Down, got %d", foreground.PlaylistSelectIndex)
	}

	model, cmd := foreground.Update(tea.KeyMsg{Type: tea.KeyEnter})
	foreground = model.(*ForegroundModel)
	if !foreground.IsSubmitting {
		t.Fatalf("expected IsSubmitting to be true on Enter")
	}
	if cmd == nil {
		t.Fatalf("expected command on Enter")
	}

	model, _ = foreground.Update(types.AddToPlaylistResponseMsg{Success: true})
	foreground = model.(*ForegroundModel)
	if foreground.ActiveModal != types.ModalTypeNone {
		t.Fatalf("expected modal to close on Success, got %v", foreground.ActiveModal)
	}
}

func TestForegroundModel_DuplicateConfirmModal(t *testing.T) {
	foreground := NewForegroundModel()

	model, _ := foreground.Update(types.PromptDuplicateConfirmMsg{
		PlaylistID:   testPlaylistID,
		PlaylistName: "Chill Vibes",
		TrackID:      testTrackID,
		TrackTitle:   testTrackTitle,
	})
	foreground = model.(*ForegroundModel)

	if foreground.ActiveModal != types.ModalTypeDuplicateConfirm {
		t.Fatalf("expected ModalTypeDuplicateConfirm, got %v", foreground.ActiveModal)
	}

	view := foreground.View()
	if view == "" {
		t.Fatalf("expected non-empty view for DuplicateConfirm modal")
	}

	model, _ = foreground.Update(tea.KeyMsg{Type: tea.KeyRight})
	foreground = model.(*ForegroundModel)
	if foreground.ConfirmDuplicateIndex != 1 {
		t.Fatalf("expected ConfirmDuplicateIndex to be 1 (No), got %d", foreground.ConfirmDuplicateIndex)
	}
}
