# Screen Moments

A simple application for capturing screenshots and receiving RTMP video streams from OBS.

## Features

- Screenshot capture using hotkey (Ctrl+Alt+F1)
- RTMP server to receive video streams from OBS

## RTMP Streaming from OBS

To stream video from OBS to Screen Moments:

1. Open OBS Studio
2. Go to Settings -> Stream
3. Select "Custom..." as the service
4. Set the Server to: `rtmp://localhost:1935/live`
5. Set the Stream Key to: `screenmoments`
6. Click "OK" to save settings
7. Click "Start Streaming" in OBS

The Screen Moments application will automatically receive and display the RTMP stream status.

## Building from Source

```bash
# Get dependencies
go mod tidy

# Run the application
go run .

# Build the application
go build .
```

## Requirements

- Go 1.20 or later
- OBS Studio (for RTMP streaming)

## Installation

1. Clone the repository:
```
git clone https://github.com/yourusername/screenmoments.git
cd screenmoments
```

2. Install dependencies:
```
go mod tidy
```

3. Build the application:
```
go build
```

## Usage

Run the application:
```
./screenmoments.exe
```

Press Ctrl+Alt+F1 to trigger the registered hotkey.
Press Ctrl+C to exit the application.

## Customization

To change the hotkey combination, modify the `main.go` file and adjust the following line:
```go
if !win.RegisterHotKey(0, hotkeyID, MOD_CONTROL|MOD_ALT, VK_F1) {
```

Available modifiers:
- MOD_ALT (Alt key)
- MOD_CONTROL (Ctrl key)
- MOD_SHIFT (Shift key)
- MOD_WIN (Windows key)

Virtual key codes are defined as constants in the code. 