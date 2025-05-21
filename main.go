package main

func main() {
	// Create and display the application window
	CreateWindow()

	// Start the hotkey listener in a separate goroutine
	go StartHotkeyListener()

	// Run the main window
	RunMainWindow()
}
