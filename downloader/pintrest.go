package downloader

import (
	"net/http"

	"golang.org/x/net/html"
)

func (d *Download) GetPintrestUrl() (string, error) {
	resp, err := http.Get(d.Url.String())
	if err != nil {
		return "", err
	}
	h, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	var imageUrl string
	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					imageUrl = a.Val
					return
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if imageUrl != "" {
				break
			}
			link(c)
		}
	}
	link(h)
	return imageUrl, nil
}
