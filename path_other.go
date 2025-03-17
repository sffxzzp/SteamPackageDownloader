//go:build linux || darwin

package main

import (
	"runtime"
)

func getSteamPath() string {
	path := ""
	switch runtime.GOOS {
	case "darwin":
		path = "Steam.app/contents/macOS/"
		break
	case "linux":
		path = "~/.local/share/Steam/"
		break
	}
	return path
}
