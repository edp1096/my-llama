package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func openBrowser(uri string) {
	switch os := runtime.GOOS; os {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", uri).Start()
	case "linux":
		exec.Command("xdg-open", uri).Start()
	case "darwin":
		exec.Command("open", uri).Start()
	default:
		fmt.Printf("%s: unsupported platform", os)
	}
}

func findModelFiles(path string) (result []string, err error) {
	result = []string{}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".bin") {
			result = append(result, file.Name())
		}
	}

	return result, nil
}
