# A Channel Downloader Written in Go

## How to run yourself

1. Install Go on their machine if it's not already installed. They can download it from the official Go website.

2. Clone the repository to their local machine using Git. They can do this by running the following command in their terminal:

```bash
git clone https://github.com/thejayhaykid/youtube-downloader.git
```

Navigate into the cloned repository's directory:

```bash
cd youtube-downloader
```

Run the program with the desired YouTube channel ID as an argument:

```bash
go run main.go CHANNEL_ID
```

They should replace `CHANNEL_ID` with the actual ID of the YouTube channel.

Please note that this program requires a YouTube Data API key to function. The key should be stored in a `.env` file in the root of the project with the variable name YOUTUBE_DATA_API_KEY. Use the `.env.example` as an example of what to put in there.
