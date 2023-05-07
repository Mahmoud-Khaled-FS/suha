package app

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Mahmoud-Khaled-FS/suha/utils"
)

// create struct to handle watch user actians
type Watcher struct {
	App      *App
	Watch    bool
	UrlList  []*url.URL
	lastCopy string
	autoNum  int
}

// Start Watching cmd
func (a *App) watch() {
	watcher := &Watcher{
		App: a,
	}

	utils.PrintWatch(a.Auto)

	initClipboard()

	go watcher.watchClipboard()
	watcher.handleUserInput()
}

// Run in his own thread
func (w *Watcher) watchClipboard() {
	w.Watch = true
	for {
		if !w.Watch {
			break
		}
		w.getClip()
		time.Sleep(100 * time.Millisecond) // Wait before checking again
	}
}

// To handle user commands when app watch user copies
func (w *Watcher) handleUserInput() {
loop:
	for {
		var s string
		fmt.Scanln(&s)
		s = strings.ToLower(s)
		utils.MoveTop(1)
		utils.ClearLine()
		switch s {
		case "q":
			w.Watch = false
			w.App.UrlList = w.UrlList
			w.App.DownloadList()
			break loop
		case "c":
			w.Watch = false
			fmt.Println("Exit without download!")
			os.Exit(0)
		case "l":
			w.removeElement()
		case "e":
			for range w.UrlList {
				w.removeElement()
			}
		}
	}
}

func (w *Watcher) removeElement() {
	if len(w.UrlList) > 0 {
		w.UrlList = w.UrlList[1:]
		utils.MoveTop(1)
		utils.ClearLine()
	}
}

func (w *Watcher) getClip() {
	cmd := exec.Command("powershell", "get-clipboard")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	textCopy := strings.TrimSpace(string(out))
	if textCopy != w.lastCopy && textCopy != "" {
		w.lastCopy = textCopy
		url, err := url.ParseRequestURI(textCopy)
		if err != nil {
			return
		}
		if w.App.Auto {
			w.autoNum += 1
			w.App.Name = ""
			w.App.Url = url
			utils.PrintCopiedUrl(url, w.autoNum, utils.InStack)
			err := w.App.DownloadFile()
			utils.ClearText()
			if err != nil {
				utils.PrintCopiedUrl(url, w.autoNum, utils.Failed)
			} else {
				utils.PrintCopiedUrl(url, w.autoNum, utils.Success)
			}
			fmt.Println("")
			return
		} else {
			w.UrlList = append(w.UrlList, url)
		}
		utils.PrintCopiedUrl(url, len(w.UrlList), utils.InStack)
		fmt.Println("")
	}
}

func initClipboard() {
	cmd := exec.Command("powershell", "echo Suha|clip")
	cmd.Output()
}

//
