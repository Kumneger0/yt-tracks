package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"log"
	"log/slog"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gofrs/flock"
	"github.com/kumneger0/yt-tracks/internal/config"
	logSetup "github.com/kumneger0/yt-tracks/internal/logger"
	"github.com/kumneger0/yt-tracks/internal/queue"
	"github.com/kumneger0/yt-tracks/internal/youtube"
	ytMusicClient "github.com/kumneger0/yt-tracks/internal/yt-music-client"
	"go.dalton.dog/bubbleup"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kumneger0/yt-tracks/internal/mpris"
	"github.com/kumneger0/yt-tracks/internal/types"
	"github.com/kumneger0/yt-tracks/internal/ui"
)

var (
	Program *tea.Program
)

func newRootCmd(version string, debug bool, serverURL string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yt-tracks",
		Short: "youtube music player",
		RunE: func(cmd *cobra.Command, args []string) error {
			lockFilePath := filepath.Join(os.TempDir(), "yt-tracks.lock")

			fileLock := flock.New(lockFilePath)
			locked, err := fileLock.TryLock()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error trying to acquire lock: %v\n", err)
				os.Exit(1)
			}

			if !locked {
				showAnotherProcessIsRunning(lockFilePath)
				os.Exit(1)
			}
			defer func() {
				if err := fileLock.Unlock(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not unlock file: %v\n", err)
				}
				_ = os.Remove(lockFilePath)
			}()

			if runtime.GOOS != "windows" {
				pid := os.Getpid()
				if err := os.WriteFile(lockFilePath, []byte(strconv.Itoa(pid)), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not write PID to lock file: %v\n", err)
				}
			}
			return root(cmd, serverURL)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if memFile != "" && debug {
				f, err := os.Create(memFile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				defer f.Close()
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}
			if cpuFile != "" {
				pprof.StopCPUProfile()
			}
		},
	}

	cmd.AddCommand(newVersionCmd(version))
	cmd.AddCommand(ytTracksLog())
	cmd.AddCommand(ManCmd(cmd))
	cmd.AddCommand(newExtractCookieCmd(serverURL))
	return cmd
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func showAnotherProcessIsRunning(lockFilePath string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "Another instance of yt-tracks is already running.\n")
		return
	}
	pidBytes, readErr := os.ReadFile(lockFilePath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return
		}
		fmt.Fprintf(os.Stderr, "Error reading lock file: %v\n", readErr)
		os.Exit(1)
	}
	pid, parseErr := strconv.Atoi(string(pidBytes))

	if parseErr != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PID from lock file: %v\n", parseErr)
		os.Exit(1)
	}

	if isProcessRunning(pid) {
		fmt.Fprintf(os.Stderr, "Another instance of yt-tracks is already running (PID: %d).\n", pid)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Another instance of yt-tracks is not running (stale lock file for PID %d).\n", pid)
	fmt.Fprintf(os.Stderr, "Please try removing %s and running again if this persists.\n", lockFilePath)
	os.Exit(1)
}

type CLIFlags struct {
	DebugDir           string
	CacheDir           string
	DisableCache       bool
	CookiesFromBrowser string
	Cookies            string
	ConfigFromFile     *config.Config
}

func getCLIFlags(cmd *cobra.Command) *CLIFlags {
	debugDir, err := cmd.Flags().GetString("debug-dir")
	configFromFile := config.GetUserConfig(runtime.GOOS)

	if !cmd.Flags().Changed("debug-dir") && configFromFile.DebugDir != nil {
		debugDir = *configFromFile.DebugDir
	}

	if err != nil {
		slog.Error(err.Error())
	}

	cacheDir, err := cmd.Flags().GetString("cache-dir")
	if err != nil {
		slog.Error(err.Error())
	}
	if !cmd.Flags().Changed("cache-dir") && configFromFile.CacheDir != nil {
		cacheDir = *configFromFile.CacheDir
	}
	isCacheDisabled, err := cmd.Flags().GetBool("disable-cache")
	if err != nil {
		slog.Error(err.Error())
	}
	if !cmd.Flags().Changed("disable-cache") {
		isCacheDisabled = configFromFile.CacheDisabled
	}

	cookiesFromBrowser, err := cmd.Flags().GetString("cookies-from-browser")

	if !cmd.Flags().Changed("cookies-from-browser") {
		if configFromFile.YtDlpArgs != nil && configFromFile.YtDlpArgs.CookiesFromBrowser != nil {
			cookiesFromBrowser = *configFromFile.YtDlpArgs.CookiesFromBrowser
		}
	}

	if err != nil {
		slog.Error(err.Error())
	}

	cookiesFile, err := cmd.Flags().GetString("cookies")

	if err != nil {
		slog.Error(err.Error())
	}

	if !cmd.Flags().Changed("cookies") {
		if configFromFile.YtDlpArgs != nil && configFromFile.YtDlpArgs.Cookies != nil {
			cookiesFile = *configFromFile.YtDlpArgs.Cookies
		}
	}

	return &CLIFlags{
		DebugDir:           debugDir,
		CacheDir:           cacheDir,
		DisableCache:       isCacheDisabled,
		CookiesFromBrowser: cookiesFromBrowser,
		Cookies:            cookiesFile,
		ConfigFromFile:     configFromFile,
	}
}

func root(cmd *cobra.Command, serverURL string) error {
	flags := getCLIFlags(cmd)

	if err := os.MkdirAll(flags.DebugDir, 0755); err != nil {
		fmt.Printf("failed to create debug directory '%s': %v\n", flags.DebugDir, err)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(flags.DebugDir)
	if err != nil {
		fmt.Printf("failed to stat debug directory '%s': %v\n", flags.DebugDir, err)
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Printf("the debug path '%v' is not a directory\n", flags.DebugDir)
		os.Exit(1)
	}

	ytDlpArgs := config.YtDlpArgs{
		CookiesFromBrowser: nil,
		Cookies:            nil,
	}

	if flags.CookiesFromBrowser != "" {
		ytDlpArgs.CookiesFromBrowser = &flags.CookiesFromBrowser
	}

	if flags.Cookies != "" {
		ytDlpArgs.Cookies = &flags.Cookies
	}

	config.SetConfig(&config.Config{
		DebugDir:      &flags.DebugDir,
		CacheDisabled: flags.DisableCache,
		CacheDir:      &flags.CacheDir,
		YtDlpArgs:     &ytDlpArgs,
		SkipOnNoMatch: flags.ConfigFromFile.SkipOnNoMatch,
	})

	logger := logSetup.Init(flags.DebugDir)
	defer logger.Close()

	slog.Info("starting the application")
	debsCheekResults := doAllDepsInstalled()

	var missingDeps []DebsCheckResult
	for _, dep := range debsCheekResults {
		if dep.Installed == false {
			missingDeps = append(missingDeps, dep)
		}
	}

	coreDepsPath := &youtube.CoreDepsPath{}
	if len(missingDeps) > 0 {
		for _, dep := range missingDeps {
			fmt.Fprintf(os.Stderr, "Error: %s is missing. Please install %s using your system package manager.\n", dep.ToolName, dep.ToolName)
		}
		return fmt.Errorf("missing required system dependencies")
	}

	for _, dep := range debsCheekResults {
		if dep.ToolName == FFmpeg {
			coreDepsPath.FFmpeg = dep.Path
		}
		if dep.ToolName == YTDlp {
			coreDepsPath.YtDlp = dep.Path
		}
	}

	ins, messageChan, err := mpris.GetDbusInstance()
	if err != nil {
		slog.Error(err.Error())
	}

	client := ytMusicClient.GetYtMusicClient(serverURL)
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		slog.Error(err.Error())
		if ins != nil {
			_ = ins.Conn.Close()
		}
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	model := ui.Model{
		BreadcrumbItems: []types.Breadcrumb{{Name: "Home", Icon: "⌂"}},
		FocusedOn:       ui.SideView,
		DBusConn:        ins,
		MainViewMode:    ui.HomePageMode,
		YtMusicClient:   client,
		CoreDepsPath:    coreDepsPath,
		Width:           termWidth - 4,
		Height:          termHeight - 4,
	}

	dims := ui.CalculateLayoutDimensions(&model)
	model.LibraryWidth = dims.SidebarWidth
	model.MainViewWidth = dims.MainWidth
	model.PlayerSectionHeight = dims.InputHeight
	model.Queue = queue.NewRingQueue()

	model.SearchResult = list.New([]list.Item{}, ui.CustomDelegate{Model: &model}, dims.MainWidth, dims.ContentHeight)
	model.SearchResult.SetShowTitle(false)
	ui.RemoveListDefaults(&model.SearchResult)

	model.HomePageList = list.New([]list.Item{}, ui.CustomDelegate{Model: &model}, dims.MainWidth, dims.ContentHeight)
	model.HomePageList.SetShowTitle(false)
	ui.RemoveListDefaults(&model.HomePageList)

	model.RelatedList = list.New([]list.Item{}, ui.CustomDelegate{Model: &model}, dims.SidebarWidth, dims.ContentHeight)
	model.RelatedList.Title = "Related"
	ui.RemoveListDefaults(&model.RelatedList)

	sideBarItems := []struct{ name, icon string }{{name: "Home", icon: "⌂"}, {name: "Library", icon: "🔖"}}
	var SideBarMenuList []list.Item
	for _, item := range sideBarItems {
		SideBarMenuList = append(SideBarMenuList, types.SidebarItem{
			Name: item.name,
			Icon: item.icon,
		})
	}
	playlistItems := list.New([]list.Item{}, ui.CustomDelegate{Model: &model}, dims.MainWidth, dims.ContentHeight)
	playlistItems.SetShowTitle(false)
	ui.RemoveListDefaults(&playlistItems)

	input := textinput.New()
	input.Placeholder = "Search tracks, artists, albums..."
	input.Prompt = "> "
	input.CharLimit = 256

	model.Alert = *bubbleup.NewAlertModel(80, true, 10*time.Second)

	model.Search = input

	model.SideBarList = list.New(SideBarMenuList, ui.CustomDelegate{Model: &model}, dims.SidebarWidth, dims.ContentHeight)
	model.SideBarList.Title = "yt-tracks"
	ui.RemoveListDefaults(&model.SideBarList)

	queueList := list.New([]list.Item{}, ui.CustomDelegate{Model: &model}, dims.SidebarWidth, dims.ContentHeight)
	queueList.Title = "Queue"
	ui.RemoveListDefaults(&queueList)
	model.QueueList = queueList

	model.SelectedPlayListItems = playlistItems
	model.Queue = queue.NewRingQueue()
	model.UpdateListDimensions()

	fgModel := ui.NewForegroundModel()
	manager := ui.Manager{
		State:        ui.Background,
		WindowWidth:  termWidth,
		WindowHeight: termHeight,
		Foreground:   fgModel,
		Background:   model,
		OverlayMode:  ui.Search,
	}

	err = runProgram(manager, messageChan)

	if err != nil {
		slog.Error(err.Error())
	}

	if ins != nil {
		_ = ins.Conn.Close()
	}

	return err
}

func runProgram(manager ui.Manager, messageChan *chan types.DBusMessage) error {
	Program := tea.NewProgram(manager, tea.WithAltScreen(), tea.WithMouseCellMotion())

	go func() {
		if messageChan == nil {
			return
		}
		for v := range *messageChan {
			Program.Send(v)
		}
	}()

	go func() {
		if types.PlayedSecondsUpdateChan == nil {
			return
		}
		for data := range types.PlayedSecondsUpdateChan {
			Program.Send(data)
		}
	}()

	_, err := Program.Run()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

type CoreDependency string

const (
	FFmpeg  CoreDependency = "ffmpeg"
	FFprobe CoreDependency = "ffprobe"
	YTDlp   CoreDependency = "yt-dlp"
)

type DebsCheckResult struct {
	ToolName  CoreDependency
	Installed bool
	Path      string
}

func doAllDepsInstalled() []DebsCheckResult {
	toolNames := []CoreDependency{FFmpeg, FFprobe, YTDlp}
	results := []DebsCheckResult{}
	for _, toolName := range toolNames {
		pathFound, err := exec.LookPath(string(toolName))
		if err != nil {
			results = append(results, DebsCheckResult{
				ToolName:  toolName,
				Installed: false,
				Path:      "",
			})
			continue
		}
		results = append(results, DebsCheckResult{
			ToolName:  toolName,
			Installed: true,
			Path:      pathFound,
		})
	}
	return results
}

var (
	cpuFile string
	memFile string
)

func Execute(version string, debug bool, serverURL string) error {
	cmd := newRootCmd(version, debug, serverURL)
	if debug {
		cmd.PersistentFlags().StringVar(&cpuFile, "cpuprofile", "", "write cpu profile to `file`")
		cmd.PersistentFlags().StringVar(&memFile, "memprofile", "", "write memory profile to `file`")

		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if cpuFile != "" {
				f, err := os.Create(cpuFile)
				if err != nil {
					return fmt.Errorf("could not create CPU profile: %w", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					f.Close()
					return fmt.Errorf("could not start CPU profile: %w", err)
				}
			}
			return nil
		}
	}

	defaultDebugDir := filepath.Join(config.GetStateDir(runtime.GOOS), "logs")
	cmd.Flags().StringP("debug-dir", "d", defaultDebugDir, "a path to store app logs")
	cmd.Flags().StringP("cache-dir", "c", config.GetCacheDir(runtime.GOOS), "a path to store app cache")
	cmd.Flags().Bool("disable-cache", false, "disable cache")

	cmd.Flags().String("cookies-from-browser", "", "The name of the browser to load cookies from this option is used by yt-dlp see yt-dlp docs to see supported browsers")
	cmd.Flags().String("cookies", "", "cookies file the option you pass for this flag will be passed to yt-dlp checkout yt-dlp docs to learn more about this flag")

	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("Error executing root command: %w", err)
	}
	return nil
}
