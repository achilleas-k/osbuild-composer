package dnfjson

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gobwas/glob"
)

type rpmCache struct {
	// root path for the cache
	root string

	// individual repository cache data
	repoElements map[string]pathInfo

	// list of known repository IDs, sorted by mtime
	repoRecency []string

	// total cache size
	size uint64

	// max cache size
	maxSize uint64
}

func newRPMCache(path string, maxSize uint64) *rpmCache {
	r := &rpmCache{
		root:         path,
		repoElements: make(map[string]pathInfo),
		size:         0,
		maxSize:      maxSize,
	}
	// collect existing cache paths and timestamps
	r.updateInfo()
	return r
}

// updateInfo updates the repoPaths and repoRecency fields of the rpmCache.
func (r *rpmCache) updateInfo() {
	repos := make(map[string]pathInfo)
	repoIDs := make([]string, 0)
	cacheEntries, _ := os.ReadDir(r.root)

	var totalSize uint64

	// collect the paths grouped by their repo ID (first 64 characters of a file or directory name)
	for _, entry := range cacheEntries {
		eInfo, err := entry.Info()
		if err != nil {
			// skip it
			continue
		}

		fname := entry.Name()
		if len(fname) < 64 {
			// unknown file in cache; ignore
			continue
		}
		repoID := fname[:64]
		repo, ok := repos[repoID]
		if !ok {
			// new repo ID
			repoIDs = append(repoIDs, repoID)
		}
		mtime := eInfo.ModTime()
		ePath := filepath.Join(r.root, entry.Name())

		// calculate and add entry size
		size, err := dirSize(ePath)
		if err != nil {
			// skip it
			continue
		}
		repo.size += size
		totalSize += size

		// add path
		repo.paths = append(repo.paths, ePath)

		// if for some reason the mtimes of the various entries of a single
		// repository are out of sync, use the most recent one
		if repo.mtime.Before(mtime) {
			repo.mtime = mtime
		}

		// update the collection
		repos[repoID] = repo
	}
	sortFunc := func(idx, jdx int) bool {
		ir := repos[repoIDs[idx]]
		jr := repos[repoIDs[jdx]]
		return ir.mtime.Before(jr.mtime)
	}

	// sort IDs by mtime (oldest first)
	sort.Slice(repoIDs, sortFunc)

	r.size = totalSize
	r.repoElements = repos
	r.repoRecency = repoIDs
}

// create a file only if it doesn't exist.
func createExclusive(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
}

type cacheLock struct {
	path string
}

// Create an exclusive write lock (wlock) file in the rpmCache to prevent
// concurrent access to the global cache directory.  This is not required when
// depsolving since dnf takes care of locking individual repository cache
// directories.  The primary purpose of this lock is to prevent multiple
// concurrent cache deletions and to prevent reading or updating the caches
// while they are being deleted.
// An exclusive lock can only be acquired if no other locks, exclusive or not,
// exist.  If an exclusive lock is requested, no other lock can be acquired
// until it is released.
func (r *rpmCache) wlock() (*cacheLock, error) {
	ldir := r.lockdir()

	{ // create a hold file
		holdfile := filepath.Join(ldir, "exclusive.hold") // signals that an exclusive lock is waiting
		var fp *os.File
		var err error
		for fp, err = createExclusive(holdfile); errors.Is(err, os.ErrExist); fp, err = createExclusive(holdfile) {
			// will stop when err is nil or a different kind of error
		}
		if err != nil {
			// file does not exist, but a different kind of error occurred
			return nil, err
		}
		defer fp.Close()
		defer os.Remove(holdfile) // remove the hold if we get the lock (or fail)
	}

	lockfile := filepath.Join(ldir, "exclusive")
	{ // create the lock file
		var fp *os.File
		var err error
		for fp, err = createExclusive(lockfile); errors.Is(err, os.ErrExist); fp, err = createExclusive(lockfile) {
			// will stop when err is nil or a different kind of error
		}
		if err != nil {
			// file does not exist, but a different kind of error occurred
			return nil, err
		}
		defer fp.Close()
	}
	return &cacheLock{path: lockfile}, nil
}

// Remove the lock file.  Returns error if the file removal fails.  If the lock
// file does not exist, this is a no-op and returns nil.
func (cl *cacheLock) unlock() error {
	lockfile := cl.path
	if _, err := os.Stat(lockfile); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(lockfile)
}

func (r *rpmCache) shrink() error {
	lock, err := r.wlock()
	if err != nil {
		return err
	}
	defer func() {
		if err := lock.unlock(); err != nil {
			// the unlock function only returns with an error if the cache is
			// actually locked and it fails to remove the lockfile.
			// panic if we failed to unlock, otherwise we'd stay locked forever
			panic(fmt.Sprintf("failed to unlock cache after shrink: %s", err))
		}
	}()
	// start deleting until we drop below r.maxSize
	nDeleted := 0
	for idx := 0; idx < len(r.repoRecency) && r.size >= r.maxSize; idx++ {
		repoID := r.repoRecency[idx]
		nDeleted++
		repo, ok := r.repoElements[repoID]
		if !ok {
			// cache inconsistency?
			// ignore and let the ID be removed from the recency list
			continue
		}
		for _, gPath := range repo.paths {
			if err := os.RemoveAll(gPath); err != nil {
				return err
			}
		}
		r.size -= repo.size
		delete(r.repoElements, repoID)
	}

	// update recency list
	r.repoRecency = r.repoRecency[nDeleted:]
	return nil
}

// Update file atime and mtime on the filesystem to time t for all files in the
// root of the cache that match the repo ID.  This should be called whenever a
// repository is used.
// This function does not update the internal cache info.  A call to
// updateInfo() should be made after touching one or more repositories.
func (r *rpmCache) touchRepo(repoID string, t time.Time) error {
	repoGlob, err := glob.Compile(fmt.Sprintf("%s*", repoID))
	if err != nil {
		return err
	}

	// we only touch the top-level directories and files of the cache
	cacheEntries, err := os.ReadDir(r.root)
	if err != nil {
		return err
	}

	for _, cacheEntry := range cacheEntries {
		if repoGlob.Match(cacheEntry.Name()) {
			path := filepath.Join(r.root, cacheEntry.Name())
			if err := os.Chtimes(path, t, t); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create the cache lock directory and return the path.
func (r *rpmCache) lockdir() string {
	d := filepath.Join(r.root, ".composer.lock")
	if err := os.MkdirAll(d, 0770); err != nil {
		panic(fmt.Sprintf("cache lock directory creation failed: %s", err))
	}
	return d
}

// A collection of directory paths, their total size, and their most recent
// modification time.
type pathInfo struct {
	paths []string
	size  uint64
	mtime time.Time
}

func dirSize(path string) (uint64, error) {
	var size uint64
	sizer := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += uint64(info.Size())
		return nil
	}
	err := filepath.Walk(path, sizer)
	return size, err
}
