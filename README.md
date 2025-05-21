# ScreenMoments

A simple Windows application that demonstrates global hotkey registration using Go.

## Features

- Registers a global hotkey (Ctrl+Alt+F1)
- Performs an action when the hotkey is triggered
- Clean exit with Ctrl+C

## Requirements

- Go 1.18 or later
- Windows operating system

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