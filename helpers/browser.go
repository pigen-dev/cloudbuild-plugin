package helpers

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd string
	var args []string
	fmt.Println("Opening browser with URL:", url)
	switch {
	case isWindows():
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case isMac():
		cmd = "open"
		args = []string{url}
	case isLinux():
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

// Platform detection helpers
func isWindows() bool {
	return runtime.GOOS == "windows"
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}