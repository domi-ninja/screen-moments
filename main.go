package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Define hotkey constants
const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008

	// Virtual key codes
	VK_F1 = 0x70
	VK_F2 = 0x71
	VK_A  = 0x41
	VK_B  = 0x42
	// Add more as needed

	// Windows messages
	WM_HOTKEY = 0x0312
)

// Windows API functions
var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
)

// RegisterHotKey wraps the Windows RegisterHotKey function
func RegisterHotKey(hwnd windows.Handle, id int32, fsModifiers, vk uint32) bool {
	ret, _, _ := procRegisterHotKey.Call(
		uintptr(hwnd),
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
	)
	return ret != 0
}

// UnregisterHotKey wraps the Windows UnregisterHotKey function
func UnregisterHotKey(hwnd windows.Handle, id int32) bool {
	ret, _, _ := procUnregisterHotKey.Call(
		uintptr(hwnd),
		uintptr(id),
	)
	return ret != 0
}

// MSG represents a Windows message
type MSG struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// GetMessage wraps the Windows GetMessage function
func GetMessage(msg *MSG, hwnd windows.Handle, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int(ret)
}

// TranslateMessage wraps the Windows TranslateMessage function
func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

// DispatchMessage wraps the Windows DispatchMessage function
func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(uintptr(unsafe.Pointer(msg)))
	return ret
}

func main() {
	// Register a hotkey (Ctrl+Alt+F1 in this example)
	hotkeyID := 1
	if !RegisterHotKey(0, int32(hotkeyID), MOD_CONTROL|MOD_ALT, VK_F1) {
		fmt.Println("Failed to register hotkey")
		os.Exit(1)
	}

	fmt.Println("Hotkey registered: Ctrl+Alt+F1")
	fmt.Println("Press Ctrl+C to exit")

	// Message loop to capture hotkey events
	go func() {
		var msg MSG
		for {
			result := GetMessage(&msg, 0, 0, 0)
			fmt.Println("result", result)
			if result == 0 {
				// WM_QUIT message
				return
			}
			if result == -1 {
				// Error
				return
			}

			if msg.Message == WM_HOTKEY {
				fmt.Println("Hotkey triggered!")
				// Add your code to handle the hotkey here
			}
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}()

	// Wait for Ctrl+C to exit
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Unregister the hotkey before exiting
	UnregisterHotKey(0, int32(hotkeyID))
	fmt.Println("Hotkey unregistered. Exiting...")
}
