package merge

import (
	"os"
	"regexp"
)

func findDuplicateFolders(name string, entries []os.DirEntry) []os.DirEntry {
	isDuplicateName := duplicateFolderPattern(name).MatchString

	var duplicates []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == name {
			continue
		}
		if isDuplicateName(entry.Name()) {
			duplicates = append(duplicates, entry)
		}
	}

	return duplicates
}

func duplicateFolderPattern(name string) *regexp.Regexp {
	return regexp.MustCompile(`^` + regexp.QuoteMeta(name) + `\(\d+\)$`)
}
