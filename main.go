package main

func main() {
	// Create and display the application window
	CreateWindow()

	// Start the RTMP server in a separate goroutine
	go StartRTMPServer()

	// Start the hotkey listener in a separate goroutine
	go StartHotkeyListener()

	// Run the main window
	RunMainWindow()

	// Cleanup when application exits
	StopRTMPServer()
}
