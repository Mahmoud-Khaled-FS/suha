package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func main() {
	// TODO: Name flag
	// nameFlag := flag.String("n", "", "File name")
	args := os.Args
	if len(args) == 1 {
		panicWithUsage("ERROR: No Url Founded!")
	}
	urlPath := string(args[1])
	flag.Parse()
	if urlPath == "" {
		panicWithUsage("ERROR: Url should be the first argument")
	}
	u, err := url.ParseRequestURI(urlPath)
	if err != nil {
		panicWithUsage(fmt.Sprintf("ERROR: expected url but got '%s'", urlPath))
	}
	downloadFile("fille.png", u)
}

func downloadFile(fileName string, url *url.URL) {
	file, err := os.Create(fileName)
	if err != nil {
		panicRed(fmt.Sprintf("ERROR: Can not create file with name '%s'\n%s", fileName, err))
	}
	defer file.Close()
	res, err := http.Get(url.String())
	fmt.Println(res.StatusCode)
	if err != nil || res.StatusCode > 299 {
		file.Close()
		os.Remove(fileName)
		panicRed(fmt.Sprintf("ERROR: url '%s' not correct! can not download\n%s", url.String(), err))
	}
	defer res.Body.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		file.Close()
		os.Remove(fileName)
		panicRed(fmt.Sprintf("ERROR: Failed to download file %s\n%s", fileName, err))
	}
	greenMsg("Download Success!")
}

func panicWithUsage(msg string) {
	fmt.Println("usage suha <url>")
	panicRed(msg)
}

func panicRed(msg string) {
	fmt.Print(colorRed)
	fmt.Println(msg)
	fmt.Print(colorReset)
	os.Exit(1)
}
func greenMsg(msg string) {
	fmt.Print(colorGreen)
	fmt.Println(msg)
	fmt.Print(colorReset)
}
