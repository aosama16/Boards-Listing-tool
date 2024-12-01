package core

import (
	"boards-merger/internal/utils/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadDirectory(dirPath string, recursive bool, maxDepth int) ([]string, error) {
	var jsonFiles []string

	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get an absolute path for: %v", dirPath)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %v", dirPath)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %v", dirPath)
	}

	logger.Info("Reading directory: %v", absDir)
	if !recursive && maxDepth > 0 {
		maxDepth = 0
		logger.Warn("'Max depth' is set while 'recursive' is false, setting MaxDepth to 0")
	}

	rootDepth := strings.Count(filepath.ToSlash(absDir), "/")
	walkFunc := func(path string, d os.DirEntry, pathError error) error {
		if pathError != nil {
			logger.Warn("Skipping path due to error: %v", path)
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			logger.Warn("Skipping Symbolic link: %v", path)
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get an absolute path for: %v", path)
		}

		if d.IsDir() && absPath != absDir {
			if !recursive {
				return filepath.SkipDir
			}

			currentDepth := strings.Count(filepath.ToSlash(absPath), "/")
			if currentDepth > rootDepth+maxDepth {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() && strings.ToLower(filepath.Ext(d.Name())) == ".json" {
			jsonFiles = append(jsonFiles, absPath)
			logger.Info("Found JSON file: %v", absPath)
		}

		return nil
	}

	if err := filepath.WalkDir(absDir, walkFunc); err != nil {
		return nil, err
	}

	if len(jsonFiles) == 0 {
		return nil, fmt.Errorf("no JSON files found in %v", dirPath)
	}

	return jsonFiles, nil
}
