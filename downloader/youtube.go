package downloader

import (
	"fmt"

	"github.com/kkdai/youtube/v2"
)

func (d *Download) StreamYoutubeUrl() (string, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(d.Url.String())
	if err != nil {
		return "", err
	}
	d.Name = video.Title
	streamUrl := ""
	if d.AudioOnly {
		for _, format := range video.Formats {
			if format.AudioQuality != "" && format.QualityLabel == "" {
				streamUrl, err = client.GetStreamURL(video, &format)
				if err != nil {
					return "", err
				}
				break
			}
		}
	} else {
		quality := d.Quality
		if quality == "" {
			quality = "360p"
		}
		formats := video.Formats.WithAudioChannels().Quality(quality)
		if len(formats) == 0 {
			return "", fmt.Errorf("ERROR: No formats found with quality %s", quality)
		}
		streamUrl, err = client.GetStreamURL(video, &formats[0])
		if err != nil {
			return "", err
		}
	}
	return streamUrl, nil
}
