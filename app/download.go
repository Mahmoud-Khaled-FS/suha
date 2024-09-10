package app

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Mahmoud-Khaled-FS/suha/downloader"
	"github.com/Mahmoud-Khaled-FS/suha/mime"
	"github.com/Mahmoud-Khaled-FS/suha/utils"
)

// Download file
func (a *App) DownloadFile() error {
	downloader := downloader.Download{
		Fource:    a.Fource,
		Name:      a.Name,
		OutDir:    a.OutDir,
		Url:       a.Url,
		Quality:   a.Quality,
		AudioOnly: a.AudioOnly,
	}
	return downloader.Download()
}

// Check the nam of file from header first then form url
func (a *App) checkName(res *http.Response) {
	fileEx := ""
	contentHeader, ok := res.Header["Content-Type"]
	if ok && len(contentHeader) >= 1 {
		mimeHeader := contentHeader[0]
		mimeHeader = strings.Split(mimeHeader, ";")[0]
		fileEx = mime.GetExt(mimeHeader)
	}
	fileName := ""
	if a.Name != "" {
		fileName = a.Name
	} else {
		fileName = a.fileNameFromHeader(res.Header)
	}
	if fileName == "download" {
		path := strings.Split(res.Request.URL.Path, "/")
		lastPath := path[len(path)-1]
		nameFromUrl, err := url.QueryUnescape(lastPath)
		if err != nil {
			nameFromUrl = lastPath
		}
		if strings.Contains(lastPath, fileEx) && fileEx != ".html" {
			a.Name = nameFromUrl
			return
		}
		fileName = lastPath
	}

	if fileEx != "" && !strings.Contains(fileName, fileEx) {
		fileName += fileEx
	}
	a.Name = fileName
}

// Check if Content-Disposition contain filename or not and return the name
// will return UTF-8 string if name filename is UTF-8
func (a *App) fileNameFromHeader(header http.Header) string {
	dispositionHeaderssd, ok := header["Content-Disposition"]
	if !ok || len(dispositionHeaderssd) == 0 {
		return "download"
	}
	head := dispositionHeaderssd[0]
	nameUtfReg, _ := regexp.Compile(`filename\*=UTF-8''([\w%\-\.]+)(?:; ?|$)`)
	name := strings.Replace(nameUtfReg.FindString(head), "filename*=UTF-8''", "", 1)
	if name != "" {
		decodedName, err := url.QueryUnescape(name)
		if err == nil {
			name = decodedName
		}
	}
	if name != "" {
		return name
	}
	nameAciiReg, _ := regexp.Compile(`filename=(["']?)(.*?[^\\])(?:; ?|$)`)
	name = strings.Replace(nameAciiReg.FindString(head), "filename=\"", "", 1)
	name = strings.Replace(name, "\";", "", 1)
	if name != "" {
		return name
	} else {
		return "download"
	}
}

// Download list urls with DownloadFile()
func (a *App) DownloadList() {
	utils.MoveTop(len(a.UrlList))
	for i, url := range a.UrlList {
		a.Url = url
		a.Name = ""
		utils.ClearText()
		err := a.DownloadFile()
		utils.ClearText()
		if err != nil {
			utils.PrintCopiedUrl(url, i+1, utils.Failed)
		} else {
			utils.PrintCopiedUrl(url, i+1, utils.Success)
		}
		utils.MoveBottom(1)
	}
}
