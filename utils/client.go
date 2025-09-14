package utils

import (
	"ytminer/ui"
	"ytminer/youtube"
)

// CreateYouTubeClient creates a YouTube client with standardized error handling
func CreateYouTubeClient() (*youtube.Client, error) {
	client, err := youtube.NewClient()
	if err != nil {
		ui.DisplayError("Failed to create YouTube client: " + err.Error())
		return nil, err
	}
	return client, nil
}

// HandleError handles errors in a standardized way
func HandleError(err error, message string) {
	if err != nil {
		ui.DisplayError(message + ": " + err.Error())
	}
}
