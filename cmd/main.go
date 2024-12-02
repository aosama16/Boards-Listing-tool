package main

import (
	"boards-merger/internal/core"
	"boards-merger/internal/utils/logger"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	dirPathFlag := flag.String("path", "", "Path to the directory containing JSON files")
	recursiveFlag := flag.Bool("r", false, "Enable recursive directory traversal (default: disabled)")
	loggingFlag := flag.Bool("l", false, "Enable logs (default: disabled)")
	depthFlag := flag.Int("depth", 10, "Maximum depth for directory traversal, used only when recursive is set")
	flag.Parse()

	dirPath := *dirPathFlag
	if len(dirPath) == 0 {
		fmt.Print("Enter the path to the directory: ")
		fmt.Scanln(&dirPath)
	}

	recursive := *recursiveFlag
	depth := *depthFlag
	if !recursive {
		depth = 0
	}

	if *loggingFlag {
		logger.Enable()
	} else {
		logger.Disable()
	}

	jsonList, err := core.ReadDirectory(dirPath, recursive, depth)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	boards, err := core.ProcessJsonFiles(jsonList)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	out, _ := json.MarshalIndent(boards, "", "  ")
	os.Stdout.Write(out)
	fmt.Println()
}
