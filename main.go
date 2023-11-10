package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
