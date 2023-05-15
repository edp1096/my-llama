package main // import "github.com/edp1096/my-llama"

import (
	"flag"
	"fmt"
	"os"

	"github.com/edp1096/my-llama/cgollama"
	"github.com/shirou/gopsutil/v3/cpu"
)

var (
	cpuPhysicalNUM = 0
	cpuLogicalNUM  = 0

	isBrowserOpen = false

	modelPath = "./"

	modelFilename  string   = ""
	modelFilenames []string = []string{}
)

func main() {
	l, err := cgollama.New()
	if err != nil {
		fmt.Printf("Creating LLama instance failed: %s", err)
		os.Exit(1)
	}

	l.Hello()

	cpuPhysicalNUM, _ = cpu.Counts(false)
	cpuLogicalNUM, _ = cpu.Counts(true)

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&isBrowserOpen, "b", false, "open browser automatically")

	err = flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Parsing program arguments failed: %s", err)
		os.Exit(1)
	}

	modelFilenames, err = findModelFiles(modelPath)
	if err != nil {
		fmt.Printf("Finding model files failed: %s", err)
		os.Exit(1)
	}
	// fmt.Println(modelFnames)

	if len(modelFilenames) == 0 {
		fmt.Println("No model files found.")
		fmt.Println("Press enter to download vicuna model file and open the model search page.")
		fmt.Println("Press Ctrl+C, if you want to exit.")
		fmt.Scanln()

		openBrowser(weightsSearchURL)
		downloadVicuna()

		modelFilenames, _ = findModelFiles(modelPath)
	}
	modelFilename = modelFilenames[0]

	if _, err := os.Stat(modelFilename); os.IsNotExist(err) {
		// Because, model will be downloaded if not exists, may be not reachable here
		fmt.Printf("Model file %s does not exist", modelFilename)
		os.Exit(1)
	}

	fmt.Println("CPU cores physical/logical:", cpuPhysicalNUM, "/", cpuLogicalNUM)
}
