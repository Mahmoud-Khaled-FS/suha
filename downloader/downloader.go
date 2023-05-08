package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Mahmoud-Khaled-FS/suha/mime"
	"github.com/Mahmoud-Khaled-FS/suha/utils"
)

type Downloader interface {
	Download() error
}

type Download struct {
	Fource  bool
	Name    string
	OutDir  string
	Url     *url.URL
	Quality string
}

func (d *Download) checkName(res *http.Response) {
	fileEx := ""
	contentHeader, ok := res.Header["Content-Type"]
	if ok && len(contentHeader) >= 1 {
		mimeHeader := contentHeader[0]
		mimeHeader = strings.Split(mimeHeader, ";")[0]
		fileEx = mime.GetExt(mimeHeader)
	}
	fileName := ""
	if d.Name != "" {
		fileName = d.Name
	} else {
		fileName = d.fileNameFromHeader(res.Header)
	}
	if fileName == "download" {
		path := strings.Split(res.Request.URL.Path, "/")
		lastPath := path[len(path)-1]
		nameFromUrl, err := url.QueryUnescape(lastPath)
		if err != nil {
			nameFromUrl = lastPath
		}
		if strings.Contains(lastPath, fileEx) && fileEx != ".html" {
			d.Name = nameFromUrl
			return
		}
		fileName = lastPath
	}

	if fileEx != "" && !strings.Contains(fileName, fileEx) {
		fileName += fileEx
	}
	d.Name = fileName
}

func (a *Download) fileNameFromHeader(header http.Header) string {
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

func (d *Download) Download() error {
	url, err := d.getProvider()
	if err != nil {
		return err
	}
	res, err := http.Get(url)
	if err != nil || res.StatusCode > 299 {
		return fmt.Errorf(fmt.Sprintf("ERROR: url '%s' not correct! can not download\n%s", d.Url.String(), err))
	}
	defer res.Body.Close()
	d.checkName(res)
	path := filepath.Join(d.OutDir, d.Name)
	// Check if file exsits
	if !d.Fource && utils.IsFileExist(path) {
		return fmt.Errorf(fmt.Sprintf("ERROR: File '%s' already exist", path))
	}
	utils.CreateDir(d.OutDir)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("ERROR: Can not create file with name '%s'\n%s", d.Name, err))
	}
	defer file.Close()
	pr := utils.NewProgress(res.Header)
	if pr == nil {
		_, err = io.Copy(file, res.Body)
	} else {
		_, err = io.Copy(file, io.TeeReader(res.Body, pr))
	}
	if err != nil {
		file.Close()
		os.Remove(path)
		return fmt.Errorf(fmt.Sprintf("ERROR: Failed to download file %s\n%s", d.Name, err))
	}
	return nil
}

func (d *Download) getProvider() (string, error) {
	host := d.Url.Host
	if strings.Contains(host, "youtube") {
		return d.StreamYoutubeUrl()
	}
	if strings.Contains(host, "pinterest") {
		return d.GetPintrestUrl()
	}
	return d.Url.String(), nil
}
