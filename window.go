package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow
var statusLabel *walk.Label
var streamStatusLabel *walk.Label

func CreateWindow() {
	// Initialize the application
	app := walk.App()
	app.SetProductName("Screen Moments")

	// Create the main window
	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title:    "Screen Moments",
		MinSize:  Size{Width: 500, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "Screenshot",
				Layout: VBox{},
				Children: []Widget{
					PushButton{
						Text: "Capture Screenshot (Ctrl+Alt+F1)",
						OnClicked: func() {
							log.Println("Button clicked")
							// Call the same function used by the hotkey
							CaptureScreenshot()
						},
					},
					Label{
						Text: "Press Ctrl+Alt+F1 to capture a screenshot",
					},
				},
			},
			GroupBox{
				Title:  "RTMP Stream",
				Layout: VBox{},
				Children: []Widget{
					Label{
						Text: "RTMP URL: rtmp://localhost:1935/live",
					},
					Label{
						Text: "Stream Key: screenmoments",
					},
					Label{
						Text: "Set up OBS to stream to the above URL and stream key",
					},
					Label{
						AssignTo: &streamStatusLabel,
						Text:     "RTMP server not started",
					},
				},
			},
			GroupBox{
				Title:  "Status",
				Layout: VBox{},
				Children: []Widget{
					Label{
						AssignTo: &statusLabel,
						Text:     "Ready",
					},
				},
			},
		},
	}).Create(); err != nil {
		log.Fatal(err)
	}

	// Show the window
	mainWindow.Show()
}

func RunMainWindow() {
	// Start the main window loop
	mainWindow.Run()
}

func CloseWindow() {
	if mainWindow != nil {
		mainWindow.Close()
	}
}
