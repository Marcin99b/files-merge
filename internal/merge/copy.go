package merge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

func copyTreeInto(srcDir, destDir string, copiedFilePaths *[]string) error {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", destDir, err)
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())

		if entry.IsDir() {
			if err := copyTreeInto(srcPath, filepath.Join(destDir, entry.Name()), copiedFilePaths); err != nil {
				return err
			}
			continue
		}

		if err := copyFileInto(srcPath, destDir, entry.Name(), copiedFilePaths); err != nil {
			return err
		}
	}

	return nil
}

func copyFileInto(srcPath, destDir, fileName string, copiedFilePaths *[]string) error {
	destPath := uniqueDestPath(destDir, fileName, *copiedFilePaths)

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("opening %s: %w", srcPath, err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", destPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copying %s to %s: %w", srcPath, destPath, err)
	}

	*copiedFilePaths = append(*copiedFilePaths, destPath)
	return nil
}

func uniqueDestPath(destDir, fileName string, copiedFilePaths []string) string {
	destPath := filepath.Join(destDir, fileName)
	if !isPathTaken(destPath, copiedFilePaths) {
		return destPath
	}

	ext := filepath.Ext(fileName)
	base := fileName[:len(fileName)-len(ext)]

	for suffix := 1; ; suffix++ {
		candidate := filepath.Join(destDir, fmt.Sprintf("%s(%d)%s", base, suffix, ext))
		if !isPathTaken(candidate, copiedFilePaths) {
			return candidate
		}
	}
}

func isPathTaken(path string, copiedFilePaths []string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return slices.Contains(copiedFilePaths, path)
}
