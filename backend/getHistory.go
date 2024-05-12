package backend

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type HistoryInfo map[string]FeedContentFilterInfo

func readHistoryFromFile(filename string) (HistoryInfo, error) {
	// Check if file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return HistoryInfo{}, nil
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	var history HistoryInfo
	err = json.Unmarshal(file, &history)
	if err != nil {
		var emptyArray []interface{}
		err2 := json.Unmarshal(file, &emptyArray)
		if err2 == nil && len(emptyArray) == 0 {
			return HistoryInfo{}, nil
		}

		log.Fatalf("Unable to parse JSON: %v", err)
	}

	return history, nil
}

func GetHistory() HistoryInfo {
	var historyFilePath string
	if os.Getenv("DEV_MODE") == "true" {
		historyFilePath = "data/history.json"
	} else {
		configDir, _ := os.UserConfigDir()
		historyFilePath = filepath.Join(configDir, "MrRSS", "data", "history.json")
	}
	historyList, err := readHistoryFromFile(historyFilePath)
	if err != nil {
		log.Fatalf("Unable to read history: %v", err)
	}
	return historyList
}
