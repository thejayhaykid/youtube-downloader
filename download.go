package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Andreychik32/ytdl"
	"github.com/rs/zerolog"
)

func downloadList(ctx context.Context, videoURLs []string) (err error) {
	client := ytdl.Client{
		HTTPClient: nil,
		Logger:     zerolog.Nop(),
	}

	if len(videoURLs) == 0 {
		return fmt.Errorf("no video URLs provided")
	}

	if ctx == nil {
		return fmt.Errorf("no context provided")
	}

	for _, videoURL := range videoURLs {
		videoInfo, err := client.GetVideoInfo(ctx, videoURL)
		if err != nil {
			fmt.Printf("Failed to get video info for %s: %v\n", videoURL, err)
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
			fmt.Printf("No audio-only format found for %s\n", videoURL)
			continue
		}

		// Create a safe filename
		datePublished := videoInfo.DatePublished.Format("20060102")
		title := strings.ReplaceAll(videoInfo.Title, "/", "-")
		filename := fmt.Sprintf("%s_%s", datePublished, title)

		// Download the audio
		audioFile, err := os.Create(filename + ".mp3")
		if err != nil {
			fmt.Printf("Failed to create file for %s: %v\n", videoURL, err)
			continue
		}
		defer audioFile.Close()

		ctx := context.Background() // Create a new context
		err = client.Download(ctx, videoInfo, audioFormat, audioFile)
		if err != nil {
			fmt.Printf("Failed to download audio for %s: %v\n", videoURL, err)
			continue
		}

		// Save the video info to a JSON file
		infoFile, err := os.Create(filename + ".json")
		if err != nil {
			fmt.Printf("Failed to create info file for %s: %v\n", videoURL, err)
			continue
		}
		defer infoFile.Close()

		err = json.NewEncoder(infoFile).Encode(videoInfo)
		if err != nil {
			fmt.Printf("Failed to write info file for %s: %v\n", videoURL, err)
			continue
		}
	}

	return nil
}
