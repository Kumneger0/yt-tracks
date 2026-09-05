package ui

import (
	"errors"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/queue"
	"github.com/kumneger0/yt-tracks/internal/types"
	"go.dalton.dog/bubbleup"
)

func newTestModel() Model {
	model := Model{
		FocusedOn:    SideView,
		MainViewMode: HomePageMode,
		Alert:        *bubbleup.NewAlertModel(80, false, 5*time.Second),
		Width:        120,
		Height:       40,
	}
	dims := CalculateLayoutDimensions(&model)
	delegate := CustomDelegate{Model: &model}
	model.SideBarList = list.New(nil, delegate, dims.SidebarWidth, dims.ContentHeight)
	model.SelectedPlayListItems = list.New(nil, delegate, dims.MainWidth, dims.ContentHeight)
	model.HomePageList = list.New(nil, delegate, dims.MainWidth, dims.ContentHeight)
	model.SearchResult = list.New(nil, delegate, dims.MainWidth, dims.ContentHeight)
	model.RelatedList = list.New(nil, delegate, dims.SidebarWidth, dims.ContentHeight)
	model.QueueList = list.New(nil, delegate, dims.SidebarWidth, dims.ContentHeight)
	model.Queue = queue.NewRingQueue()
	model.Search = textinput.New()
	return model
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	model := newTestModel()
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	result, _ := model.Update(msg)
	updated := result.(Model)

	if updated.Width != 116 {
		t.Errorf("Width: want 116, got %d", updated.Width)
	}
	if updated.Height != 36 {
		t.Errorf("Height: want 36, got %d", updated.Height)
	}
	if updated.LibraryWidth == 0 {
		t.Error("LibraryWidth was not set")
	}
	if updated.MainViewWidth == 0 {
		t.Error("MainViewWidth was not set")
	}
}

func TestUpdate_SearchingMsg(t *testing.T) {
	model := newTestModel()
	model.IsSearchLoading = false

	result, _ := model.Update(types.SearchingMsg{})
	updated := result.(Model)

	if !updated.IsSearchLoading {
		t.Error("IsSearchLoading should be true after SearchingMsg")
	}
}

func TestUpdate_PlayedSecondsUpdateMsg_NilTrack(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = nil

	result, _ := model.Update(types.PlayedSecondsUpdateMsg{CurrentSeconds: 42.0})
	updated := result.(Model)

	if updated.PlayedSeconds != 0 {
		t.Errorf("PlayedSeconds: want 0, got %f", updated.PlayedSeconds)
	}
}

func TestUpdate_PlayedSecondsUpdateMsg_UpdatesSeconds(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{DurationSeconds: 200},
		},
	}

	result, _ := model.Update(types.PlayedSecondsUpdateMsg{CurrentSeconds: 55.5})
	updated := result.(Model)

	if updated.PlayedSeconds != 55.5 {
		t.Errorf("PlayedSeconds: want 55.5, got %f", updated.PlayedSeconds)
	}
}

func TestUpdate_LikeUnlikeTrackMsg(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		isLiked: false,
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "vid1"},
		},
	}

	result, _ := model.Update(types.LikeUnlikeTrackResponseMsg{TrackID: "vid1", Liked: true})
	updated := result.(Model)

	if !updated.SelectedTrack.isLiked {
		t.Error("track should be liked after LikeUnlikeTrackMsg{Like: true}")
	}
}

func TestUpdate_LikeUnlikeTrackMsg_DifferentID(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		isLiked: false,
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "vid1"},
		},
	}

	result, _ := model.Update(types.LikeUnlikeTrackResponseMsg{TrackID: "other", Liked: true})
	updated := result.(Model)

	if updated.SelectedTrack.isLiked {
		t.Error("track should remain unliked when TrackID doesn't match")
	}
}

func TestUpdate_CheckUserSavedTrackMsg(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		isLiked: false,
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "vid1"},
		},
	}

	result, _ := model.Update(types.CheckUserSavedTrackResponseMsg{Saved: true})
	updated := result.(Model)

	if !updated.SelectedTrack.isLiked {
		t.Error("isLiked should be true after CheckUserSavedTrackResponseMsg{Saved: true}")
	}
}

func TestUpdate_HomePageResponseMsg_Success(t *testing.T) {
	model := newTestModel()
	model.IsSearchLoading = true

	resp := &musicpb.GetHomePageResponse{
		Sections: []*musicpb.HomePageSection{
			{Title: "Quick picks"},
			{Title: "Listen again"},
		},
	}

	result, _ := model.Update(types.HomePageResponseMsg{Response: resp})
	updated := result.(Model)

	if updated.IsSearchLoading {
		t.Error("IsSearchLoading should be false after HomePageResponseMsg")
	}
	if updated.MainViewMode != HomePageMode {
		t.Errorf("MainViewMode: want %q, got %q", HomePageMode, updated.MainViewMode)
	}
	if updated.HomePageViewMode != HomePageSectionView {
		t.Errorf("HomePageViewMode: want HomePageSectionView, got %d", updated.HomePageViewMode)
	}
	if updated.HomePageData == nil {
		t.Fatal("HomePageData should not be nil")
	}
	items := updated.HomePageList.Items()
	if len(items) != 2 {
		t.Fatalf("HomePageList items: want 2, got %d", len(items))
	}
}

func TestUpdate_HomePageResponseMsg_Error(t *testing.T) {
	model := newTestModel()
	model.IsSearchLoading = true

	result, cmd := model.Update(types.HomePageResponseMsg{Err: errors.New("network down")})
	updated := result.(Model)

	if updated.IsSearchLoading {
		t.Error("IsSearchLoading should be false on error")
	}
	if cmd == nil {
		t.Error("cmd should not be nil (alert batch)")
	}
}

func TestUpdate_LyricsMsg_NoLyrics(t *testing.T) {
	model := newTestModel()
	model.CurrentLyrics = &musicpb.GetLyricsResponse{}

	result, _ := model.Update(types.LyricsMsg{LyricsResponse: nil})
	updated := result.(Model)

	if updated.CurrentLyrics != nil {
		t.Error("CurrentLyrics should be nil when no lyrics found")
	}
}

func TestUpdate_LyricsMsg_WithLyrics(t *testing.T) {
	model := newTestModel()

	lyrics := &musicpb.GetLyricsResponse{
		Lyrics: "Hello world",
	}

	result, _ := model.Update(types.LyricsMsg{LyricsResponse: lyrics})
	updated := result.(Model)

	if updated.CurrentLyrics == nil {
		t.Fatal("CurrentLyrics should not be nil")
	}
	if updated.CurrentLyrics.Lyrics != "Hello world" {
		t.Errorf("Lyrics: want %q, got %q", "Hello world", updated.CurrentLyrics.Lyrics)
	}
}

func TestUpdate_GetLibraryMsg_Error(t *testing.T) {
	model := newTestModel()
	model.IsSearchLoading = true

	result, cmd := model.Update(types.GetLibraryMsg{Err: errors.New("auth failed")})
	updated := result.(Model)

	if updated.IsSearchLoading {
		t.Error("IsSearchLoading should be false on error")
	}
	if cmd == nil {
		t.Error("cmd should not be nil (alert)")
	}
}

func TestUpdate_KeyMsg_Quit(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Fatal("cmd should not be nil after quit key")
	}
}

func TestUpdate_KeyMsg_TabCyclesFocus(t *testing.T) {
	tests := []struct {
		name     string
		start    FocusedOn
		key      tea.KeyType
		wantNext FocusedOn
	}{
		{"SideView->tab->MainView", SideView, tea.KeyTab, MainView},
		{"MainView->tab->QueueList", MainView, tea.KeyTab, QueueList},
		{"QueueList->tab->Player", QueueList, tea.KeyTab, Player},
		{"Player->tab->SideView", Player, tea.KeyTab, SideView},
		{"MainView->shift+tab->SideView", MainView, tea.KeyShiftTab, SideView},
		{"Player->shift+tab->QueueList", Player, tea.KeyShiftTab, QueueList},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			model := newTestModel()
			model.FocusedOn = testCase.start
			result, _ := model.Update(tea.KeyMsg{Type: testCase.key})
			updated := result.(Model)
			if updated.FocusedOn != testCase.wantNext {
				t.Errorf("FocusedOn: want %q, got %q", testCase.wantNext, updated.FocusedOn)
			}
		})
	}
}

func TestUpdate_KeyMsg_CtrlK_OpensSearchBar(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
	updated := result.(Model)

	if updated.FocusedOn != SearchBar {
		t.Errorf("FocusedOn: want SearchBar, got %q", updated.FocusedOn)
	}
}

func TestUpdate_KeyMsg_CtrlQ_TogglesRightColumn(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = QueueList
	model.RightColumnMode = RightColumnQueue

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlQ})
	updated := result.(Model)

	if updated.RightColumnMode == RightColumnQueue {
		t.Error("RightColumnMode should have toggled away from RightColumnQueue")
	}
}

func TestUpdate_LyricsView_JKScrolling(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView
	model.MainViewMode = LyricsMode
	model.LyricsView.SetContent("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	updated := result.(Model)
	if updated.MainViewMode != LyricsMode {
		t.Errorf("MainViewMode: want LyricsMode, got %v", updated.MainViewMode)
	}

	result2, _ := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	updated2 := result2.(Model)
	if updated2.MainViewMode != LyricsMode {
		t.Errorf("MainViewMode: want LyricsMode, got %v", updated2.MainViewMode)
	}
}

func TestUpdate_QueueList_RemovalAndNavigation(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = QueueList
	model.RightColumnMode = RightColumnQueue

	trackObj := types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "song1", Title: "Song 1"}}
	model.Queue = queue.NewRingQueue()
	model.Queue.AddTrack(&trackObj)
	model.SyncQueueList()
	model.QueueList.Select(1)

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := result.(Model)
	if updated.Queue.Len() != 0 {
		t.Errorf("Queue items: want 0 after removal, got %d", updated.Queue.Len())
	}

	resultTab, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	updatedTab := resultTab.(Model)
	if updatedTab.FocusedOn != Player {
		t.Errorf("FocusedOn after tab: want Player, got %s", updatedTab.FocusedOn)
	}

	resultShiftTab, _ := model.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	updatedShiftTab := resultShiftTab.(Model)
	if updatedShiftTab.FocusedOn != MainView {
		t.Errorf("FocusedOn after shift+tab: want MainView, got %s", updatedShiftTab.FocusedOn)
	}
}

func TestUpdate_QueueList_MultiTrackRemovalAndPlayback(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = QueueList
	model.RightColumnMode = RightColumnQueue

	t1 := types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "song1", Title: "Song 1"}}
	t2 := types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "song2", Title: "Song 2"}}
	t3 := types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "song3", Title: "Song 3"}}

	model.Queue = queue.NewRingQueue()
	model.Queue.AddTrack(&t1)
	model.Queue.AddTrack(&t2)
	model.Queue.AddTrack(&t3)
	model.SyncQueueList()

	// Select row 2 (which is track 2, since row 0 is header "Queue", row 1 is track 1, row 2 is track 2)
	model.QueueList.Select(2)

	// Press 'r' to remove selected row 2 (Song 2)
	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := result.(Model)

	if updated.Queue.Len() != 2 {
		t.Fatalf("Queue items: want 2 after removal of middle track, got %d", updated.Queue.Len())
	}

	tracks := updated.Queue.AllTracks()
	if len(tracks) == 2 && (tracks[0].Track.VideoId != "song1" || tracks[1].Track.VideoId != "song3") {
		t.Errorf("Remaining tracks want [song1, song3], got [%s, %s]", tracks[0].Track.VideoId, tracks[1].Track.VideoId)
	}

	// Select row 1 (song1) and press Enter
	updated.QueueList.Select(1)
	resEnter, _ := updated.Update(tea.KeyMsg{Type: tea.KeyEnter})
	mPlayed := resEnter.(Model)
	if mPlayed.SelectedTrack == nil || mPlayed.SelectedTrack.Track == nil || mPlayed.SelectedTrack.Track.VideoId != "song1" {
		t.Errorf("Played track: want song1, got %v", mPlayed.SelectedTrack)
	}
}

func TestUpdate_QueueList_ContextItemSelection(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = QueueList
	model.RightColumnMode = RightColumnQueue

	t1 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "ctx1", Title: "Context 1"}}
	t2 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "ctx2", Title: "Context 2"}}
	t3 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "ctx3", Title: "Context 3"}}

	model.SetPlaybackContext([]*types.PlaylistTrackObject{t1, t2, t3}, "My Playlist", 0)
	model.SyncQueueList()

	// Select row 2 in QueueList (row 0 is header "Next from My Playlist", row 1 is ctx1, row 2 is ctx2)
	model.QueueList.Select(2)

	resEnter, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	mPlayed := resEnter.(Model)

	if mPlayed.PlaylistContextIndex != 1 {
		t.Errorf("PlaylistContextIndex: want 1 for ctx2, got %d", mPlayed.PlaylistContextIndex)
	}
	if mPlayed.SelectedTrack == nil || mPlayed.SelectedTrack.Track.VideoId != "ctx2" {
		t.Errorf("Played track: want ctx2, got %v", mPlayed.SelectedTrack)
	}
}

func TestUpdate_PlayerActions(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = Player
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "v1", Title: "Test Song"},
		},
	}

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	_, cmdLike := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if cmdLike == nil {
		t.Error("cmdLike should not be nil when pressing 'l' on selected track")
	}

	resultLyrics, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlL})
	updatedLyrics := resultLyrics.(Model)
	if updatedLyrics.MainViewMode != LyricsMode {
		t.Errorf("MainViewMode after ctrl+l: want LyricsMode, got %v", updatedLyrics.MainViewMode)
	}
	if updatedLyrics.FocusedOn != MainView {
		t.Errorf("FocusedOn after ctrl+l: want MainView, got %s", updatedLyrics.FocusedOn)
	}
}

func TestUpdate_SearchBar_FocusSubmitCancel(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView

	resFocus, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlK})
	updatedFocus := resFocus.(Model)
	if updatedFocus.FocusedOn != SearchBar {
		t.Fatalf("FocusedOn: want SearchBar, got %s", updatedFocus.FocusedOn)
	}

	updatedFocus.Search.SetValue("Kendrick")

	resSubmit, cmdSubmit := updatedFocus.Update(tea.KeyMsg{Type: tea.KeyEnter})
	submittedModel := resSubmit.(Model)
	if cmdSubmit == nil {
		t.Error("cmdSubmit should not be nil when pressing enter in SearchBar")
	}

	_, cmdQ := submittedModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmdQ != nil {
		msg := cmdQ()
		if msg == tea.Quit() {
			t.Error("q in SearchBar should not quit application")
		}
	}

	resCancel, _ := submittedModel.Update(tea.KeyMsg{Type: tea.KeyEscape})
	updatedCancel := resCancel.(Model)
	if updatedCancel.FocusedOn == SearchBar {
		t.Errorf("FocusedOn after escape: should not be SearchBar, got %s", updatedCancel.FocusedOn)
	}
}

func TestUpdate_SongRelatedContent_ArtistSelection(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = QueueList
	model.RightColumnMode = RightColumnRelated

	artistItem := types.SongRelatedContentItem{
		SongRelatedContent: &musicpb.SongRelatedContent{
			Title:       "Kendrick Lamar",
			BrowseId:    "UCq3V-gFmAdGlc8c-1M8-g",
			ContentType: "artist",
			Subscribers: "10M",
		},
	}

	model.RelatedList = list.New([]list.Item{artistItem}, CustomDelegate{Model: &model}, 20, 10)

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("cmd should not be nil when selecting related artist")
	}
}

func TestParseLRCTimestamp(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantMs int32
		wantOK bool
	}{
		{"1-digit fraction", "01:19.6", 79600, true},
		{"2-digit fraction", "00:19.67", 19670, true},
		{"3-digit fraction", "02:10.123", 130123, true},
		{"no fraction", "01:30", 90000, true},
		{"invalid format", "invalid", 0, false},
		{"invalid minutes", "xx:10.00", 0, false},
		{"invalid seconds", "01:yy.00", 0, false},
		{"negative minutes", "-01:10.00", 0, false},
		{"seconds equal to 60", "01:60.00", 0, false},
		{"seconds above 60", "01:75.00", 0, false},
		{"fraction longer than 3 digits", "01:10.1234", 0, false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotMs, gotOK := parseLRCTimestamp(testCase.input)
			if gotOK != testCase.wantOK {
				t.Errorf("parseLRCTimestamp(%q) ok = %v, want %v", testCase.input, gotOK, testCase.wantOK)
			}
			if gotMs != testCase.wantMs {
				t.Errorf("parseLRCTimestamp(%q) ms = %d, want %d", testCase.input, gotMs, testCase.wantMs)
			}
		})
	}
}

func TestParseSyncedLyrics(t *testing.T) {
	synced := `[ar: Rick Astley]
[ti: Never Gonna Give You Up]
[00:19.67] Line 1
[00:23.56] Line 2
[00:27.92] Line 3
`
	lines := parseSyncedLyrics(synced)
	if len(lines) != 3 {
		t.Fatalf("parseSyncedLyrics lines length: want 3, got %d", len(lines))
	}

	if lines[0].Text != "Line 1" || lines[0].StartTime != 19670 || lines[0].EndTime != 23560 {
		t.Errorf("Line 0: got text=%q, start=%d, end=%d", lines[0].Text, lines[0].StartTime, lines[0].EndTime)
	}
	if lines[1].Text != "Line 2" || lines[1].StartTime != 23560 || lines[1].EndTime != 27920 {
		t.Errorf("Line 1: got text=%q, start=%d, end=%d", lines[1].Text, lines[1].StartTime, lines[1].EndTime)
	}
	if lines[2].Text != "Line 3" || lines[2].StartTime != 27920 || lines[2].EndTime != 0 {
		t.Errorf("Line 2 (last): got text=%q, start=%d, end=%d (want end=0)", lines[2].Text, lines[2].StartTime, lines[2].EndTime)
	}
}

func TestUpdate_WatchPlaylistItemsMsg_Error(t *testing.T) {
	model := newTestModel()

	result, cmd := model.Update(types.WatchPlaylistItemsMsg{Err: errors.New("fetch failed")})
	updated := result.(Model)

	if cmd == nil {
		t.Error("cmd should not be nil on error (alert expected)")
	}
	if updated.Queue.Len() != 0 {
		t.Errorf("Queue items: want 0 on error, got %d", updated.Queue.Len())
	}
}

func TestUpdate_WatchPlaylistItemsMsg_NilResponse(t *testing.T) {
	model := newTestModel()

	result, cmd := model.Update(types.WatchPlaylistItemsMsg{WatchPlaylistItems: nil})
	updated := result.(Model)

	if cmd != nil {
		t.Error("cmd should be nil when WatchPlaylistItems is nil")
	}
	if updated.Queue.Len() != 0 {
		t.Errorf("Queue items: want 0 on nil response, got %d", updated.Queue.Len())
	}
}

func TestUpdate_WatchPlaylistItemsMsg_AddsTracks(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "current-video", Title: "Current Song"},
		},
	}

	watchResp := &musicpb.GetWatchPlaylistItemsResponse{
		Tracks: []*musicpb.Song{
			{VideoId: "watch1", Title: "Watch Track 1"},
			{VideoId: "watch2", Title: "Watch Track 2"},
			{VideoId: "watch3", Title: "Watch Track 3"},
		},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID:           "current-video",
		WatchPlaylistItems: watchResp,
	}

	result, _ := model.Update(msg)
	updated := result.(Model)

	if len(updated.PlaybackContext) != 3 {
		t.Fatalf("Queue items: want 3, got %d", len(updated.PlaybackContext))
	}

	tracks := updated.PlaybackContext
	expectedIDs := []string{"watch1", "watch2", "watch3"}
	for i, id := range expectedIDs {
		if tracks[i].Track.VideoId != id {
			t.Errorf("Track %d: want VideoId %q, got %q", i, id, tracks[i].Track.VideoId)
		}
	}
}

func TestUpdate_WatchPlaylistItemsMsg_SourceMismatch(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "current-video", Title: "Current Song"},
		},
	}

	watchResp := &musicpb.GetWatchPlaylistItemsResponse{
		Tracks: []*musicpb.Song{
			{VideoId: "watch1", Title: "Watch Track 1"},
		},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID:           "different-video",
		WatchPlaylistItems: watchResp,
	}

	result, cmd := model.Update(msg)
	updated := result.(Model)

	if cmd != nil {
		t.Error("cmd should be nil when SourceID doesn't match current track")
	}
	if updated.Queue.Len() != 0 {
		t.Errorf("Queue items: want 0 on source mismatch, got %d", updated.Queue.Len())
	}
}

func TestUpdate_WatchPlaylistItemsMsg_NoSelectedTrack(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = nil

	watchResp := &musicpb.GetWatchPlaylistItemsResponse{
		Tracks: []*musicpb.Song{
			{VideoId: "watch1", Title: "Watch Track 1"},
			{VideoId: "watch2", Title: "Watch Track 2"},
		},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID:           "some-video",
		WatchPlaylistItems: watchResp,
	}

	result, _ := model.Update(msg)
	updated := result.(Model)

	if len(updated.PlaybackContext) != 2 {
		t.Errorf("Queue items: want 2 when no selected track, got %d", len(updated.PlaybackContext))
	}
}

func TestUpdate_WatchPlaylistItemsMsg_EmptyTracks(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "current-video", Title: "Current Song"},
		},
	}

	watchResp := &musicpb.GetWatchPlaylistItemsResponse{
		Tracks: []*musicpb.Song{},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID:           "current-video",
		WatchPlaylistItems: watchResp,
	}

	result, _ := model.Update(msg)
	updated := result.(Model)

	if updated.Queue.Len() != 0 {
		t.Errorf("Queue items: want 0 with empty tracks, got %d", updated.Queue.Len())
	}
}

func TestUpdate_WatchPlaylistItemsMsg_AppendsToExistingQueue(t *testing.T) {
	model := newTestModel()
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "current-video", Title: "Current Song"},
		},
	}

	existing := types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "existing1", Title: "Existing Track"}}
	model.Queue.AddTrack(&existing)
	model.SyncQueueList()

	watchResp := &musicpb.GetWatchPlaylistItemsResponse{
		Tracks: []*musicpb.Song{
			{VideoId: "current-video", Title: "Current Song"},
			{VideoId: "watch1", Title: "Watch Track 1"},
			{VideoId: "watch2", Title: "Watch Track 2"},
		},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID:           "current-video",
		WatchPlaylistItems: watchResp,
	}

	result, _ := model.Update(msg)
	updated := result.(Model)
	model.SyncQueueList()

	items := updated.QueueList.Items()
	var tracks []*types.PlaylistTrackObject
	for _, item := range items {
		if track, ok := item.(types.PlaylistTrackObject); ok {
			tracks = append(tracks, &track)
		}
	}

	if len(tracks) != 3 {
		t.Fatalf("Queue items: want 3 (1 existing + 2 new), got %d", len(tracks))
	}

	if tracks[0].Track.VideoId != "existing1" {
		t.Errorf("Track 0: want existing1, got %q", tracks[0].Track.VideoId)
	}
	if tracks[1].Track.VideoId != "watch1" {
		t.Errorf("Track 1: want watch1, got %q", tracks[1].Track.VideoId)
	}
	if tracks[2].Track.VideoId != "watch2" {
		t.Errorf("Track 2: want watch2, got %q", tracks[2].Track.VideoId)
	}
}

func TestUpdate_WatchPlaylistItems_EmptyContext_FirstTrackSelected(t *testing.T) {
	model := newTestModel()
	model.PlaybackContext = nil
	model.PlaylistContextIndex = 0
	model.SelectedTrack = &SelectedTrack{
		PlaylistTrackObject: types.PlaylistTrackObject{
			Track: &musicpb.Song{VideoId: "current", Title: "Current Song"},
		},
	}

	msg := types.WatchPlaylistItemsMsg{
		SourceID: "current",
		WatchPlaylistItems: &musicpb.GetWatchPlaylistItemsResponse{
			Tracks: []*musicpb.Song{
				{VideoId: "current", Title: "Current Song"},
				{VideoId: "next1", Title: "Next 1"},
				{VideoId: "next2", Title: "Next 2"},
				{VideoId: "next3", Title: "Next 3"},
			},
		},
	}

	result, _ := model.Update(msg)
	updated := result.(Model)

	if len(updated.PlaybackContext) != 3 {
		t.Fatalf("PlaybackContext length: want 3, got %d", len(updated.PlaybackContext))
	}

	if updated.PlaylistContextIndex != 0 {
		t.Errorf("PlaylistContextIndex: want 0, got %d", updated.PlaylistContextIndex)
	}
	if updated.PlaybackContext[0].Track.VideoId != "next1" {
		t.Errorf("First next track: want next1, got %s", updated.PlaybackContext[0].Track.VideoId)
	}
}

func TestUpdate_HomePageEnter_NonVideoBeforePlayable_CorrectIndex(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView
	model.MainViewMode = HomePageMode
	model.HomePageViewMode = HomePageContentView

	dims := CalculateLayoutDimensions(&model)
	delegate := CustomDelegate{Model: &model}
	items := []list.Item{
		types.HomePageContentItem{ItemTitle: "Some Album", BrowseID: "MPRE_album1", ContentType: "album", VideoID: ""},
		types.HomePageContentItem{ItemTitle: "Song A", VideoID: "vidA", Artists: []*musicpb.Artist{{Name: "Artist A"}}},
		types.HomePageContentItem{ItemTitle: "Song B", VideoID: "vidB", Artists: []*musicpb.Artist{{Name: "Artist B"}}},
	}
	model.HomePageList = list.New(items, delegate, dims.MainWidth, dims.ContentHeight)
	model.HomePageList.Title = "Test Section"

	model.HomePageList.Select(1)

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := result.(Model)

	if len(updated.PlaybackContext) != 2 {
		t.Fatalf("PlaybackContext length: want 2, got %d", len(updated.PlaybackContext))
	}

	if updated.PlaylistContextIndex != 0 {
		t.Errorf("PlaylistContextIndex: want 0 (filtered position of Song A), got %d", updated.PlaylistContextIndex)
	}

	if updated.PlaybackContext[0].Track.VideoId != "vidA" {
		t.Errorf("PlaybackContext[0]: want vidA, got %s", updated.PlaybackContext[0].Track.VideoId)
	}
	if updated.PlaybackContext[1].Track.VideoId != "vidB" {
		t.Errorf("PlaybackContext[1]: want vidB, got %s", updated.PlaybackContext[1].Track.VideoId)
	}
}

func TestUpdate_PlayTrackFromList_DuplicateVideoId(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView
	model.MainViewMode = NormalMode

	dims := CalculateLayoutDimensions(&model)
	delegate := CustomDelegate{Model: &model}
	items := []list.Item{
		types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "songA", Title: "Song A (1st)"}},
		types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "songB", Title: "Song B"}},
		types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "songA", Title: "Song A (2nd)"}},
		types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "songC", Title: "Song C"}},
	}
	model.SelectedPlayListItems = list.New(items, delegate, dims.MainWidth, dims.ContentHeight)
	model.SelectedPlayListItems.Title = "Test Playlist"

	model.SelectedPlayListItems.Select(2)

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := result.(Model)

	if len(updated.PlaybackContext) != 4 {
		t.Fatalf("PlaybackContext length: want 4, got %d", len(updated.PlaybackContext))
	}
	if updated.PlaylistContextIndex != 2 {
		t.Errorf("PlaylistContextIndex: want 2 (2nd occurrence of songA), got %d", updated.PlaylistContextIndex)
	}

	if updated.PlaybackContext[3].Track.VideoId != "songC" {
		t.Errorf("Next track after 2nd songA: want songC at index 3, got %s", updated.PlaybackContext[3].Track.VideoId)
	}
}

func TestUpdate_HomePageEnter_DuplicateVideoId(t *testing.T) {
	model := newTestModel()
	model.FocusedOn = MainView
	model.MainViewMode = HomePageMode
	model.HomePageViewMode = HomePageContentView

	dims := CalculateLayoutDimensions(&model)
	delegate := CustomDelegate{Model: &model}
	items := []list.Item{
		types.HomePageContentItem{ItemTitle: "Song X (1st)", VideoID: "vidX"},
		types.HomePageContentItem{ItemTitle: "Song Y", VideoID: "vidY"},
		types.HomePageContentItem{ItemTitle: "Song X (2nd)", VideoID: "vidX"},
		types.HomePageContentItem{ItemTitle: "Song Z", VideoID: "vidZ"},
	}
	model.HomePageList = list.New(items, delegate, dims.MainWidth, dims.ContentHeight)
	model.HomePageList.Title = "Test Section"

	model.HomePageList.Select(2)

	result, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := result.(Model)

	if len(updated.PlaybackContext) != 4 {
		t.Fatalf("PlaybackContext length: want 4, got %d", len(updated.PlaybackContext))
	}
	if updated.PlaylistContextIndex != 2 {
		t.Errorf("PlaylistContextIndex: want 2 (2nd occurrence of vidX), got %d", updated.PlaylistContextIndex)
	}

	if updated.PlaybackContext[3].Track.VideoId != "vidZ" {
		t.Errorf("Next track after 2nd vidX: want vidZ at index 3, got %s", updated.PlaybackContext[3].Track.VideoId)
	}
}
