package utils

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ShowLoading displays a simple loading animation with the given message
func ShowLoading(message string) func() {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Margin(1, 0)
	
	// Show loading message
	fmt.Printf("\n%s\n", style.Render(message))
	
	// Create animation frames
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	// Start animation in goroutine
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, frame := range frames {
					fmt.Printf("\r%s %s", frame, message)
					time.Sleep(100 * time.Millisecond)
					select {
					case <-done:
						return
					default:
					}
				}
			}
		}
	}()
	
	// Return function to stop loading
	return func() {
		done <- true
		fmt.Printf("\r%s %s\n", 
			lipgloss.NewStyle().Foreground(lipgloss.Color("#27AE60")).Render("✅"), 
			lipgloss.NewStyle().Bold(true).Render("Complete!"))
	}
}
