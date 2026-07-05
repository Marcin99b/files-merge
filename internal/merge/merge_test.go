package merge

import (
	"path/filepath"
	"slices"
	"testing"
)

func TestDirectories(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeSourceFile(t, src, "DCIM/photo1.jpg", "a1")
	writeSourceFile(t, src, "DCIM(0)/photo1.jpg", "b1")
	writeSourceFile(t, src, "DCIM(0)/photo2.jpg", "b2")
	writeSourceFile(t, src, "DCIM(1)/photo3.jpg", "c1")
	writeSourceFile(t, src, "Movies/clip.mp4", "m1")

	results, err := Directories(src, dest)
	if err != nil {
		t.Fatalf("Directories: %v", err)
	}

	resultsByFolderName := make(map[string]Result)
	for _, result := range results {
		resultsByFolderName[result.FolderName] = result
	}

	t.Run("does not process duplicate folders as their own group", func(t *testing.T) {
		if _, ok := resultsByFolderName["DCIM(0)"]; ok {
			t.Errorf("DCIM(0) should have been merged into DCIM")
		}
		if _, ok := resultsByFolderName["DCIM(1)"]; ok {
			t.Errorf("DCIM(1) should have been merged into DCIM")
		}
	})

	t.Run("records the duplicate folders merged into DCIM", func(t *testing.T) {
		dcim, ok := resultsByFolderName["DCIM"]
		if !ok {
			t.Fatalf("expected a result for DCIM, got %v", resultsByFolderName)
		}

		want := []string{"DCIM(0)", "DCIM(1)"}
		if !slices.Equal(dcim.DuplicateFolderNames, want) {
			t.Errorf("DCIM.DuplicateFolderNames = %v, want %v", dcim.DuplicateFolderNames, want)
		}
	})

	t.Run("merges file contents, giving clashing names a unique suffix", func(t *testing.T) {
		assertFileContent(t, filepath.Join(dest, "DCIM", "photo1.jpg"), "a1")
		assertFileContent(t, filepath.Join(dest, "DCIM", "photo1(1).jpg"), "b1")
		assertFileContent(t, filepath.Join(dest, "DCIM", "photo2.jpg"), "b2")
		assertFileContent(t, filepath.Join(dest, "DCIM", "photo3.jpg"), "c1")
		assertFileContent(t, filepath.Join(dest, "Movies", "clip.mp4"), "m1")
	})

	t.Run("leaves source files untouched", func(t *testing.T) {
		assertFileContent(t, filepath.Join(src, "DCIM", "photo1.jpg"), "a1")
		assertFileContent(t, filepath.Join(src, "DCIM(0)", "photo1.jpg"), "b1")
	})
}
