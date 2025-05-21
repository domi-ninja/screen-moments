package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
)

const (
	rtmpAddr = ":1935"         // Default RTMP port
	rtmpApp  = "live"          // Application name for RTMP URL
	rtmpKey  = "screenmoments" // Stream key for RTMP
)

var (
	rtmpServer  *rtmp.Server
	videoStream *pubsub.Queue
	mutex       sync.Mutex
	isStreaming bool
)

// StartRTMPServer initializes and starts an RTMP server
func StartRTMPServer() {
	log.Println("Starting RTMP server on", rtmpAddr)

	// Create a new RTMP server
	rtmpServer = &rtmp.Server{
		Addr: rtmpAddr,
	}

	// Create a pubsub queue for the video stream
	videoStream = pubsub.NewQueue()
	videoStream.SetMaxGopCount(1) // Only keep the latest GOP (Group of Pictures)

	// Handle incoming RTMP streams (from OBS)
	rtmpServer.HandlePublish = func(conn *rtmp.Conn) {
		mutex.Lock()
		isStreaming = true
		mutex.Unlock()

		// Update UI to show we're receiving
		updateStreamStatus("Receiving stream from OBS")

		streams, err := conn.Streams()
		if err != nil {
			log.Println("Error getting streams:", err)
			return
		}

		// Log information about the received streams
		log.Printf("Received %d streams from OBS", len(streams))
		for i, stream := range streams {
			log.Printf("Stream #%d: Type=%v", i, stream.Type())
		}

		// Copy the stream to our queue
		err = avutil.CopyFile(videoStream, conn)
		if err != nil {
			log.Println("Stream ended with error:", err)
		}

		mutex.Lock()
		isStreaming = false
		mutex.Unlock()

		// Update UI when stream ends
		updateStreamStatus("Stream ended")
	}

	// Handle RTMP play requests (for potential playback)
	rtmpServer.HandlePlay = func(conn *rtmp.Conn) {
		cursor := videoStream.Latest()
		avutil.CopyFile(conn, cursor)
	}

	// Start the RTMP server in a goroutine
	go func() {
		err := rtmpServer.ListenAndServe()
		if err != nil {
			log.Println("RTMP server error:", err)
			updateStreamStatus(fmt.Sprintf("RTMP server error: %v", err))
		}
	}()

	updateStreamStatus("RTMP server started. Ready for OBS stream.")
}

// StopRTMPServer logs that the RTMP server is stopping
// Note: The joy4 RTMP server doesn't provide a way to gracefully shut down
// The server will automatically close when the application exits
func StopRTMPServer() {
	if rtmpServer != nil {
		log.Println("RTMP server stopping. It will close automatically when the application exits.")
	}
}

// IsStreamActive returns whether a stream is currently active
func IsStreamActive() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return isStreaming
}

// updateStreamStatus updates the UI with the current streaming status
func updateStreamStatus(status string) {
	if mainWindow != nil && streamStatusLabel != nil {
		// This will run on the UI thread
		mainWindow.Synchronize(func() {
			streamStatusLabel.SetText(status)
		})
	}
	log.Println(status)
}
