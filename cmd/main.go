package main

import (
	"boards-merger/internal/core"
	"boards-merger/internal/utils/logger"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	logger.Enable()

	testpath := "./test/example-boards"
	// testpath := "./test/sandbox"
	// testpath := "./test/conflict"
	// testpath := "./test/symlinks"
	jsonList, err := core.ReadDirectory(testpath, false, 10)
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
