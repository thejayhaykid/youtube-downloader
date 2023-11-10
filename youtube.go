// youtube.go
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetVideoURLs(channelId string) []string {
	ctx := context.Background()

	// Replace with your API key
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create YouTube service: %v", err)
	}

	call := service.Search.List([]string{"id"}).
		ChannelId(channelId).
		MaxResults(50).
		Type("video")

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	videoURLs := make([]string, len(response.Items))
	for i, item := range response.Items {
		videoURLs[i] = "https://www.youtube.com/watch?v=" + item.Id.VideoId
	}

	return videoURLs
}
