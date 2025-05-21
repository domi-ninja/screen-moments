package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var mainWindow *walk.MainWindow
var statusLabel *walk.Label

func CreateWindow() {
	// Initialize the application
	app := walk.App()
	app.SetProductName("Screen Moments")

	// Create the main window
	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title:    "Screen Moments",
		MinSize:  Size{Width: 400, Height: 300},
		Layout:   VBox{},
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
			Label{
				AssignTo: &statusLabel,
				Text:     "Ready",
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
