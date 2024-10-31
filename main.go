package main

import (
	"fmt"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	robotgo "github.com/go-vgo/robotgo"
	gohook "github.com/robotn/gohook"
)

var (
	recording   bool
	actions     []string
	lastKeyTime time.Time
	myApp       = app.New() // Declare myApp as a global variable
)

func main() {
	myWindow := myApp.NewWindow("Keyboard Recorder")

	// Create a label to show status
	statusLabel := widget.NewLabel("Press F2 to start recording, F10 to stop, and F4 to play.")

	// Create a button to save actions
	saveButton := widget.NewButton("Save Actions", func() {
		saveActions(myWindow) // Pass the window to the saveActions function
	})

	// Set the content of the window
	myWindow.SetContent(container.NewVBox(
		statusLabel,
		saveButton,
	))

	// Set size for the window
	myWindow.Resize(fyne.NewSize(800, 600))

	// Start listening for global keyboard events
	go listenForGlobalKeys()

	myWindow.ShowAndRun()
}

// listenForGlobalKeys listens for global keyboard events and processes them.
func listenForGlobalKeys() {
	// Start the hook
	chanHook := gohook.Start()
	defer gohook.End()

	for ev := range chanHook {
		if ev.Kind == gohook.KeyHold {
			fmt.Printf("Key pressed: %d\n", ev.Keycode)
			handleKeyHold(int(ev.Keycode))
		}
	}
}

// handleKeyHold processes the key hold events and manages recording/playback.
func handleKeyHold(keycode int) {
	currentTime := time.Now()

	if recording {
		// Calculate delay since the last key press
		delay := 0
		if !lastKeyTime.IsZero() {
			delay = int(currentTime.Sub(lastKeyTime).Milliseconds())
		}
		// Store the keycode and delay in a simpler format
		actions = append(actions, fmt.Sprintf("%d,%d", keycode, delay))
		fmt.Printf("Recorded key: %d, Delay: %d ms\n", keycode, delay)
	}

	// Update the last key time
	lastKeyTime = currentTime

	// Check for specific keys to control recording and playback
	switch keycode {
	case 60: // F2 key
		startRecording()
	case 68: // F10 key
		stopRecording()
	case 62: // F4 key
		playRecording()
	}
}

// startRecording initializes the recording process.
func startRecording() {
	recording = true
	actions = []string{}      // Clear previous actions
	lastKeyTime = time.Time{} // Reset last key time
	fmt.Println("Recording started...")
}

// stopRecording ends the recording process.
func stopRecording() {
	recording = false
	fmt.Println("Recording stopped.")
}

// playRecording plays back the recorded actions with the recorded delays.
func playRecording() {
	recording = false
	fmt.Println("Playing back recorded actions...")

	for _, action := range actions {
		var keycode, delay int
		fmt.Sscanf(action, "%d,%d", &keycode, &delay) // Extract keycode and delay

		// Simulate the delay for playback
		time.Sleep(time.Duration(delay) * time.Millisecond)

		// Simulate the key press using robotgo
		robotgo.KeyTap(getKeyString(keycode)) // Convert keycode to string

		// Display the action (simulating pressing the key)
		fmt.Printf("Simulated Key Press: %d\n", keycode)
	}

	fmt.Println("Playback finished.")
}

// getKeyString converts a keycode to its corresponding string representation.
func getKeyString(keycode int) string {
	for k, v := range gohook.Keycode {
		if int(v) == keycode {
			return k
		}
	}
	return ""
}

// saveActions prompts the user to choose a file location to save the recorded actions.
// saveActions prompts the user to choose a file location to save the recorded actions.
func saveActions(window fyne.Window) {
	dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}
		defer uc.Close()

		for _, action := range actions {
			_, err := uc.Write([]byte(action + "\n"))
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
		}
		fmt.Println("Actions saved successfully.")
	}, window)
}
