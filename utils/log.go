package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Mahmoud-Khaled-FS/suha/constant"
)

func PanicWithUsage(msg string) {
	fmt.Println("usage suha <url>")
	PanicRed(msg)
}

func PanicRed(msg string) {
	fmt.Print(constant.ColorRed)
	fmt.Println(msg)
	fmt.Print(constant.ColorReset)
	os.Exit(1)
}

func GreenMsg(msg string) {
	fmt.Print(constant.ColorGreen)
	fmt.Println(msg)
	fmt.Print(constant.ColorReset)
}

func PrintWelcome() {
	fmt.Println("")
	fmt.Println("")
	fmt.Print(constant.ColorPurple)
	fmt.Println("  .d8888. db    db db   db  .d8b. ")
	fmt.Println("  88'  YP 88    88 88   88 d8' `8b")
	fmt.Println("  `8bo.   88    88 88ooo88 88ooo88")
	fmt.Println("    `Y8b. 88    88 88~~~88 88~~~88")
	fmt.Println("  db   8D 88b  d88 88   88 88   88")
	fmt.Println("  `8888Y' ~Y8888P' YP   YP YP   YP")
	fmt.Print(constant.ColorReset)
}

func PrintWatch(auto bool) {
	PrintWelcome()
	fmt.Println("")
	fmt.Print(constant.BoldText + constant.ColorPurple)
	fmt.Println("Welcome to Suha v0.1, the ultimate CLI application for file downloading!\033[0m")
	fmt.Print(constant.ColorReset)
	fmt.Println("")
	if !auto {
		fmt.Println("To initiate a download, simply enter the appropriate download command:")
		fmt.Println("- Enter '\033[1m\033[36md\033[0m' -> download files.")
		fmt.Println("- Enter '\033[1m\033[36mq\033[0m' -> download files and exit.")
		fmt.Println("- Enter '\033[1m\033[36me\033[0m' -> remove all files in the list.")
		fmt.Println("- Enter '\033[1m\033[36ml\033[0m' -> remove last file in list.")
		fmt.Println("- Enter '\033[1m\033[36mc\033[0m' -> exit.")
		fmt.Println("")
	}
	fmt.Print(constant.BoldText + constant.ColorBlue)
	fmt.Println("Happy downloading with Suha!")
	fmt.Println("")
	fmt.Print(constant.ColorReset)
}

type URLState int

const (
	InStack URLState = iota
	Failed
	Success
)

func (u *URLState) String() string {
	if *u == Failed {
		return "[Failed]"
	} else if *u == Success {
		return "[Success]"
	}
	return ""
}

func (u *URLState) Color() string {
	if *u == Failed {
		return constant.ColorRed
	} else if *u == Success {
		return constant.ColorGreen
	}
	return constant.ColorWhite
}

func PrintCopiedUrl(url *url.URL, number int, state URLState) {
	strUrl := url.String()
	if len(strUrl) > 70 {
		strUrl = strUrl[0:70] + "..."
	}
	fmt.Print(state.Color())
	fmt.Printf("[%d] -> %s   %s", number, strUrl, state.String())
	fmt.Print(constant.ColorReset)
}

func MoveTop(lineNumber int) {
	fmt.Printf("\033[%dA", lineNumber)
}
func MoveBottom(lineNumber int) {
	fmt.Printf("\033[%dB", lineNumber)
}
func MoveLeft(lineNumber int) {
	fmt.Printf("\033[%dD", lineNumber)
}
func ClearLine() {
	fmt.Print("\033[2K")
}
func ClearText() {
	MoveLeft(1000)
	fmt.Print(strings.Repeat(" ", 100))
	MoveLeft(1000)
}

type Progress struct {
	TotalSize int64
	BytesRead int64
	Unit      string
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n = len(p)
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {

	ClearText()
	total := convertFromBytes(pr.TotalSize, pr.Unit)
	bytesWrited := convertFromBytes(pr.BytesRead, pr.Unit)
	pre := int(float64(bytesWrited) / float64(total) * 100)
	preString := fmt.Sprintf("%d", pre) + "%"
	progressBar := strings.Repeat("=", pre/2)
	spaceNeeded := strings.Repeat(" ", 50-len(progressBar))

	fmt.Print(constant.ColorGreen)
	fmt.Printf("[%s]%s %s %d%s / %d%s", progressBar+spaceNeeded, constant.ColorReset, preString, bytesWrited, pr.Unit, total, pr.Unit)
}

func NewProgress(headers http.Header) *Progress {
	lengthHeader := headers["Content-Length"]
	if len(lengthHeader) <= 0 {
		return nil
	}
	length, err := strconv.ParseInt(lengthHeader[0], 10, 64)
	unit := constant.KbStr
	if length > constant.Mb {
		unit = constant.MbStr
	} else if length > constant.Gb {
		unit = constant.GbStr
	}

	if err != nil {
		return nil
	}
	return &Progress{
		TotalSize: length,
		Unit:      unit,
	}
}

func convertFromBytes(bytes int64, unit string) int64 {
	switch unit {
	case constant.KbStr:
		return bytes / 1000
	case constant.MbStr:
		return bytes / constant.Mb
	case constant.GbStr:
		return bytes / constant.Gb
	}
	return bytes / 1000
}
