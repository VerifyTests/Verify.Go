package diff

import (
	"github.com/VerifyTests/Verify.Go/utils"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFinderMultiMatchDir(t *testing.T) {
	createTestDir(t, "../_testdata/DirForSearch/", time.Now().Add(-3*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir2", time.Now().Add(-2*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir1", time.Now())

	var path = filepath.Join("../_testdata/DirForSearch", "*", "TextFile1.txt")
	var finder = newDiffFinder()

	result, found := finder.TryFind(path)

	if !found || !utils.File.Exists(result) {
		t.Fatalf("should find the file at path: %s", path)
	}
}

func TestFinderMultiMatchDirReverseOrder(t *testing.T) {
	createTestDir(t, "../_testdata/DirForSearch/", time.Now().Add(-3*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir1", time.Now().Add(-2*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir2", time.Now())

	var path = filepath.Join("../_testdata/DirForSearch", "*", "TextFile1.txt")

	var finder = newDiffFinder()

	result, found := finder.TryFind(path)

	if !found || !utils.File.Exists(result) {
		t.Fatalf("should find the file at path: %s", path)
	}
}

func TestFinderFindFullPath(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/DirForSearch/dir2/TextFile2.txt")
	var finder = newDiffFinder()

	result, found := finder.TryFind(path)

	if !found || !utils.File.Exists(result) {
		t.Fatalf("should find the file at path: %s", path)
	}
}

func TestFinderNonExistingFile(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/DirForSearch/dir2/TextFile2.bin")
	var finder = newDiffFinder()

	_, found := finder.TryFind(path)

	if found {
		t.Fatalf("should not find the non-existing file at path: %s", path)
	}
}

func TestFinderWildcardInDirectory(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/*/dir1/TextFile1.txt")
	var finder = newDiffFinder()

	result, found := finder.TryFind(path)

	if !found || !utils.File.Exists(result) {
		t.Fatalf("should find the file at path: %s", path)
	}
}

func TestFinderWildcardMissing(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/*/dir3/TextFile1.txt")
	var finder = newDiffFinder()

	_, found := finder.TryFind(path)

	if found {
		t.Fatalf("should not find the non-existing file at path: %s", path)
	}
}

func createTestDir(t *testing.T, dir string, time time.Time) {
	if _, err := os.Stat(dir); err == nil {
		_ = os.Chtimes(dir, time, time)
	} else {
		mke := os.Mkdir(dir, os.ModeSticky|os.ModePerm)
		if mke != nil {
			t.Fatalf("should not have errors creating directories: %s", mke)
		}

		che := os.Chtimes(dir, time, time)
		if mke != nil {
			t.Fatalf("should not have errors accessing directory times: %s", che)
		}
	}
}
