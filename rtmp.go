package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/mp4"
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
	saveToFile  bool
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

		// Create a file writer if saving is enabled
		var fileWriter *mp4.Muxer
		var recordFile *os.File

		if saveToFile {
			// Create recordings directory if it doesn't exist
			recordingsDir := "recordings"
			if err := os.MkdirAll(recordingsDir, 0755); err != nil {
				log.Println("Error creating recordings directory:", err)
			} else {
				// Create a file with timestamp in the filename
				timestamp := time.Now().Format("2006-01-02_15-04-05")
				filename := filepath.Join(recordingsDir, fmt.Sprintf("screenmoments_%s.mp4", timestamp))

				recordFile, err = os.Create(filename)
				if err != nil {
					log.Println("Error creating recording file:", err)
				} else {
					fileWriter = mp4.NewMuxer(recordFile)

					// Write header
					if err := fileWriter.WriteHeader(streams); err != nil {
						log.Println("Error writing MP4 header:", err)
						recordFile.Close()
						fileWriter = nil
					} else {
						log.Println("Recording to file:", filename)
						updateStreamStatus(fmt.Sprintf("Recording to %s", filename))
					}
				}
			}
		}

		// Copy packets one by one to both queue and file
		for {
			pkt, err := conn.ReadPacket()
			if err != nil {
				log.Println("Stream ended:", err)
				break
			}

			// Write to queue
			if err := videoStream.WritePacket(pkt); err != nil {
				log.Println("Error writing to queue:", err)
			}

			// Write to file if recording
			if fileWriter != nil {
				if err := fileWriter.WritePacket(pkt); err != nil {
					log.Println("Error writing to file:", err)
				}
			}
		}

		// Close file if we were recording
		if fileWriter != nil {
			if err := fileWriter.WriteTrailer(); err != nil {
				log.Println("Error writing trailer:", err)
			}
			recordFile.Close()
			log.Println("Recording saved")
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
		err := avutil.CopyFile(conn, cursor)
		if err != nil {
			log.Println("On playing, Error copying file:", err)
		}
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

// ToggleStreamRecording enables or disables saving the stream to a file
func ToggleStreamRecording(enabled bool) {
	mutex.Lock()
	defer mutex.Unlock()
	saveToFile = enabled

	if saveToFile {
		log.Println("Stream recording enabled")
	} else {
		log.Println("Stream recording disabled")
	}
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
