package merge

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUniqueDestPath(t *testing.T) {
	dir := t.TempDir()
	writeSourceFile(t, dir, "photo.jpg", "existing")

	t.Run("returns a numbered suffix when the file already exists on disk", func(t *testing.T) {
		got := uniqueDestPath(dir, "photo.jpg", nil)
		want := filepath.Join(dir, "photo(1).jpg")
		if got != want {
			t.Errorf("uniqueDestPath() = %q, want %q", got, want)
		}
	})

	t.Run("returns the next numbered suffix when one is already claimed in this run", func(t *testing.T) {
		alreadyCopied := []string{filepath.Join(dir, "photo(1).jpg")}
		got := uniqueDestPath(dir, "photo.jpg", alreadyCopied)
		want := filepath.Join(dir, "photo(2).jpg")
		if got != want {
			t.Errorf("uniqueDestPath() = %q, want %q", got, want)
		}
	})

	t.Run("returns the original name when there is no clash", func(t *testing.T) {
		got := uniqueDestPath(dir, "new.jpg", nil)
		want := filepath.Join(dir, "new.jpg")
		if got != want {
			t.Errorf("uniqueDestPath() = %q, want %q", got, want)
		}
	})
}

func TestCopyTreeIntoMergesNestedNameClashes(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	writeSourceFile(t, src, "cameraA/photo1.jpg", "a1")
	writeSourceFile(t, src, "cameraA/Screenshots/shot.png", "a-shot")
	writeSourceFile(t, src, "cameraB/photo1.jpg", "b1")
	writeSourceFile(t, src, "cameraB/Screenshots/shot.png", "b-shot")

	var copiedFilePaths []string
	if err := copyTreeInto(filepath.Join(src, "cameraA"), dest, &copiedFilePaths); err != nil {
		t.Fatalf("copyTreeInto cameraA: %v", err)
	}
	if err := copyTreeInto(filepath.Join(src, "cameraB"), dest, &copiedFilePaths); err != nil {
		t.Fatalf("copyTreeInto cameraB: %v", err)
	}

	assertFileContent(t, filepath.Join(dest, "photo1.jpg"), "a1")
	assertFileContent(t, filepath.Join(dest, "photo1(1).jpg"), "b1")
	assertFileContent(t, filepath.Join(dest, "Screenshots", "shot.png"), "a-shot")
	assertFileContent(t, filepath.Join(dest, "Screenshots", "shot(1).png"), "b-shot")

	if len(copiedFilePaths) != 4 {
		t.Errorf("len(copiedFilePaths) = %d, want 4 (%v)", len(copiedFilePaths), copiedFilePaths)
	}
}

func writeSourceFile(t *testing.T, root, relativePath, content string) {
	t.Helper()

	fullPath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", relativePath, err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if string(got) != want {
		t.Errorf("content of %s = %q, want %q", path, got, want)
	}
}
