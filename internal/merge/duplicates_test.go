package merge

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFindDuplicateFolders(t *testing.T) {
	dir := t.TempDir()

	mustMkdir(t, dir, "DCIM")
	mustMkdir(t, dir, "DCIM(0)")
	mustMkdir(t, dir, "DCIM(1)")
	mustMkdir(t, dir, "Movies")
	mustMkdir(t, dir, "Movies(0)")
	writeSourceFile(t, dir, "DCIM(9)", "a file, not a directory")

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	duplicates := findDuplicateFolders("DCIM", entries)

	var duplicateNames []string
	for _, duplicate := range duplicates {
		duplicateNames = append(duplicateNames, duplicate.Name())
	}

	want := []string{"DCIM(0)", "DCIM(1)"}
	if !slices.Equal(duplicateNames, want) {
		t.Errorf("findDuplicateFolders(DCIM) = %v, want %v", duplicateNames, want)
	}
}

func mustMkdir(t *testing.T, root, name string) {
	t.Helper()

	if err := os.Mkdir(filepath.Join(root, name), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", name, err)
	}
}
