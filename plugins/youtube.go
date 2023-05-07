package main

import (
	"fmt"

	"github.com/kkdai/youtube/v2"
)

func main() {
	client := youtube.Client{}
	video, err := client.GetVideo("ugS4FsrW0fM")
	fmt.Println(video.Formats.WithAudioChannels()[1].Quality)
	fmt.Println(video.Formats.WithAudioChannels()[1].QualityLabel)
	fmt.Println(video.Formats.WithAudioChannels()[1].AudioQuality)
	// stream, _, _ := client.GetStream(video, &video.Formats[0])

	// f, _ := os.Create("video.mp4")

	// io.Copy(f, stream)

	fmt.Println(err)
}
