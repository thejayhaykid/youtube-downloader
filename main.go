package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Andreychik32/ytdl"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a channel ID as an argument")
		os.Exit(1)
	}

	channelID := os.Args[1]

	// Replace with the actual list of video URLs
	videoURLs := GetVideoURLs(channelID)

	ctx := context.Background()

	err = DownloadList(ctx, videoURLs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func DownloadList(ctx context.Context, videoURLs []string) (err error) {
	client := &ytdl.Client{
		HTTPClient: &http.Client{},
	}

	if len(videoURLs) == 0 {
		return fmt.Errorf("no video URLs provided")
	}

	if ctx == nil {
		return fmt.Errorf("no context provided")
	}

	errorFile, err := os.Create("errors.txt")
	if err != nil {
		return fmt.Errorf("failed to create error file: %w", err)
	}
	defer errorFile.Close()

	for _, videoURL := range videoURLs {
		fmt.Println("Downloading", videoURL)
		videoInfo, err := client.GetVideoInfo(ctx, videoURL)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to get video info for %s: %v\n", videoURL, err)
			continue
		}

		// Find an audio-only format
		var audioFormat *ytdl.Format
		for _, format := range videoInfo.Formats {
			if format.AudioEncoding != "" && format.VideoEncoding == "" {
				audioFormat = format
				break
			}
		}

		if audioFormat == nil {
			fmt.Fprintf(errorFile, "No audio-only format found for %s\n", videoURL)
			continue
		}

		// Create a safe filename
		datePublished := videoInfo.DatePublished.Format("20060102")
		title := strings.ReplaceAll(videoInfo.Title, "/", "-")
		filename := fmt.Sprintf("%s_%s", datePublished, title)

		// Download the audio
		audioFile, err := os.Create(filename + ".mp3")
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to create file for %s: %v\n", videoURL, err)
			continue
		}
		defer audioFile.Close()

		ctx := context.Background() // Create a new context
		err = client.Download(ctx, videoInfo, audioFormat, audioFile)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to download audio for %s: %v\n", videoURL, err)
			continue
		}

		// Save the video info to a JSON file
		infoFile, err := os.Create(filename + ".json")
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to create info file for %s: %v\n", videoURL, err)
			continue
		}
		defer infoFile.Close()

		err = json.NewEncoder(infoFile).Encode(videoInfo)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to write info file for %s: %v\n", videoURL, err)
			continue
		}
	}

	return nil
}

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
