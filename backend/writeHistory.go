package backend

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func writeHistoryToFile(filename string, history []FeedContentFilterInfo) error {
	if history == nil {
		history = []FeedContentFilterInfo{}
	}

	historyMap := make(map[string]FeedContentFilterInfo)

	// Calculate the hash of the history
	for _, item := range history {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(item.Title+item.Content)))
		historyMap[hash] = item
	}

	data, err := json.MarshalIndent(historyMap, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func WriteHistory(history []FeedContentFilterInfo) error {
	if history == nil {
		history = []FeedContentFilterInfo{}
	}

	var historyFilePath string
	if os.Getenv("DEV_MODE") == "true" {
		historyFilePath = "data/history.json"
	} else {
		configDir, _ := os.UserConfigDir()
		historyFilePath = filepath.Join(configDir, "MrRSS", "data", "history.json")
	}

	err := writeHistoryToFile(historyFilePath, history)
	if err != nil {
		return err
	}

	return nil
}
