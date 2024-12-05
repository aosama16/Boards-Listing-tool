package model

import (
	"boards-merger/internal/utils/logger"
	"encoding/json"
	"fmt"
)

type BoardsInfo struct {
	Boards   []Board  `json:"boards"`
	MetaData MetaData `json:"_metadata"`
}

type MetaData struct {
	UniqueVendors int `json:"unique_vendors"`
	TotalBoards   int `json:"total_boards"`
}

func (boardinfo *BoardsInfo) UnmarshalJSON(data []byte) error {
	var tempBoards struct {
		Boards []json.RawMessage `json:"boards"`
	}

	err := json.Unmarshal(data, &tempBoards)
	if err != nil {
		return err
	}

	for _, rawBoard := range tempBoards.Boards {
		var board Board
		if err := json.Unmarshal(rawBoard, &board); err != nil {
			logger.Warn("Skipping board due to board parsing error: %v", err.Error())
		} else {
			boardinfo.Boards = append(boardinfo.Boards, board)
		}
	}

	// Try to unmarshal a single JSON board object
	if len(tempBoards.Boards) == 0 {
		logger.Info("Attempting to parse a single board")
		var singleboard Board
		if err := json.Unmarshal(data, &singleboard); err != nil {
			logger.Error("Failed parsing a single board object: %v", err.Error())
			return fmt.Errorf("failed to parse JSON boards list or a single board object")
		}
		boardinfo.Boards = append(boardinfo.Boards, singleboard)
	}

	return nil
}
