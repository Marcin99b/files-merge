package merge

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type CopyFailure struct {
	Path string
	Err  error
}

func copyTreeInto(srcDir, destDir string, copiedFilePaths *[]string, failures *[]CopyFailure) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		*failures = append(*failures, CopyFailure{Path: srcDir, Err: fmt.Errorf("creating directory %s: %w", destDir, err)})
		return
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		*failures = append(*failures, CopyFailure{Path: srcDir, Err: fmt.Errorf("reading directory %s: %w", srcDir, err)})
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())

		if entry.IsDir() {
			copyTreeInto(srcPath, filepath.Join(destDir, entry.Name()), copiedFilePaths, failures)
			continue
		}

		copyFileInto(srcPath, destDir, entry.Name(), copiedFilePaths, failures)
	}
}

func copyFileInto(srcPath, destDir, fileName string, copiedFilePaths *[]string, failures *[]CopyFailure) {
	destPath := uniqueDestPath(destDir, fileName, *copiedFilePaths)

	src, err := os.Open(srcPath)
	if err != nil {
		*failures = append(*failures, CopyFailure{Path: srcPath, Err: fmt.Errorf("opening %s: %w", srcPath, err)})
		return
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		*failures = append(*failures, CopyFailure{Path: srcPath, Err: fmt.Errorf("creating %s: %w", destPath, err)})
		return
	}

	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()

	if copyErr != nil {
		_ = os.Remove(destPath)
		*failures = append(*failures, CopyFailure{Path: srcPath, Err: fmt.Errorf("copying %s to %s: %w", srcPath, destPath, copyErr)})
		return
	}
	if closeErr != nil {
		_ = os.Remove(destPath)
		*failures = append(*failures, CopyFailure{Path: srcPath, Err: fmt.Errorf("closing %s: %w", destPath, closeErr)})
		return
	}

	*copiedFilePaths = append(*copiedFilePaths, destPath)
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
