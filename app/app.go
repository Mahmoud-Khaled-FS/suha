package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Mahmoud-Khaled-FS/suha/utils"
)

type Mode int

const (
	NORMAL Mode = iota
	WATCH
)

func (m Mode) String() string {
	switch m {
	case NORMAL:
		return "NORMAL"
	case WATCH:
		return "WATCH"
	default:
		return ""
	}
}

type App struct {
	Url       *url.URL
	Fource    bool
	Name      string
	Mode      Mode
	OutDir    string
	UrlList   []*url.URL
	Auto      bool
	Quality   string
	AudioOnly bool
}

func New(args []string) *App {
	app := &App{
		Mode: NORMAL,
	}
	app.parseArgs(args)
	app.handleDir()
	return app
}

func (a *App) parseArgs(args []string) {
	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if arg[0] == '-' {
			switch arg {
			case "-n":
				if len(args) <= i+1 {
					utils.PanicRed("ERROR: expected name after -n but got no arg")
				}
				skipNext = true
				a.Name = args[i+1]
			case "-o":
				if len(args) <= i+1 {
					utils.PanicRed("ERROR: expected Dir after -o but got no arg")
				}
				skipNext = true
				a.OutDir = args[i+1]
			case "-q":
				if len(args) <= i+1 {
					utils.PanicRed("ERROR: expected quality after -q but got no arg")
				}
				skipNext = true
				a.Quality = args[i+1]
			case "-f":
				a.Fource = true
			case "-w":
				a.Mode = WATCH
			case "-a":
				a.Auto = true
			case "--audio-only":
			case "--mp3":
				a.AudioOnly = true
			}
		} else {
			fileUrl, err := url.ParseRequestURI(arg)
			if err != nil {
				utils.PanicWithUsage(fmt.Sprintf("ERROR: expected url but got '%s'", arg))
			}
			a.Url = fileUrl
		}
	}
}

func (a *App) handleDir() {
	if a.OutDir == "" {
		dir := getwd()
		a.OutDir = dir
	}
	isAbs := filepath.IsAbs(a.OutDir)
	if !isAbs {
		dir := getwd()
		a.OutDir = filepath.Join(dir, a.OutDir)
	}
}

func (a *App) Start() {
	if a.Mode == NORMAL {
		err := a.DownloadFile()
		if err != nil {
			utils.PanicRed(err.Error())
		}
		utils.GreenMsg("\nDownload Success!")
	} else if a.Mode == WATCH {
		a.watch()
	} else {
		utils.PanicRed(fmt.Sprintf("ERROR: Unvalid mode %s", a.Mode))
	}
}

func getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		utils.PanicRed(fmt.Sprintf("ERROR: current directory missed\n%s", err))
	}
	return dir
}
