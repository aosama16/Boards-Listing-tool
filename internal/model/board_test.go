package model_test

import (
	"boards-merger/internal/model"
	"boards-merger/internal/utils/logger"
	"boards-merger/internal/utils/testutils"
	"encoding/json"
	"reflect"
	"testing"
)

func TestBoardUnmarshal(t *testing.T) {
	logger.Disable()

	tests := []struct {
		name          string
		setup         string
		expectedErr   bool
		expectedBoard model.Board
	}{
		{
			name: "Valid JSON all fields and extra",
			setup: `{
                "name":   "Board1",
                "vendor": "VendorA",
                "has_wifi": false,
                "core": "CoreX",
                "extra_feature_1": "no",
                "extra_feature_2": "yes"
			}`,
			expectedErr: false,
			expectedBoard: model.Board{
				Name:    "Board1",
				Vendor:  "VendorA",
				Core:    "CoreX",
				HasWiFi: &testutils.BoolFalse,
				ExtraEntries: map[string]interface{}{
					"extra_feature_1": "no",
					"extra_feature_2": "yes",
				},
			},
		},
		{
			name:          "Empty JSON",
			setup:         `{ }`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
		{
			name: "Invalid JSON",
			setup: `{ {
                "name":   "Board1",
                "vendor": "VendorA",
                "has_wifi": false,
                "core": "CoreX",
                "extra_feature_1": "no",
                "extra_feature_2": "yes"
                }`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
		{
			name: "Missing name",
			setup: `{
                "vendor": "VendorA"
            }`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
		{
			name: "Missing name capitalization",
			setup: `{
                "NAME":   "Board1",
                "vendor": "VendorA"
			}`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
		{
			name: "Missing vendor",
			setup: `{
                "name":   "Board1"
			}`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
		{
			name: "Missing vendor capitalization",
			setup: `{
                "name":   "Board1",
                "VENDOR": "VendorA"
			}`,
			expectedErr:   true,
			expectedBoard: model.Board{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var board model.Board
			err := json.Unmarshal([]byte(test.setup), &board)

			if (err != nil) != test.expectedErr {
				t.Fatalf("Unexpected err: %v", err.Error())
			}

			if !test.expectedErr {
				testutils.CompareBoards(t, test.expectedBoard, board)
			}
		})
	}
}

func TestBoardMarshal(t *testing.T) {
	board := model.Board{
		Name:    "Board1",
		Vendor:  "VendorA",
		Core:    "CoreX",
		HasWiFi: nil,
		ExtraEntries: map[string]interface{}{
			"extra_feature_1": "no",
			"extra_feature_2": "yes",
		},
	}

	expectedMap := map[string]interface{}{
		"name":            "Board1",
		"vendor":          "VendorA",
		"core":            "CoreX",
		"extra_feature_1": "no",
		"extra_feature_2": "yes",
	}

	data, err := json.Marshal(board)
	if err != nil {
		t.Fatalf("Unexpected marshalling err: %v", err.Error())
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(data, &resultMap)
	if err != nil {
		t.Fatalf("Unexpected unmarshalling err: %v", err.Error())
	}

	if !reflect.DeepEqual(expectedMap, resultMap) {
		t.Errorf("unexpected JSON output: got %s, expected %s", resultMap, expectedMap)
	}
}

func TestBoardMerge(t *testing.T) {
	board1 := model.Board{
		Name:    "Board1",
		Vendor:  "VendorA",
		HasWiFi: &testutils.BoolFalse,
		ExtraEntries: map[string]interface{}{
			"extra_feature_1": "no",
		},
	}
	board2 := model.Board{
		Name:    "Board1",
		Vendor:  "VendorA",
		Core:    "CoreX",
		HasWiFi: nil,
		ExtraEntries: map[string]interface{}{
			"extra_feature_2": "yes",
		},
	}
	board3 := model.Board{
		Name:   "Board1",
		Vendor: "VendorB",
	}

	expectedBoard := model.Board{
		Name:    "Board1",
		Vendor:  "VendorA",
		Core:    "CoreX",
		HasWiFi: &testutils.BoolFalse,
		ExtraEntries: map[string]interface{}{
			"extra_feature_1": "no",
			"extra_feature_2": "yes",
		},
	}

	err := board1.Merge(board3)
	if err == nil {
		t.Fatalf("Merge should have failed boards have different name/vendor:"+
			"(Board1 name: %v, Board1 vendor: %v), (Board3 name: %v, Board3 vendor: %v)", board1.Name, board1.Vendor, board3.Name, board3.Vendor)
	}

	err = board1.Merge(board2)
	if err != nil {
		t.Fatalf("Unexpected err: %v", err.Error())
	}

	testutils.CompareBoards(t, expectedBoard, board1)
}
