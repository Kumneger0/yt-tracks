package ui

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/gen/genconnect"
	"github.com/kumneger0/yt-tracks/internal/queue"
	"github.com/kumneger0/yt-tracks/internal/types"
	"github.com/kumneger0/yt-tracks/internal/youtube"
	"go.dalton.dog/bubbleup"
)

type FocusedOn string

const (
	SideView  FocusedOn = "SIDE_VIEW"
	MainView  FocusedOn = "MAIN_VIEW"
	Player    FocusedOn = "PLAYER"
	SearchBar FocusedOn = "SEARCH_BAR"
	QueueList FocusedOn = "QUEUE_LIST"
)

type MainViewMode string

const (
	SearchResultMode MainViewMode = "SEARCH_RESULT_MODE"
	NormalMode       MainViewMode = "NORMAL_MODE"
	LyricsMode       MainViewMode = "LYRICS_MODE"
	HomePageMode     MainViewMode = "HOME_PAGE_MODE"
)

type HomePageViewMode int

const (
	HomePageSectionView HomePageViewMode = iota
	HomePageContentView
)

type RightColumnMode string

const (
	RightColumnQueue   RightColumnMode = "RIGHT_COLUMN_QUEUE"
	RightColumnRelated RightColumnMode = "RIGHT_COLUMN_RELATED"
)

type SpotifySearchResult struct {
	Tracks, Artists, Albums, Playlists list.Model
}

type SelectedTrack struct {
	types.PlaylistTrackObject
	isLiked bool
}

type MusicQueueList struct {
	list.Model
	PaginationInfo *types.PaginationInfo
}

type Model struct {
	BreadcrumbItems       []types.Breadcrumb
	SideBarList           list.Model
	Alert                 bubbleup.AlertModel
	SelectedPlayListItems list.Model
	LyricsView            viewport.Model
	FocusedOn             FocusedOn
	MainViewMode
	PlayerProcess        *types.Player
	playbackCancel       context.CancelFunc
	SelectedTrack        *SelectedTrack
	PlayedSeconds        float64
	Height               int
	Width                int
	LibraryWidth         int
	MainViewWidth        int
	PlayerSectionHeight  int
	Search               textinput.Model
	Queue                *queue.RingQueue
	QueueList            list.Model
	PlaybackContext      []*types.PlaylistTrackObject
	PlaybackContextName  string
	PendingContextName   string
	PlaylistContextIndex int
	PlayHistory          []*types.PlaylistTrackObject
	PlayHistoryIndex     int
	YtMusicClient        genconnect.MusicServiceClient
	DBusConn             *Instance
	SearchQuery          string
	IsSearchLoading      bool
	SearchResult         list.Model
	PaginationInfo       *types.PaginationInfo
	IsOnPagination       bool
	CoreDepsPath         *youtube.CoreDepsPath
	HomePageData         *musicpb.GetHomePageResponse
	HomePageList         list.Model
	HomePageViewMode     HomePageViewMode
	RelatedList          list.Model
	RightColumnMode      RightColumnMode
	CurrentLyrics        *musicpb.GetLyricsResponse
}

type Instance struct {
	Props *prop.Properties
	Conn  *dbus.Conn
}

func (m Model) Init() tea.Cmd {
	homePageFeed := func() tea.Msg {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		homePage, err := m.YtMusicClient.GetHomePage(ctx, &musicpb.GetHomePageRequest{})
		if err != nil {
			slog.Error(err.Error())
			return types.HomePageResponseMsg{
				Response: nil,
				Err:      err,
			}
		}
		return types.HomePageResponseMsg{
			Response: homePage,
			Err:      nil,
		}
	}
	return tea.Batch(m.Alert.Init(), SendLoadingCmd(), homePageFeed, tea.SetWindowTitle("yt-tracks"))
}

func renderBreadcrumbs(items []types.Breadcrumb) string {
	if len(items) == 0 {
		return ""
	}

	parts := make([]string, 0, len(items)*2)
	for idx, item := range items {
		label := item.Name
		if item.Icon != "" {
			label = fmt.Sprintf("%s %s", item.Icon, item.Name)
		}
		itemStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
		if idx == len(items)-1 {
			itemStyle = itemStyle.Foreground(textPrimary).Bold(true)
		} else {
			itemStyle = itemStyle.Foreground(accentColor)
		}

		parts = append(parts, itemStyle.Render(label))
		if idx < len(items)-1 {
			parts = append(parts, lipgloss.NewStyle().Foreground(textDim).Render("▸"))
		}
	}

	return lipgloss.NewStyle().Padding(0, 0, 0, 1).Render(strings.Join(parts, " "))
}

func (m Model) View() string {
	dimensions := calculateLayoutDimensions(&m)
	sideBarView := getStyle(&m, dimensions.ContentHeight, dimensions.SidebarWidth, SideView).Render(m.SideBarList.View())
	searchBar := renderSearchBar(&m, dimensions.MainWidth)
	breadcrumb := renderBreadcrumbs(m.BreadcrumbItems)
	var mainView string
	switch {
	case m.IsSearchLoading:
		loadingText := dimmerStyle.Render("  ⟳ Loading...")
		mainView = getStyle(&m, dimensions.ContentHeight, dimensions.MainWidth, MainView).Render(
			lipgloss.JoinVertical(lipgloss.Top, searchBar, breadcrumb, loadingText),
		)
	case m.MainViewMode == SearchResultMode:
		resultHeader := titleStyle.Render("  Search Results")
		mainView = getStyle(&m, dimensions.ContentHeight, dimensions.MainWidth, MainView).Render(
			lipgloss.JoinVertical(lipgloss.Top, searchBar, resultHeader, lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(m.SearchResult.View())),
		)
	case m.MainViewMode == LyricsMode:
		trackName := ""
		if m.SelectedTrack != nil && m.SelectedTrack.Track != nil {
			trackName = " • " + m.SelectedTrack.Track.Title
		}
		lyricsHeader := titleStyle.Render("  📝 Lyrics" + trackName)
		lyricsPadded := lipgloss.NewStyle().Padding(1, 2).Render(m.LyricsView.View())
		mainView = getStyle(&m, dimensions.ContentHeight, dimensions.MainWidth, MainView).Render(
			lipgloss.JoinVertical(lipgloss.Top, searchBar, breadcrumb, lyricsHeader, lyricsPadded),
		)
	case m.MainViewMode == HomePageMode:
		mainView = getStyle(&m, dimensions.ContentHeight, dimensions.MainWidth, MainView).Render(
			lipgloss.JoinVertical(lipgloss.Top, searchBar, breadcrumb, lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(m.HomePageList.View())),
		)
	default:
		mainView = getStyle(&m, dimensions.ContentHeight, dimensions.MainWidth, MainView).
			Render(lipgloss.JoinVertical(lipgloss.Top, searchBar, breadcrumb, lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(m.SelectedPlayListItems.View())))
	}

	var playingView string

	if m.SelectedTrack != nil && m.SelectedTrack.Track != nil {
		playedSeconds := int(m.PlayedSeconds)
		currentPosition := time.Second * time.Duration(playedSeconds)
		total := time.Duration(m.SelectedTrack.Track.DurationSeconds) * time.Second
		playingView = renderNowPlaying(&m, currentPosition, total)
	}

	controls := renderPlayerControls(&m)
	playingCombined := strings.TrimSpace(playingView) + "\n" + controls

	playing := getPlayerStyles(&m, dimensions).
		Foreground(playerFg).
		Render(playingCombined)
	var rightColumnView string
	showRelated := m.RightColumnMode == RightColumnRelated
	if showRelated {
		rightColumnView = m.RelatedList.View()
	} else {
		rightColumnView = m.QueueList.View()
	}
	queueList := getStyle(&m, dimensions.ContentHeight, dimensions.SidebarWidth, QueueList).Render(rightColumnView)
	combinedView := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, sideBarView, mainView, queueList),
		playing,
	)
	return m.Alert.Render(combinedView)
}

func formatTime(d time.Duration) string {
	totalSeconds := int(time.Duration(math.Max(float64(d), 0)).Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

type LayoutDimensions struct {
	SidebarWidth  int
	MainWidth     int
	ContentHeight int
	InputHeight   int
}

func CalculateLayoutDimensions(model *Model) LayoutDimensions {
	if model.Width <= 0 || model.Height <= 0 {
		model.Width = 100
		model.Height = 30
	}
	sidebarWidth := model.Width * 22 / 100
	inputHeight := min(max(model.Height*10/100, 2), 3)
	mainCenterArea := max(model.Width-(sidebarWidth*2)-2, 10)

	return LayoutDimensions{
		SidebarWidth:  sidebarWidth,
		MainWidth:     mainCenterArea,
		ContentHeight: model.Height * 90 / 100,
		InputHeight:   inputHeight,
	}
}

func calculateLayoutDimensions(m *Model) LayoutDimensions {
	return CalculateLayoutDimensions(m)
}

func (m *Model) UpdateListDimensions() {
	dimensions := CalculateLayoutDimensions(m)
	m.SideBarList.SetSize(dimensions.SidebarWidth, dimensions.ContentHeight)
	m.SelectedPlayListItems.SetSize(dimensions.MainWidth, dimensions.ContentHeight-4)
	m.SearchResult.SetSize(dimensions.MainWidth, dimensions.ContentHeight-4)
	m.HomePageList.SetSize(dimensions.MainWidth, dimensions.ContentHeight-4)
	m.QueueList.SetSize(dimensions.SidebarWidth, dimensions.ContentHeight)
	if len(m.RelatedList.Items()) > 0 {
		m.RelatedList.SetSize(dimensions.SidebarWidth, dimensions.ContentHeight)
	}
}

func (m *Model) SetPlaybackContext(tracks []*types.PlaylistTrackObject, name string, currentIndex int) {
	m.PlaybackContext = tracks
	m.PlaybackContextName = name
	m.PlaylistContextIndex = currentIndex
}

func (m *Model) SyncQueueList() tea.Cmd {
	var items []list.Item
	if m.Queue != nil && m.Queue.Len() > 0 {
		userQueueTracks := m.Queue.AllTracks()
		items = append(items, types.HomePageSectionItem{SectionTitle: "Queue"})
		for _, t := range userQueueTracks {
			if t != nil {
				items = append(items, *t)
			}
		}
	}

	if len(m.PlaybackContext) > 0 {
		contextName := m.PlaybackContextName
		if contextName != "" {
			items = append(items, types.HomePageSectionItem{SectionTitle: "Next from " + contextName})
		}
		if contextName == "" {
			items = append(items, types.HomePageSectionItem{SectionTitle: "Next Up"})
		}

		limit := len(m.PlaybackContext)
		startIdx := m.PlaylistContextIndex
		if startIdx >= 0 && startIdx < limit && m.PlaybackContext[startIdx] != nil &&
			m.SelectedTrack != nil && m.SelectedTrack.Track.VideoId == m.PlaybackContext[startIdx].Track.VideoId {
			startIdx++
		}
		if startIdx < 0 || startIdx >= limit {
			startIdx = 0
		}
		for idx := startIdx; idx < limit; idx++ {
			if m.PlaybackContext[idx] != nil {
				items = append(items, *m.PlaybackContext[idx])
			}
		}
	}

	if len(items) == 0 && len(m.PlayHistory) > 0 {
		items = append(items, types.HomePageSectionItem{SectionTitle: "Recently Played"})
		for _, track := range slices.Backward(m.PlayHistory) {
			if track != nil {
				items = append(items, *track)
			}
		}
	}
	return m.QueueList.SetItems(items)
}

func RemoveListDefaults(listToRemoveDefaults *list.Model) {
	if listToRemoveDefaults != nil {
		listToRemoveDefaults.SetShowFilter(false)
		listToRemoveDefaults.SetShowPagination(false)
		listToRemoveDefaults.SetShowHelp(false)
		listToRemoveDefaults.SetShowStatusBar(false)
	}
}

func removeListDefaults(listToRemoveDefaults *list.Model) {
	RemoveListDefaults(listToRemoveDefaults)
}

func (m *Model) IsPlaying() bool {
	return m != nil && m.PlayerProcess != nil && m.PlayerProcess.OtoPlayer != nil && m.PlayerProcess.OtoPlayer.IsPlaying()
}
