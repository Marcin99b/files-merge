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
	var failures []CopyFailure
	copyTreeInto(filepath.Join(src, "cameraA"), dest, &copiedFilePaths, &failures)
	copyTreeInto(filepath.Join(src, "cameraB"), dest, &copiedFilePaths, &failures)

	assertFileContent(t, filepath.Join(dest, "photo1.jpg"), "a1")
	assertFileContent(t, filepath.Join(dest, "photo1(1).jpg"), "b1")
	assertFileContent(t, filepath.Join(dest, "Screenshots", "shot.png"), "a-shot")
	assertFileContent(t, filepath.Join(dest, "Screenshots", "shot(1).png"), "b-shot")

	if len(copiedFilePaths) != 4 {
		t.Errorf("len(copiedFilePaths) = %d, want 4 (%v)", len(copiedFilePaths), copiedFilePaths)
	}
	if len(failures) != 0 {
		t.Errorf("len(failures) = %d, want 0 (%v)", len(failures), failures)
	}
}

func TestCopyTreeIntoSkipsUnreadableFilesAndKeepsGoing(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root, file permissions have no effect")
	}

	src := t.TempDir()
	dest := t.TempDir()

	writeSourceFile(t, src, "before.jpg", "before")
	writeSourceFile(t, src, "corrupt.jpg", "unreadable")
	writeSourceFile(t, src, "after.jpg", "after")

	unreadablePath := filepath.Join(src, "corrupt.jpg")
	if err := os.Chmod(unreadablePath, 0o000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(unreadablePath, 0o644) })

	var copiedFilePaths []string
	var failures []CopyFailure
	copyTreeInto(src, dest, &copiedFilePaths, &failures)

	assertFileContent(t, filepath.Join(dest, "before.jpg"), "before")
	assertFileContent(t, filepath.Join(dest, "after.jpg"), "after")

	if _, err := os.Stat(filepath.Join(dest, "corrupt.jpg")); !os.IsNotExist(err) {
		t.Errorf("corrupt.jpg should not exist in the destination, stat err = %v", err)
	}

	if len(failures) != 1 {
		t.Fatalf("len(failures) = %d, want 1 (%v)", len(failures), failures)
	}
	if failures[0].Path != unreadablePath {
		t.Errorf("failures[0].Path = %q, want %q", failures[0].Path, unreadablePath)
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
