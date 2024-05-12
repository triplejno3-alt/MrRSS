package backend

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type FeedInfo struct {
	Title string
	Link  string
}

func readFeedFromFile(filename string) []FeedInfo {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	var feeds []FeedInfo
	err = json.Unmarshal(file, &feeds)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %v", err)
	}

	return feeds
}

func GetFeedList() []FeedInfo {
	var feedFilePath string
	if os.Getenv("DEV_MODE") == "true" {
		feedFilePath = "data/feeds.json"
	} else {
		configDir, _ := os.UserConfigDir()
		feedFilePath = filepath.Join(configDir, "MrRSS", "data", "feeds.json")
	}
	feedList := readFeedFromFile(feedFilePath)
	return feedList
}
