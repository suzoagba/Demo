package handlers

import (
	"log"
	"os/exec"
	"runtime"
)

func Open(path string) { // https://github.com/0x434D53/openinbrowser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open", path}
	case "windows":
		args = []string{"cmd", "/c", "start", path}
	default:
		args = []string{"xdg-open", path}
	}
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Printf("openinbrowser: %v, this function is not working in Docker. Open link manually.", err)
	}
}
