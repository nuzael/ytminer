package utils

import (
	"fmt"
	"ytminer/youtube"
)

// CreateYouTubeClient creates a YouTube client with standardized error handling
func CreateYouTubeClient() (*youtube.Client, error) {
	client, err := youtube.NewClient()
	if err != nil {
		fmt.Printf("Failed to create YouTube client: %v\n", err)
		return nil, err
	}
	return client, nil
}

// HandleError handles errors in a standardized way
func HandleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	}
}
