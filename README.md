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

---

## Getting a YouTube Data API Key

To get a YouTube Data API key, you need to create a project in the Google Cloud Console and enable the YouTube Data API v3 for that project. Here are the steps:

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).

2. Click on "Select a project" at the top of the page, then click on "New Project", and give your project a name.

3. Once the project is created, select it and go to the "Library" page in the API & Services section.

4. Search for "YouTube Data API v3" and click on it, then click "Enable".

5. After the API is enabled, go to the "Credentials" page in the API & Services section.

6. Click on "Create Credentials", then select "API key".

7. A new API key will be created. You can restrict the key's usage to the YouTube Data API v3 and your project by clicking on "Restrict key" and setting the appropriate options.

Remember to secure your API key as it can be used to perform any API request on behalf of your application.
