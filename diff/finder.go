package diff

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type diffFinder struct {
	logger Logger
}

var finder = newDiffFinder()

// newDiffFinder creates a new diffFinder instance
func newDiffFinder() *diffFinder {
	return &diffFinder{
		logger: newLogger("finder"),
	}
}

// TryFindExe tries to find the diff tool at all the specified paths
func (f *diffFinder) TryFindExe(paths []string) (exePath string, found bool) {
	ps := f.unique(paths)
	for _, p := range ps {
		if r, found := f.TryFind(p); found {
			return r, found
		}
	}
	return "", false
}

// TryFind tries to find the diff tool at the specified paths
func (f *diffFinder) TryFind(path string) (result string, found bool) {
	expanded := os.ExpandEnv(path)
	if !strings.ContainsRune(expanded, '*') {
		if utils.File.Exists(expanded) {
			result = expanded
			found = true
			return
		}

		f.logger.Info("could not find file: %s", path)
		return "", false
	}

	var filePart = filepath.Base(expanded)
	var directoryPart = utils.File.GetDirectoryName(expanded)
	var directories = f.getDirectories(directoryPart)

	for _, dir := range directories {
		if strings.ContainsRune(dir, '*') {
			panic("wildcard in file part currently not supported.")
		}

		var filePath = filepath.Join(dir, filePart)
		if utils.File.Exists(filePath) {
			return filePath, true
		}
	}

	f.logger.Info("could not find file: %s", path)
	return "", false
}

func (f *diffFinder) unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (f *diffFinder) getDirectories(directory string) []string {
	var expanded = os.ExpandEnv(directory)
	if !strings.ContainsRune(expanded, '*') {
		dirExists, _ := utils.File.FileOrDirectoryExists(directory)
		if dirExists {
			return []string{directory}
		}
	}

	var segments = strings.Split(expanded, string(filepath.Separator))
	var currentRoots = []string{segments[0] + string(filepath.Separator)}

	for i, segment := range segments {
		if i == 0 {
			continue //skip the first
		}

		var newRoots = make([]string, 0)
		for _, root := range currentRoots {
			if strings.ContainsRune(segment, '*') {
				matches, _ := f.getDirectoriesFromRoot(root, segment)
				newRoots = append(newRoots, matches...)
			} else {
				var newRoot = filepath.Join(root, segment)
				var dirExists, _ = utils.File.FileOrDirectoryExists(newRoot)
				if dirExists {
					newRoots = append(newRoots, newRoot)
				}
			}
		}

		if len(newRoots) == 0 {
			return nil
		}

		currentRoots = newRoots
	}

	return currentRoots
}

type dirInfo struct {
	name string
	date time.Time
}

func (f *diffFinder) getDirectoriesFromRoot(root, segment string) ([]string, error) {
	matches := make([]dirInfo, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if matched, _ := filepath.Match(segment, info.Name()); matched {
				matches = append(matches, dirInfo{name: path, date: info.ModTime()})
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].date.Before(matches[j].date)
	})

	sorted := make([]string, 0)
	for _, match := range matches {
		sorted = append(sorted, match.name)
	}
	return sorted, nil
}
