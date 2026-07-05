package merge

import (
	"fmt"
	"os"
	"path/filepath"
)

type Result struct {
	FolderName           string
	DuplicateFolderNames []string
	CopiedFilePaths      []string
	Failures             []CopyFailure
}

func Directories(sourcePath, destinationPath string) ([]Result, error) {
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("reading source directory %s: %w", sourcePath, err)
	}

	sourceRoot, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("resolving source directory %s: %w", sourcePath, err)
	}

	destinationRoot, err := filepath.Abs(destinationPath)
	if err != nil {
		return nil, fmt.Errorf("resolving destination directory %s: %w", destinationPath, err)
	}

	processedFolderNames := make(map[string]bool)
	var results []Result

	for _, entry := range entries {
		if !entry.IsDir() || processedFolderNames[entry.Name()] {
			continue
		}

		duplicates := findDuplicateFolders(entry.Name(), entries)

		processedFolderNames[entry.Name()] = true
		duplicateFolderNames := make([]string, 0, len(duplicates))
		for _, duplicate := range duplicates {
			processedFolderNames[duplicate.Name()] = true
			duplicateFolderNames = append(duplicateFolderNames, duplicate.Name())
		}

		copiedFilePaths, failures := mergeFolderGroup(sourceRoot, destinationRoot, entry.Name(), duplicateFolderNames)

		results = append(results, Result{
			FolderName:           entry.Name(),
			DuplicateFolderNames: duplicateFolderNames,
			CopiedFilePaths:      copiedFilePaths,
			Failures:             failures,
		})
	}

	return results, nil
}

func mergeFolderGroup(sourceRoot, destinationRoot, folderName string, duplicateFolderNames []string) ([]string, []CopyFailure) {
	destDir := filepath.Join(destinationRoot, folderName)

	var copiedFilePaths []string
	var failures []CopyFailure

	copyTreeInto(filepath.Join(sourceRoot, folderName), destDir, &copiedFilePaths, &failures)

	for _, duplicateFolderName := range duplicateFolderNames {
		copyTreeInto(filepath.Join(sourceRoot, duplicateFolderName), destDir, &copiedFilePaths, &failures)
	}

	return copiedFilePaths, failures
}
