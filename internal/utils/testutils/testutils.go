package testutils

import (
	"boards-merger/internal/model"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	BoolTrue  = true
	BoolFalse = false
)

func WriteToFile(t *testing.T, filePath string, data string) {
	t.Helper()
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatalf("failed to create file %v: %v", filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		t.Fatalf("failed to write string to file %v: %v", filePath, err)
	}
}

func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	return dir
}

func CreateTempFile(t *testing.T, dir string, filename string) string {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create temporary file %v: %v", filePath, err)
	}
	defer file.Close()
	return filePath
}

func CompareBoards(t *testing.T, expected model.Board, actual model.Board) {
	t.Helper()
	if expected.Name != actual.Name {
		t.Fatalf("Unexpected board name: got %v, expected %v", actual.Name, expected.Name)
	}

	if expected.Vendor != actual.Vendor {
		t.Fatalf("Unexpected board vendor: got %v, expected %v", actual.Vendor, expected.Vendor)
	}

	if expected.Core != actual.Core {
		t.Fatalf("Unexpected board core: got %v, expected %v", actual.Core, expected.Core)
	}

	if (expected.HasWiFi == nil) != (actual.HasWiFi == nil) {
		t.Fatalf("Unexpected board has_wifi references: got 'board.has_wifi' == nil is %v, expected 'board.has_wifi' == nil is %v", (expected.HasWiFi == nil), (actual.HasWiFi == nil))
	}

	if (len(expected.ExtraEntries) == 0) && (len(actual.ExtraEntries) == 0) {
		return
	}

	if !reflect.DeepEqual(expected.ExtraEntries, actual.ExtraEntries) {
		t.Fatalf("Unexpected board extra entries: got %v, expected %v", actual.ExtraEntries, expected.ExtraEntries)
	}
}

func CompareBoardsInfo(t *testing.T, expected *model.BoardsInfo, actual *model.BoardsInfo) {
	t.Helper()
	if (expected == nil) && (actual == nil) {
		return
	}

	if (expected == nil) != (actual == nil) {
		t.Fatalf("Unexpected boards references: got 'board == nil' is %v, expected 'board == nil' is %v", (actual == nil), (expected == nil))
	}

	if len(expected.Boards) != len(actual.Boards) {
		t.Fatalf("Unexpected boards length: got %v, expected %v", len(actual.Boards), len(expected.Boards))
	}

	for i := range expected.Boards {
		CompareBoards(t, expected.Boards[i], actual.Boards[i])
	}
}
