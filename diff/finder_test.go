package diff

import (
	"github.com/stretchr/testify/assert"
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
	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.True(t, found)
	assert.FileExists(t, result)
}

func TestFinderMultiMatchDirReverseOrder(t *testing.T) {
	createTestDir(t, "../_testdata/DirForSearch/", time.Now().Add(-3*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir1", time.Now().Add(-2*24*time.Hour))
	createTestDir(t, "../_testdata/DirForSearch/dir2", time.Now())

	var path = filepath.Join("../_testdata/DirForSearch", "*", "TextFile1.txt")

	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.True(t, found)
	assert.FileExists(t, result)
}

func TestFinderFindFullPath(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/DirForSearch/dir2/TextFile2.txt")
	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.True(t, found)
	assert.FileExists(t, result)
}

func TestFinderNonExistingFile(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/DirForSearch/dir2/TextFile2.bin")
	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.False(t, found)
	assert.Empty(t, result)
}

func TestFinderWildcardInDirectory(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/*/dir1/TextFile1.txt")
	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.True(t, found)
	assert.FileExists(t, result)
}

func TestFinderWildcardMissing(t *testing.T) {
	var path, _ = filepath.Abs("../_testdata/*/dir3/TextFile1.txt")
	var finder = newFinder()

	result, found := finder.TryFind(path)

	assert.False(t, found)
	assert.Empty(t, result)
}

func createTestDir(t *testing.T, dir string, time time.Time) {
	if _, err := os.Stat(dir); err == nil {
		_ = os.Chtimes(dir, time, time)
	} else {
		mke := os.Mkdir(dir, 0700)
		assert.NoError(t, mke)

		che := os.Chtimes(dir, time, time)
		assert.NoError(t, che)
	}
}