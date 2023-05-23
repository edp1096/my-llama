package main

import (
	"fmt"
	"io"
	"net/http"
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
			if !strings.Contains(file.Name(), "dump_state.bin") {
				result = append(result, file.Name())
			}
		}
	}

	return result, nil
}

func downloadVicuna() {
	url := sampleVicunaWeightsDownloadURL
	saveFileName := sampleVicunaWeightsFileName

	// Create HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Create file
	out, err := os.Create(saveFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer out.Close()

	// Print download progress
	var written int64
	size := resp.ContentLength
	// fmt.Printf("Downloading %.2f GB...\n", float64(size)/(1024^3))
	fmt.Printf("Downloading %.2f GB...\n", float64(size)/(1024*1024*1024))

	// Copy response body to file
	progress := make([]byte, 100000)
	for {
		n, err := resp.Body.Read(progress)
		if n > 0 {
			_, err := out.Write(progress[:n])
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				return
			}
			written += int64(n)

			// Update progress
			percent := float64(written) / float64(size) * 100
			fmt.Printf("\r%.2f%% downloaded", percent)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading HTTP response:", err)
			return
		}
	}

	fmt.Println("\nDownload complete!")
}
