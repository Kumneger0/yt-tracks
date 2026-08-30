package ui

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/types"
)

type CustomDelegate struct {
	list.DefaultDelegate
	*Model
}

func (d CustomDelegate) Height() int {
	return 1
}

func (d CustomDelegate) Spacing() int {
	return 0
}

func (d CustomDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	isSelected := isItemSelected(d, m, index)
	var selectedTrack *SelectedTrack
	if d.Model != nil {
		selectedTrack = d.Model.SelectedTrack
	}
	icon, title, subtitle := extractItemDisplayInfo(item, selectedTrack)
	rendered := formatRenderedItemLine(icon, title, subtitle, m.Width(), isSelected)
	fmt.Fprint(w, rendered)
}

func isItemSelected(d CustomDelegate, m list.Model, index int) bool {
	if d.Model == nil {
		return false
	}
	title := m.Title
	switch d.Model.FocusedOn {
	case SideView:
		if strings.EqualFold(title, "yt-tracks") || strings.EqualFold(title, "Library") {
			return m.Index() == index
		}
	case MainView:
		if !strings.EqualFold(title, "Related") && !strings.EqualFold(title, "Queue") && !strings.EqualFold(title, "yt-tracks") {
			return m.Index() == index
		}
	case QueueList:
		if strings.EqualFold(title, "Related") || strings.EqualFold(title, "Queue") {
			return m.Index() == index
		}
	}
	return false
}

const (
	defaultPlaylistSubtitle = "Playlist"
	contentTypeSong         = "song"
	contentTypeAlbum        = "album"
	contentTypeVideo        = "video"
)

func extractItemDisplayInfo(item list.Item, selectedTrack *SelectedTrack) (icon, title, subtitle string) {
	switch item := item.(type) {
	case types.SongItem:
		icon = SongIcon
		if item.Song != nil {
			title = item.Title
			subtitle = formatArtistNames(item.Artists)
		}
	case types.SearchResultSongItem:
		icon = SongIcon
		if item.SearchResultSong != nil {
			title = item.Title
			subtitle = formatArtistNames(item.Artists)
		}
	case types.SearchResultArtistItem:
		icon = ArtistIcon
		if item.SearchResultArtist != nil {
			title = item.Name
			if item.Subscribers != "" {
				subtitle = fmt.Sprintf("%s — %s subscribers", item.Name, item.Subscribers)
			} else {
				subtitle = "Artist"
			}
		}
	case types.SearchResultPlaylistItem:
		icon = PlaylistIcon
		if item.SearchResultPlaylist != nil {
			title = item.Title
			if item.Author != "" {
				subtitle = item.Author
			} else {
				subtitle = defaultPlaylistSubtitle
			}
		}
	case types.SearchResultAlbumItem:
		icon = AlbumIcon
		if item.SearchResultAlbum != nil {
			title = item.Title
			subtitle = fmt.Sprintf("%s • %s", item.Type, item.Year)
		}
	case types.SearchResultPodcastItem:
		icon = PodcastIcon
		if item.SearchResultPodcast != nil {
			title = item.Title
			if item.Author != "" {
				subtitle = fmt.Sprintf("%s — podcast", item.Author)
			} else {
				subtitle = "Podcast Show"
			}
		}
	case types.SearchResultEpisodeItem:
		icon = EpisodeIcon
		if item.SearchResultEpisode != nil {
			title = item.Title
			switch {
			case item.PodcastName != "" && item.Date != "":
				subtitle = fmt.Sprintf("%s • %s", item.PodcastName, item.Date)
			case item.PodcastName != "":
				subtitle = item.PodcastName
			default:
				subtitle = "Podcast Episode"
			}
		}
	case types.AlbumItem:
		icon = AlbumIcon
		if item.Album != nil {
			title = item.Title
			subtitle = fmt.Sprintf("%s • %s", item.Type, item.Year)
		}
	case types.PlaylistItem:
		icon = PlaylistIcon
		if item.Playlist != nil {
			title = item.Title
			if item.Author != "" {
				subtitle = item.Author
			} else {
				subtitle = defaultPlaylistSubtitle
			}
		}
	case types.FollowedArtistItem:
		icon = ArtistIcon
		if item.FollowedArtist != nil {
			title = item.Name
			if item.Subscribers != "" {
				subtitle = item.Subscribers
			} else {
				subtitle = "Artist"
			}
		}
	case types.LibraryChannelItem:
		icon = ArtistIcon
		if item.LibraryChannel != nil {
			title = item.Name
			subtitle = "Channel"
		}
	case types.PodcastItem:
		icon = PodcastIcon
		if item.Podcast != nil {
			title = item.Title
			subtitle = item.Author
		}
	case types.SongRelatedContentItem:
		if item.SongRelatedContent != nil {
			switch {
			case item.VideoId != "" || item.ContentType == contentTypeSong || item.ContentType == contentTypeVideo:
				icon = SongIcon
			case item.ContentType == "artist" || strings.HasPrefix(item.BrowseId, "UC") || item.Subscribers != "":
				icon = ArtistIcon
			case item.ContentType == contentTypeAlbum || strings.HasPrefix(item.BrowseId, "MPRE"):
				icon = AlbumIcon
			default:
				icon = PlaylistIcon
			}
			title = item.Title
			subtitle = item.Description
		}
	case types.PlaylistTrackObject:
		icon = SongIcon
		if item.Track != nil {
			title = item.Track.Title
			if selectedTrack != nil && selectedTrack.Track != nil &&
				item.Track.VideoId == selectedTrack.Track.VideoId {
				title += " (current)"
			}
			subtitle = formatArtistNames(item.Track.Artists)
		}
	case types.SidebarItem:
		icon = item.Icon
		title = item.Name
	case types.HomePageContentItem:
		switch {
		case item.VideoID != "" || item.ContentType == contentTypeSong || item.ContentType == contentTypeVideo:
			icon = SongIcon
		case item.ContentType == contentTypeAlbum || strings.HasPrefix(item.BrowseID, "MPRE"):
			icon = AlbumIcon
		default:
			icon = PlaylistIcon
		}
		title = item.ItemTitle
		if len(item.Artists) > 0 {
			subtitle = formatArtistNames(item.Artists)
		} else {
			subtitle = item.Description
		}
	case types.HomePageSectionItem:
		icon = SectionIcon
		title = item.SectionTitle
	case types.UserSavedTracksListItem:
		title = item.FilterValue()
		icon = LikedIcon
	}
	return icon, title, subtitle
}

func formatArtistNames(artists []*musicpb.Artist) string {
	if len(artists) == 0 {
		return ""
	}
	names := make([]string, 0, len(artists))
	for _, a := range artists {
		if a != nil && a.Name != "" {
			names = append(names, a.Name)
		}
	}
	return strings.Join(names, ", ")
}

func formatRenderedItemLine(icon, title, subtitle string, availableWidth int, isSelected bool) string {
	if availableWidth <= 0 {
		availableWidth = 40
	}

	prefix := fmt.Sprintf(" %s %s", icon, title)
	prefixWidth := lipgloss.Width(prefix)

	if subtitle != "" && availableWidth >= prefixWidth+8 {
		sep := " · "
		sepWidth := lipgloss.Width(sep)
		maxSubWidth := availableWidth - prefixWidth - sepWidth - 1
		if maxSubWidth > 3 {
			sub := truncateText(subtitle, maxSubWidth)
			if isSelected {
				return selectedStyle.Render(prefix) +
					selectedStyle.Foreground(lipgloss.Color("#D4D4D8")).Render(sep+sub+" ")
			}
			return normalStyle.Render(prefix) +
				dimStyle.Render(sep+sub+" ")
		}
	}

	if prefixWidth > availableWidth-1 {
		iconWidth := lipgloss.Width(fmt.Sprintf(" %s ", icon))
		maxTitleWidth := availableWidth - iconWidth - 1
		if maxTitleWidth > 3 {
			title = truncateText(title, maxTitleWidth)
			prefix = fmt.Sprintf(" %s %s", icon, title)
		}
	}
	str := prefix + " "
	if isSelected {
		return selectedStyle.Render(str)
	}
	return normalStyle.Render(str)
}

func truncateText(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxW {
		return s
	}
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes)+"…") > maxW {
		runes = runes[:len(runes)-1]
	}
	if len(runes) == 0 {
		return ""
	}
	return string(runes) + "…"
}

func renderSearchBar(m *Model, width int) string {
	if width < 20 {
		width = 20
	}
	m.Search.Width = width - 6

	box := lipgloss.NewStyle().
		Width(width).
		Padding(0, 1).
		Margin(0).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderNormal).
		Foreground(textPrimary)

	var content string
	if m.Search.Value() == "" && !m.Search.Focused() {
		content = dimmerStyle.Render("🔍 Search tracks, artists, playlists...")
	} else {
		content = strings.TrimRight(m.Search.View(), "\n")
	}
	return strings.TrimRight(box.Render(content), "\n")
}

func renderNowPlaying(m *Model, currentPosition, totalDuration time.Duration) string {
	selectedTrack := m.SelectedTrack
	if selectedTrack == nil || selectedTrack.Track == nil {
		return ""
	}
	var artistNames []string
	var artists = selectedTrack.Track.Artists
	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name)
	}
	artistName := strings.Join(artistNames, ", ")
	trackName := selectedTrack.Track.Title

	var likedIndicator string
	if selectedTrack.isLiked {
		likedIndicator = " " + LikedIcon
	} else {
		likedIndicator = " " + UnlikedIcon
	}

	barWidth := m.Width
	var progressFloat float64
	if totalDuration == 0 {
		progressFloat = 1.0
	} else {
		progressFloat = float64(currentPosition.Abs()) / float64(totalDuration.Abs()) * float64(barWidth)
	}
	progress := max(min(int(math.Max(progressFloat, 1)), barWidth), 0)

	filled := lipgloss.NewStyle().Foreground(progressFilled).Render(strings.Repeat("━", progress))
	empty := lipgloss.NewStyle().Foreground(progressEmpty).Render(strings.Repeat("─", max(barWidth-progress, 0)))

	playIcon := PauseIcon
	if m.IsPlaying() {
		playIcon = PlayIcon
	}

	trackInfo := lipgloss.NewStyle().Foreground(textPrimary).Bold(true).Render(
		fmt.Sprintf("%s %s", playIcon, trackName),
	)
	artistInfo := dimStyle.Render(fmt.Sprintf(" — %s", artistName))
	timeInfo := dimStyle.Render(fmt.Sprintf("  %s / %s", formatTime(currentPosition), formatTime(totalDuration)))
	likeInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Render(likedIndicator)

	return fmt.Sprintf("%s%s%s%s\n%s%s\n",
		trackInfo,
		artistInfo,
		timeInfo,
		likeInfo,
		filled, empty,
	)
}

func renderPlayerControls(m *Model) string {
	key := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	sep := dimmerStyle.Render("  │  ")
	label := lipgloss.NewStyle().Foreground(textSecondary)

	playPauseIcon := "▶"
	playPauseLabel := " play"
	if m.IsPlaying() {
		playPauseIcon = "⏸"
		playPauseLabel = " pause"
	}

	var parts []string
	hasTrack := m.isCurrentFocusTrack()

	switch m.FocusedOn {
	case SideView:
		parts = append(parts,
			key.Render("↵")+label.Render(" select")+dimmerStyle.Render("(enter)"),
			key.Render("✨")+label.Render(" new playlist")+dimmerStyle.Render("(ctrl+t)"),
			key.Render("→")+label.Render(" main view")+dimmerStyle.Render("(tab)"),
			key.Render("↓")+label.Render(" player")+dimmerStyle.Render("(shift+tab)"),
			key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(ctrl+k)"),
			key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
		)
	case MainView:
		if m.MainViewMode == LyricsMode {
			parts = append(parts,
				key.Render("↕")+label.Render(" scroll")+dimmerStyle.Render("(j/k)"),
				key.Render("→")+label.Render(" queue")+dimmerStyle.Render("(tab)"),
				key.Render("←")+label.Render(" sidebar")+dimmerStyle.Render("(shift+tab)"),
				key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(ctrl+k)"),
				key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
			)
		} else {
			parts = append(parts,
				key.Render("▶")+label.Render(" play")+dimmerStyle.Render("(enter)"),
				key.Render("+")+label.Render(" add queue")+dimmerStyle.Render("(a)"),
			)
			if hasTrack {
				parts = append(parts, key.Render("📋")+label.Render(" playlist")+dimmerStyle.Render("(ctrl+p)"))
			}
			parts = append(parts,
				key.Render("✨")+label.Render(" new playlist")+dimmerStyle.Render("(ctrl+t)"),
				key.Render("→")+label.Render(" queue")+dimmerStyle.Render("(tab)"),
				key.Render("←")+label.Render(" sidebar")+dimmerStyle.Render("(shift+tab)"),
				key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(ctrl+k)"),
				key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
			)
		}
	case QueueList:
		parts = append(parts,
			key.Render("▶")+label.Render(" play")+dimmerStyle.Render("(enter)"),
		)
		if hasTrack {
			parts = append(parts, key.Render("📋")+label.Render(" playlist")+dimmerStyle.Render("(ctrl+p)"))
		}
		parts = append(parts,
			key.Render("✕")+label.Render(" remove")+dimmerStyle.Render("(r)"),
			key.Render("↓")+label.Render(" player")+dimmerStyle.Render("(tab)"),
			key.Render("←")+label.Render(" main view")+dimmerStyle.Render("(shift+tab)"),
			key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(ctrl+k)"),
			key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
		)
	case Player:
		parts = append(parts,
			key.Render(playPauseIcon)+label.Render(playPauseLabel)+dimmerStyle.Render("(space)"),
			key.Render("⏮")+label.Render(" prev")+dimmerStyle.Render("(b)"),
			key.Render("⏭")+label.Render(" next")+dimmerStyle.Render("(n)"),
			key.Render("♥")+label.Render(" like")+dimmerStyle.Render("(l)"),
		)
		if hasTrack {
			parts = append(parts, key.Render("📋")+label.Render(" playlist")+dimmerStyle.Render("(ctrl+p)"))
		}
		parts = append(parts,
			key.Render("✨")+label.Render(" new playlist")+dimmerStyle.Render("(ctrl+t)"),
			key.Render("📝")+label.Render(" lyrics")+dimmerStyle.Render("(ctrl+l)"),
			key.Render("←")+label.Render(" sidebar")+dimmerStyle.Render("(tab)"),
			key.Render("→")+label.Render(" queue")+dimmerStyle.Render("(shift+tab)"),
			key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
		)
	case SearchBar:
		parts = append(parts,
			key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(enter)"),
			key.Render("✕")+label.Render(" cancel")+dimmerStyle.Render("(esc)"),
		)
		if hasTrack {
			parts = append(parts, key.Render("📋")+label.Render(" playlist")+dimmerStyle.Render("(ctrl+p)"))
		}
	default:
		parts = append(parts,
			key.Render(playPauseIcon)+label.Render(playPauseLabel)+dimmerStyle.Render("(space)"),
			key.Render("🔍")+label.Render(" search")+dimmerStyle.Render("(ctrl+k)"),
			key.Render("✕")+label.Render(" quit")+dimmerStyle.Render("(q)"),
		)
		if hasTrack {
			parts = append(parts, key.Render("📋")+label.Render(" playlist")+dimmerStyle.Render("(ctrl+p)"))
		}
	}

	return strings.Join(parts, sep)
}
