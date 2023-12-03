package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	// "github.com/Andreychik32/ytdl"
	"github.com/joho/godotenv"
	ytdl "github.com/kkdai/youtube/v2"
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

	// Create a subfolder
	err = os.MkdirAll("downloads", 0755)
	if err != nil {
		fmt.Printf("Failed to create subfolder: %v\n", err)
		os.Exit(1)
	}

	err = DownloadList(ctx, videoURLs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func DownloadList(ctx context.Context, videoURLs []string) (err error) {
	client := ytdl.Client{}

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
		video, err := client.GetVideo(videoURL)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to get video info for %s: %v\n", videoURL, err)
			continue
		}

		format := video.Formats.FindByItag(140) // Find the audio-only stream (itag 140 is for m4a audio)
		if format == nil {
			fmt.Fprintf(errorFile, "No audio stream found for %s\n", videoURL)
			continue
		}

		stream, _, err := client.GetStream(video, format)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to get audio stream for %s: %v\n", videoURL, err)
			continue
		}

		// Create a safe filename by replacing spaces with underscores and removing special characters
		safeTitle := strings.Map(func(r rune) rune {
			if r == ' ' {
				return '_'
			} else if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' {
				return -1
			}
			return r
		}, video.Title)

		publishDate, err := getPublishDate(video.ID)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to get publish date for %s: %v\n", videoURL, err)
			continue
		}

		filename := fmt.Sprintf("downloads/%s_%s", publishDate.Format("2006-01-02"), safeTitle)

		file, err := os.Create(filename + ".m4a") // Save the audio stream to a .m4a file
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to create file for %s: %v\n", videoURL, err)
			continue
		}
		defer file.Close()

		_, err = io.Copy(file, stream)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to download audio for %s: %v\n", videoURL, err)
			continue
		}

		// Save metadata to a separate file
		metadataFile, err := os.Create(filename + ".json")
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to create metadata file for %s: %v\n", videoURL, err)
			continue
		}
		defer metadataFile.Close()

		metadataEncoder := json.NewEncoder(metadataFile)
		err = metadataEncoder.Encode(video)
		if err != nil {
			fmt.Fprintf(errorFile, "Failed to write metadata for %s: %v\n", videoURL, err)
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

func convertToEmbedURL(youtubeURL string) string {
	// Extract the video ID from the YouTube URL
	// Assumes the URL is in the format "https://www.youtube.com/watch?v=VIDEO_ID"
	videoID := strings.Split(youtubeURL, "=")[1]

	// Return the embed URL
	return "https://www.youtube.com/embed/" + videoID
}

func getPublishDate(videoID string) (time.Time, error) {
	// Replace with your API key
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoID, apiKey))
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	var data struct {
		Items []struct {
			Snippet struct {
				PublishedAt time.Time `json:"publishedAt"`
			} `json:"snippet"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return time.Time{}, err
	}

	if len(data.Items) == 0 {
		return time.Time{}, errors.New("no video found")
	}

	return data.Items[0].Snippet.PublishedAt, nil
}
