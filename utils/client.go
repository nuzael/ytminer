package utils

import (
	"fmt"
	"ytminer/platform/ytapi"
)

// CreateYouTubeClient creates a YouTube client with standardized error handling
func CreateYouTubeClient() (ytapi.Client, error) {
	client, err := ytapi.New()
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
