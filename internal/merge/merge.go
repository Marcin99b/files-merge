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

		copiedFilePaths, err := mergeFolderGroup(sourceRoot, destinationRoot, entry.Name(), duplicateFolderNames)
		if err != nil {
			return results, err
		}

		results = append(results, Result{
			FolderName:           entry.Name(),
			DuplicateFolderNames: duplicateFolderNames,
			CopiedFilePaths:      copiedFilePaths,
		})
	}

	return results, nil
}

func mergeFolderGroup(sourceRoot, destinationRoot, folderName string, duplicateFolderNames []string) ([]string, error) {
	destDir := filepath.Join(destinationRoot, folderName)

	var copiedFilePaths []string
	if err := copyTreeInto(filepath.Join(sourceRoot, folderName), destDir, &copiedFilePaths); err != nil {
		return nil, err
	}

	for _, duplicateFolderName := range duplicateFolderNames {
		if err := copyTreeInto(filepath.Join(sourceRoot, duplicateFolderName), destDir, &copiedFilePaths); err != nil {
			return nil, err
		}
	}

	return copiedFilePaths, nil
}
