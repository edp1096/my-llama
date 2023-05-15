package main // import "github.com/edp1096/my-llama"

import (
	"flag"
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v3/cpu"
)

var (
	threads = 4
)

var (
	cpuPhysicalNUM = 0
	cpuLogicalNUM  = 0

	isBrowserOpen = false

	modelPath = "./"

	modelFname  string   = ""
	modelFnames []string = []string{}
)

func main() {
	cpuPhysicalNUM, _ = cpu.Counts(false)
	cpuLogicalNUM, _ = cpu.Counts(true)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&isBrowserOpen, "b", false, "open browser automatically")

	threads = cpuPhysicalNUM

	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	if modelFname == "" {
		modelFnames, err = findModelFiles(modelPath)
		if err != nil {
			fmt.Printf("Finding model files failed: %s", err)
			os.Exit(1)
		}
		// fmt.Println(modelFnames)

		if len(modelFnames) == 0 {
			fmt.Println("No model files found.")
			fmt.Println("Press enter to download vicuna model file and to open the model search page.")
			fmt.Println("Press Ctrl+C, if you want to exit.")
			fmt.Scanln()

			openBrowser(weightsSearchURL)
			downloadVicuna()

			modelFnames, _ = findModelFiles(modelPath)
		}

		modelFname = modelFnames[0]
	}

	if _, err := os.Stat(modelFname); os.IsNotExist(err) {
		// Because, model will be downloaded if not exists, may be not reachable here
		fmt.Printf("Model file %s does not exist", modelFname)
		os.Exit(1)
	}

	fmt.Println("CPU cores physical/logical/threads:", cpuPhysicalNUM, "/", cpuLogicalNUM, "/", threads)
}
