package dnfjson

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gobwas/glob"
	"github.com/osbuild/osbuild-composer/internal/rpmmd"
)

// A collection of directory paths and their total size.
type pathGroup struct {
	paths []string
	size  uint64
	mtime time.Time
}

func shrinkCache(path string, maxSize uint64) error {
	tg, curSize, err := collectTimeGroups(path)
	if err != nil {
		return err
	}

	// start deleting until we drop below maxSize
	for idx := 0; idx < len(tg) && curSize >= maxSize; idx++ {
		for _, gPath := range tg[idx].paths {
			os.RemoveAll(gPath)
		}
		curSize -= tg[idx].size
	}

	return nil
}

func entrySize(path string) (uint64, error) {
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

// Collect paths in groups based on their modification time.  Each element
// includes the common modification time, a list of paths that share that
// modification time, and the total size of the collection.  The returned list
// is sorted by ascending modification time (oldest first).
// The total size of all groups is also returned.
func collectTimeGroups(path string) ([]pathGroup, uint64, error) {
	timeGroups := make(map[time.Time]pathGroup, 0)
	cacheEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, 0, err
	}

	var totalSize uint64

	// collect the paths grouped by their mtime
	for _, entry := range cacheEntries {
		eInfo, err := entry.Info()
		if err != nil {
			// skip it
			continue
		}

		mtime := eInfo.ModTime()
		pg := timeGroups[mtime]
		// add mod time (if not already set)
		ePath := filepath.Join(path, entry.Name())

		// calculate and add entry size
		size, err := entrySize(ePath)
		if err != nil {
			// skip it
			continue
		}
		pg.size += size
		totalSize += size

		// add path
		pg.paths = append(pg.paths, ePath)

		// add time
		pg.mtime = mtime

		// update the collection
		timeGroups[mtime] = pg
	}

	// create a list of groups sorted by mtime (oldest first)
	sortedGroups := make([]pathGroup, 0, len(timeGroups))
	searcher := func(idx int) bool {
		if idx == len(timeGroups)-1 {
			// last element: return
			return false
		}
		cur := sortedGroups[idx]
		next := sortedGroups[idx+1]
		if cur.mtime.After(next.mtime) {
			return true
		}
		return false
	}

	for _, pg := range timeGroups {
		insIdx := sort.Search(len(sortedGroups), searcher)
		sortedGroups = append(append(sortedGroups[:insIdx], pg), sortedGroups[insIdx+1:]...)
	}

	return sortedGroups, totalSize, nil
}

// Update file atime and mtime to "now" for all files in the root of each
// repository cache that match the repo IDs.  This should be called whenever a
// set of repositories is used.
func touchRepoCaches(cacheRoot string, repos map[string]rpmmd.RepoConfig) error {
	for _, r := range repos {
		if err := touchRepoCache(cacheRoot, r); err != nil {
			return err
		}
	}
	return nil
}

// Update file atime and mtime to "now" for all files in the root of the cache
// that match the repo ID.  This should be called whenever a repository is
// used.
func touchRepoCache(cacheRoot string, repo rpmmd.RepoConfig) error {
	repoGlob, err := glob.Compile(fmt.Sprintf("%s*", repo.Hash()))
	if err != nil {
		return err
	}

	// we only touch the top-level directories and files of the cache
	cacheEntries, err := os.ReadDir(cacheRoot)
	if err != nil {
		return err
	}

	// use the same timestamp for all entries
	now := time.Now().Local()
	for _, cacheEntry := range cacheEntries {
		if repoGlob.Match(cacheEntry.Name()) {
			path := filepath.Join(cacheRoot, cacheEntry.Name())
			if err := os.Chtimes(path, now, now); err != nil {
				return err
			}
		}
	}
	return nil
}
