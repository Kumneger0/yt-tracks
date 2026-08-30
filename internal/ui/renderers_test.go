package ui

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestTruncateText(t *testing.T) {
	assert.Equal(t, "", truncateText("hello", 0))
	assert.Equal(t, "hello", truncateText("hello", 10))
	assert.Equal(t, "hel…", truncateText("hello world", 4))
}

func TestExtractItemDisplayInfo(t *testing.T) {
	tests := []struct {
		name         string
		item         list.Item
		selTrack     *SelectedTrack
		wantIcon     string
		wantTitle    string
		wantSubtitle string
	}{
		{
			name: "SongItem",
			item: types.SongItem{
				Song: &musicpb.Song{
					VideoId: "123",
					Title:   "My Song",
					Artists: []*musicpb.Artist{{Name: "Artist A"}, {Name: "Artist B"}},
				},
			},
			wantIcon:     "♫",
			wantTitle:    "My Song",
			wantSubtitle: "Artist A, Artist B",
		},
		{
			name: "AlbumItem",
			item: types.AlbumItem{
				Album: &musicpb.Album{
					Title: "Greatest Hits",
					Type:  "Album",
					Year:  "2024",
				},
			},
			wantIcon:     "◉",
			wantTitle:    "Greatest Hits",
			wantSubtitle: "Album • 2024",
		},
		{
			name: "PlaylistItem",
			item: types.PlaylistItem{
				Playlist: &musicpb.Playlist{
					Title:  "Party Mix",
					Author: "DJ Kool",
				},
			},
			wantIcon:     "☰",
			wantTitle:    "Party Mix",
			wantSubtitle: "DJ Kool",
		},
		{
			name: "SidebarItem",
			item: types.SidebarItem{
				Name: "Library",
				Icon: "🔖",
			},
			wantIcon:     "🔖",
			wantTitle:    "Library",
			wantSubtitle: "",
		},
		{
			name: "HomePageSectionItem",
			item: types.HomePageSectionItem{
				SectionTitle: "Quick Picks",
			},
			wantIcon:     "▸",
			wantTitle:    "Quick Picks",
			wantSubtitle: "",
		},
		{
			name: "PlaylistTrackObject (current)",
			item: types.PlaylistTrackObject{
				Track: &musicpb.Song{VideoId: "v1", Title: "Current Track"},
			},
			selTrack: &SelectedTrack{
				PlaylistTrackObject: types.PlaylistTrackObject{
					Track: &musicpb.Song{VideoId: "v1"},
				},
			},
			wantIcon:     "♫",
			wantTitle:    "Current Track (current)",
			wantSubtitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon, title, subtitle := extractItemDisplayInfo(tt.item, tt.selTrack)
			assert.Equal(t, tt.wantIcon, icon)
			assert.Equal(t, tt.wantTitle, title)
			assert.Equal(t, tt.wantSubtitle, subtitle)
		})
	}
}

func TestIsItemSelected(t *testing.T) {
	d := CustomDelegate{
		Model: &Model{
			FocusedOn: SideView,
		},
	}

	m := list.New([]list.Item{types.SidebarItem{Name: "Home"}}, d, 20, 10)
	m.Title = "yt-tracks"

	assert.True(t, isItemSelected(d, m, 0))

	d.Model.FocusedOn = MainView
	m.Title = "Search Results"
	assert.True(t, isItemSelected(d, m, 0))

	m.Title = "Queue"
	assert.False(t, isItemSelected(d, m, 0))
}

func TestCustomDelegate_Render(t *testing.T) {
	d := CustomDelegate{
		Model: &Model{
			FocusedOn: SideView,
		},
	}

	m := list.New([]list.Item{types.SidebarItem{Name: "Home", Icon: "⌂"}}, d, 30, 10)
	m.Title = "yt-tracks"

	var buf bytes.Buffer
	d.Render(&buf, m, 0, types.SidebarItem{Name: "Home", Icon: "⌂"})
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "Home")
}
