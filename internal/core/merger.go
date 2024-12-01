package core

import (
	"boards-merger/internal/model"
	"boards-merger/internal/utils/logger"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type boardRegistry map[string]model.Board
type void struct{}
type stringSet map[string]void

func hashBoard(vendor string, name string) string {
	return fmt.Sprintf("%s::%s", vendor, name)
}

func sortedKeys(m boardRegistry) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func ProcessJsonFiles(jsonFilePaths []string) (*model.BoardsInfo, error) {
	var boardsMap = make(boardRegistry)
	var vendorSet = make(stringSet)
	var boardsInfo model.BoardsInfo

	for _, path := range jsonFilePaths {
		jsonFile, err := os.ReadFile(path)
		if err != nil {
			logger.Error("Failed to read the JSON file, skipping file '%v'", path)
			fmt.Println(err.Error())
			continue
		}
		logger.Info("Parsing file: %v", path)

		var boardsList model.BoardsInfo
		if err := json.Unmarshal(jsonFile, &boardsList); err != nil {
			logger.Error("%v, skipping file '%v'", err.Error(), path)
			continue
		}

		for _, board := range boardsList.Boards {
			boardHash := hashBoard(board.Vendor, board.Name)
			if existingBoard, exists := boardsMap[boardHash]; exists {
				logger.Warn("Found a duplicate entry for board '%v' made by '%v', Attempting to merge them", board.Name, board.Vendor)
				board.Merge(existingBoard)
			}
			boardsMap[boardHash] = board

			if _, exists := vendorSet[board.Vendor]; !exists {
				vendorSet[board.Vendor] = void{}
			}
		}
	}

	if len(boardsMap) == 0 {
		return nil, fmt.Errorf("no valid boards found")
	}

	sortekBoardKeys := sortedKeys(boardsMap)
	boardsInfo.Boards = make([]model.Board, 0, len(boardsMap))
	for _, key := range sortekBoardKeys {
		boardsInfo.Boards = append(boardsInfo.Boards, boardsMap[key])
	}

	boardsInfo.MetaData = model.MetaData{
		UniqueVendors: len(vendorSet),
		TotalBoards:   len(boardsMap),
	}

	return &boardsInfo, nil
}
