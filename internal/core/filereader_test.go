package core_test

import (
	"boards-merger/internal/core"
	"boards-merger/internal/utils/logger"
	"boards-merger/internal/utils/testutils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadDirectory(t *testing.T) {
	logger.Disable()

	tests := []struct {
		name        string
		setup       func(t *testing.T) (string, func())
		recursive   bool
		maxDepth    int
		expectedErr bool
		expectedLen int
	}{
		{
			name: "Directory with JSON files",
			setup: func(t *testing.T) (string, func()) {
				dir := testutils.CreateTempDir(t)
				testutils.CreateTempFile(t, dir, "boards-1.json")
				testutils.CreateTempFile(t, dir, "boards-2.JSON")
				testutils.CreateTempFile(t, dir, "boards-3.Json")
				return dir, func() { os.RemoveAll(dir) }
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: false,
			expectedLen: 3,
		},
		{
			name: "Empty directory",
			setup: func(t *testing.T) (string, func()) {
				dir := testutils.CreateTempDir(t)
				return dir, func() { os.RemoveAll(dir) }
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "Not a directory",
			setup: func(t *testing.T) (string, func()) {
				dir := testutils.CreateTempDir(t)
				filePath := testutils.CreateTempFile(t, dir, "boards-3.Json")
				return filePath, func() { os.RemoveAll(dir) }
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "Invalid path",
			setup: func(t *testing.T) (string, func()) {
				return "Invalid path", func() {}
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "Directory with non JSON files",
			setup: func(t *testing.T) (string, func()) {
				dir := testutils.CreateTempDir(t)
				testutils.CreateTempFile(t, dir, "file-1.jsonn")
				testutils.CreateTempFile(t, dir, "file-2.jjson")
				testutils.CreateTempFile(t, dir, "json.go")
				testutils.CreateTempFile(t, dir, "file-4.txt")
				return dir, func() { os.RemoveAll(dir) }
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "Recursive directory",
			setup: func(t *testing.T) (string, func()) {
				rootDir := testutils.CreateTempDir(t)
				subDir := filepath.Join(rootDir, "subdir")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdirectory: %v", err)
				}
				testutils.CreateTempFile(t, rootDir, "boards-1.json")
				testutils.CreateTempFile(t, subDir, "boards-2.json")
				return rootDir, func() { os.RemoveAll(rootDir) }
			},
			recursive:   true,
			maxDepth:    1,
			expectedErr: false,
			expectedLen: 2,
		},
		{
			name: "Non recursive directory with sub-directory",
			setup: func(t *testing.T) (string, func()) {
				rootDir := testutils.CreateTempDir(t)
				subDir := filepath.Join(rootDir, "subdir")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdirectory: %v", err)
				}
				testutils.CreateTempFile(t, rootDir, "boards-1.json")
				testutils.CreateTempFile(t, subDir, "boards-2.json")
				return rootDir, func() { os.RemoveAll(rootDir) }
			},
			recursive:   false,
			maxDepth:    10,
			expectedErr: false,
			expectedLen: 1,
		},
		{
			name: "No permission for root directory",
			setup: func(t *testing.T) (string, func()) {
				rootDir := testutils.CreateTempDir(t)
				os.Chmod(rootDir, 0000)
				return rootDir, func() {
					os.Chmod(rootDir, 0755)
					os.RemoveAll(rootDir)
				}
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: true,
			expectedLen: 0,
		},
		{
			name: "Recursive directory with limited depth",
			setup: func(t *testing.T) (string, func()) {
				rootDir := testutils.CreateTempDir(t)
				subDir := filepath.Join(rootDir, "subdir")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatalf("failed to create subdirectory: %v", err)
				}
				testutils.CreateTempFile(t, rootDir, "boards-1.json")
				testutils.CreateTempFile(t, subDir, "boards-2.json")
				return rootDir, func() { os.RemoveAll(rootDir) }
			},
			recursive:   true,
			maxDepth:    0,
			expectedErr: false,
			expectedLen: 1,
		},
		{
			name: "Directory with symlink",
			setup: func(t *testing.T) (string, func()) {
				rootDir := testutils.CreateTempDir(t)
				testutils.CreateTempFile(t, rootDir, "file1.json")
				symlinkDir := filepath.Join(rootDir, "symlink")
				if err := os.Symlink(rootDir, symlinkDir); err != nil {
					t.Fatalf("failed to create symlink: %v", err)
				}
				return rootDir, func() { os.RemoveAll(rootDir) }
			},
			recursive:   false,
			maxDepth:    0,
			expectedErr: false,
			expectedLen: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir, cleanup := test.setup(t)
			defer cleanup()

			jsonFileList, err := core.ReadDirectory(dir, test.recursive, test.maxDepth)

			if (err != nil) != test.expectedErr {
				t.Fatalf("Unexpected err: %v", err.Error())
			}

			if len(jsonFileList) != test.expectedLen {
				t.Fatalf("Unexpected JSON file list length: got %v, expected %v", len(jsonFileList), test.expectedLen)
			}

			for _, jsonFile := range jsonFileList {
				jsonFile = strings.ToLower(jsonFile)

				if !strings.HasSuffix(jsonFile, ".json") {
					t.Errorf("Unexpected file: %v", jsonFile)
				}
			}
		})
	}
}
