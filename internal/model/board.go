package model

import (
	"boards-merger/internal/utils/logger"
	"encoding/json"
	"fmt"
	"strings"
)

type Board struct {
	Name         string
	Vendor       string
	Core         string
	HasWiFi      *bool
	ExtraEntries map[string]interface{}
}

func sanitizeMapKeys(data map[string]interface{}) map[string]interface{} {
	sanitizedData := make(map[string]interface{})
	for key, value := range data {
		sanitizedKey := strings.TrimSpace(key)
		sanitizedData[sanitizedKey] = value
	}
	return sanitizedData
}

func (board *Board) UnmarshalJSON(data []byte) error {
	var rawMap map[string]interface{}
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return err
	}

	// Trim leading and trailing spaces from keys
	sanRawMap := sanitizeMapKeys(rawMap)

	// Board 'name' is required
	name, exists := sanRawMap["name"].(string)
	name = strings.TrimSpace(name)
	if !exists || len(name) == 0 {
		return fmt.Errorf("JSON object missing critical data: Board name")
	}
	board.Name = name
	delete(sanRawMap, "name")

	// Board 'vendor' is required
	vendor, exists := sanRawMap["vendor"].(string)
	vendor = strings.TrimSpace(vendor)
	if !exists || len(vendor) == 0 {
		return fmt.Errorf("JSON object missing critical data: Board vendor")
	}
	board.Vendor = vendor
	delete(sanRawMap, "vendor")

	// Board 'core' is optional
	core, exists := sanRawMap["core"].(string)
	core = strings.TrimSpace(core)
	if exists && len(strings.TrimSpace(core)) > 0 {
		board.Core = core
		delete(sanRawMap, "core")
	}

	// Board 'has_wifi' is optional
	hasWiFi, exists := sanRawMap["has_wifi"].(bool)
	if exists {
		board.HasWiFi = &hasWiFi
		delete(sanRawMap, "has_wifi")
	}

	// Preserve all extra properties
	board.ExtraEntries = sanRawMap

	return nil
}

func (board Board) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{
		"name":   board.Name,
		"vendor": board.Vendor,
	}

	if board.Core != "" {
		result["core"] = board.Core
	}

	if board.HasWiFi != nil {
		result["has_wifi"] = *board.HasWiFi
	}

	for key, value := range board.ExtraEntries {
		result[key] = value
	}

	return json.Marshal(result)
}

func (board *Board) Merge(other Board) error {
	if board.Name != other.Name || board.Vendor != other.Vendor {
		return fmt.Errorf("cannot merge boards with different name or vendor")
	}

	if other.Core != "" {
		if board.Core == "" {
			board.Core = other.Core
		} else if board.Core != other.Core {
			logger.Warn("Two entries for %v boards have conflicting info about 'Core': {'%v', '%v'}, choosing '%v'", board.Name, board.Core, other.Core, board.Core)
		}
	}

	if other.HasWiFi != nil {
		if board.HasWiFi == nil {
			board.HasWiFi = other.HasWiFi
		} else if *board.HasWiFi != *other.HasWiFi {
			logger.Warn("Two entries for %v boards have conflicting info about 'Has Wifi', choosing '%v'", board.Name, *board.HasWiFi)
		}
	}

	if board.ExtraEntries == nil && other.ExtraEntries != nil {
		board.ExtraEntries = make(map[string]interface{})
	}
	for key, value := range other.ExtraEntries {
		// Conflicting values for the same key will be overridden
		board.ExtraEntries[key] = value
	}

	return nil
}
