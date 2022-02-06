package data

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/palantir/stacktrace"
)

// VisitFn is a callback that gets a chance to visit each file found during a
// traversal with visit().
type ValidatorFn func(path string) error

// Validate will call a validator function on the path if it's a file, or for each
// file in the path if it's a directory. this is used for running
// custom validation logic on configuration artifacts
func Validate(path string, validator ValidatorFn) error {
	f, err := os.Open(path)
	if err != nil {
		return stacktrace.Propagate(err, "cannot open file descriptor %s: %v", path)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return stacktrace.Propagate(err, "cannot stat %s", path)
	}

	if !fi.IsDir() {
		err = validator(path)
		if err != nil {
			return stacktrace.Propagate(err, "cannot validate %s", path)
		}
		return nil
	}

	contents, err := f.Readdir(-1)
	if err != nil {
		return stacktrace.Propagate(err, "cannot list %s: %v", path)
	}

	sort.Sort(dirEnts(contents))
	for _, fi := range contents {
		if fi.IsDir() {
			continue
		}

		fullPath := filepath.Join(path, fi.Name())
		err = validator(fullPath)
		if err != nil {
			return stacktrace.Propagate(err, "cannot validate %s", fullPath)
		}
	}

	return nil
}

// dirEnts applies sort.Interface to directory entries for sorting by name.
type dirEnts []os.FileInfo

func (d dirEnts) Len() int           { return len(d) }
func (d dirEnts) Less(i, j int) bool { return d[i].Name() < d[j].Name() }
func (d dirEnts) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
