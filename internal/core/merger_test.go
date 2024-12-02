package core_test

import (
	"boards-merger/internal/core"
	"boards-merger/internal/model"
	"boards-merger/internal/utils/logger"
	"boards-merger/internal/utils/testutils"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessJsonFiles(t *testing.T) {
	logger.Disable()

	tests := []struct {
		name               string
		setup              func(t *testing.T) ([]string, func())
		expectedErr        bool
		expectedBoardsInfo *model.BoardsInfo
	}{
		{
			name: "Valid JSON files with unique boards and extra entries",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `{
					"boards": [
						{"name": "Board2", "vendor": "VendorB", "has_wifi": true, "extra_feature_1": "yes"},
						{"name": "Board1", "vendor": "VendorA", "has_wifi": false, "extra_feature_1": "no", "extra_feature_2": "yes"}
					]
				}`
				testutils.WriteToFile(t, filePath, jsonContent)
				return []string{filePath}, func() { os.RemoveAll(dir) }
			},
			expectedErr: false,
			expectedBoardsInfo: &model.BoardsInfo{
				Boards: []model.Board{
					{
						Name:    "Board1",
						Vendor:  "VendorA",
						HasWiFi: &testutils.BoolFalse,
						ExtraEntries: map[string]interface{}{
							"extra_feature_1": "no",
							"extra_feature_2": "yes",
						},
					},
					{
						Name:    "Board2",
						Vendor:  "VendorB",
						HasWiFi: &testutils.BoolTrue,
						ExtraEntries: map[string]interface{}{
							"extra_feature_1": "yes",
						},
					},
				},
				MetaData: model.MetaData{
					UniqueVendors: 2,
					TotalBoards:   2,
				},
			},
		},
		{
			name: "Valid JSON files with unique single boards",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `{"name": "Board2", "vendor": "VendorB", "has_wifi": true}`
				testutils.WriteToFile(t, filePath, jsonContent)

				filePath2 := filepath.Join(dir, "boards-2.json")
				jsonContent2 := `{"name": "Board1", "vendor": "VendorA", "has_wifi": false}`
				testutils.WriteToFile(t, filePath2, jsonContent2)
				return []string{filePath, filePath2}, func() { os.RemoveAll(dir) }
			},
			expectedErr: false,
			expectedBoardsInfo: &model.BoardsInfo{
				Boards: []model.Board{
					{
						Name:    "Board1",
						Vendor:  "VendorA",
						HasWiFi: &testutils.BoolFalse,
					},
					{
						Name:    "Board2",
						Vendor:  "VendorB",
						HasWiFi: &testutils.BoolTrue,
					},
				},
				MetaData: model.MetaData{
					UniqueVendors: 2,
					TotalBoards:   2,
				},
			},
		},
		{
			name: "Empty Boards",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `{
                    "boards": [
                    ]
                }`
				testutils.WriteToFile(t, filePath, jsonContent)
				return []string{filePath}, func() { os.RemoveAll(dir) }
			},
			expectedErr:        true,
			expectedBoardsInfo: nil,
		},
		{
			name: "Invalid JSON",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `}`
				testutils.WriteToFile(t, filePath, jsonContent)
				return []string{filePath}, func() { os.RemoveAll(dir) }
			},
			expectedErr:        true,
			expectedBoardsInfo: nil,
		},
		{
			name: "Invalid JSON file path",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				return []string{filePath}, func() { os.RemoveAll(dir) }
			},
			expectedErr:        true,
			expectedBoardsInfo: nil,
		},
		{
			name: "Invalid Board",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `{
                    "boards": [
                        {"name": "Board1", "vendor": "VendorA", "has_wifi": false},
                        {"NAME": "Board1", "vendor": "VendorA", "has_wifi": false}
                    ]
                }`
				testutils.WriteToFile(t, filePath, jsonContent)
				return []string{filePath}, func() { os.RemoveAll(dir) }
			},
			expectedErr: false,
			expectedBoardsInfo: &model.BoardsInfo{
				Boards: []model.Board{
					{
						Name:    "Board1",
						Vendor:  "VendorA",
						HasWiFi: &testutils.BoolFalse,
					},
				},
				MetaData: model.MetaData{
					UniqueVendors: 1,
					TotalBoards:   1,
				},
			},
		},
		{
			name: "Duplicate boards with conflict",
			setup: func(t *testing.T) ([]string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := filepath.Join(dir, "boards-1.json")
				jsonContent := `{
                    "boards": [
                        {"name": "Board2", "vendor": "VendorB", "has_wifi": true},
						{"name": "Board1", "vendor": "VendorA", "core": "CoreY", "has_wifi": false, "extra_feature_1": "yes"}
                    ]
                }`
				testutils.WriteToFile(t, filePath, jsonContent)
				filePath2 := filepath.Join(dir, "boards-2.json")
				jsonContent2 := `{
                    "boards": [
                        {"name": "Board3", "vendor": "VendorB", "core": "CoreZ"},
						{"name": "Board1", "vendor": "VendorA", "core": "CoreX", "has_wifi": false, "extra_feature_2": "no"}
                    ]
                }`
				testutils.WriteToFile(t, filePath2, jsonContent2)

				return []string{filePath, filePath2}, func() { os.RemoveAll(dir) }
			},
			expectedErr: false,
			expectedBoardsInfo: &model.BoardsInfo{
				Boards: []model.Board{
					{
						Name:    "Board1",
						Vendor:  "VendorA",
						Core:    "CoreX",
						HasWiFi: &testutils.BoolFalse,
						ExtraEntries: map[string]interface{}{
							"extra_feature_1": "yes",
							"extra_feature_2": "no",
						},
					},
					{
						Name:    "Board2",
						Vendor:  "VendorB",
						HasWiFi: &testutils.BoolTrue,
					},
					{
						Name:   "Board3",
						Vendor: "VendorB",
						Core:   "CoreZ",
					},
				},
				MetaData: model.MetaData{
					UniqueVendors: 2,
					TotalBoards:   3,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paths, cleanup := test.setup(t)
			defer cleanup()

			actualBoardsInfo, err := core.ProcessJsonFiles(paths)

			if (err != nil) != test.expectedErr {
				t.Fatalf("unexpected error status: got %v, expected error: %v", err, test.expectedErr)
			}

			testutils.CompareBoardsInfo(t, test.expectedBoardsInfo, actualBoardsInfo)
		})
	}
}
