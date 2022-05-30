package dnfjson

import (
	"io/fs"
	"os"
	"path/filepath"
)

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

func shrinkDir(path string, maxSize uint64) error {
	curSize, err := dirSize(path)
	if err != nil {
		return err
	}
	if curSize > maxSize {
		return os.RemoveAll(path)
	}
	return nil
}
