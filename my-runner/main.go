package main // import "my-runner"

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	target := "./my-llama.exe"

	// check if my-llama.exe exists
	if _, err := os.Stat(target); err != nil {
		target = "./my-llama_cu.exe"
	}

	for {
		// run my-llama.exe
		cmd := exec.Command(target)
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}
}
