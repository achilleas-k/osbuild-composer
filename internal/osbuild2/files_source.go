package osbuild2

// Inputs for individual files

// Provides all the files, named via their content hash, specified
// via `references` in a new directory.
type FileSource struct {
	Type       string
	Origin     string
	References FileSourceChecksums
}

func (FileSource) isSource() {}

func NewFileSource(references FileSourceChecksums) *FileSource {
	return &FileSource{
		Type:   "org.osbuild.files",
		Origin: "org.osbuild.source",
	}
}

// Checksums of files to use as files input
type FileSourceChecksums []string
