package blueprint

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/osbuild/osbuild-composer/internal/common"
	"github.com/osbuild/osbuild-composer/internal/fsnode"
)

// validateModeString checks that the given string is a valid mode octal number
func validateModeString(mode string) error {
	// check that the mode string matches the octal format regular expression
	if regexp.MustCompile(`^0[0-7]{3}$`).MatchString(mode) {
		return nil
	}
	return fmt.Errorf("invalid mode %s: must be an octal number", mode)
}

// DirectoryCustomization represents a directory to be created in the image
type DirectoryCustomization struct {
	// Absolute path to the directory
	Path string `json:"path" toml:"path"`
	// Owner of the directory specified as a string (user name), int64 (UID) or nil
	User interface{} `json:"user,omitempty" toml:"user,omitempty"`
	// Owner of the directory specified as a string (group name), int64 (UID) or nil
	Group interface{} `json:"group,omitempty" toml:"group,omitempty"`
	// Permissions of the directory specified as an octal number
	Mode string `json:"mode,omitempty" toml:"mode,omitempty"`
	// EnsureParents ensures that all parent directories of the directory exist
	EnsureParents bool `json:"ensure_parents,omitempty" toml:"ensure_parents,omitempty"`
}

// Custom TOML unmarshalling for DirectoryCustomization with validation
func (d *DirectoryCustomization) UnmarshalTOML(data interface{}) error {
	var dir DirectoryCustomization

	dataMap, _ := data.(map[string]interface{})

	switch path := dataMap["path"].(type) {
	case string:
		dir.Path = path
	default:
		return fmt.Errorf("UnmarshalTOML: path must be a string")
	}

	switch user := dataMap["user"].(type) {
	case string:
		dir.User = user
	case int64:
		dir.User = user
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: user must be a string or an integer, got %T", user)
	}

	switch group := dataMap["group"].(type) {
	case string:
		dir.Group = group
	case int64:
		dir.Group = group
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: group must be a string or an integer")
	}

	switch mode := dataMap["mode"].(type) {
	case string:
		dir.Mode = mode
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: mode must be a string")
	}

	switch ensure_parents := dataMap["ensure_parents"].(type) {
	case bool:
		dir.EnsureParents = ensure_parents
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: ensure_parents must be a bool")
	}

	// try converting to fsnode.Directory to validate all values
	_, err := dir.ToFsNodeDirectory()
	if err != nil {
		return err
	}

	*d = dir
	return nil
}

// Custom JSON unmarshalling for DirectoryCustomization with validation
func (d *DirectoryCustomization) UnmarshalJSON(data []byte) error {
	type directoryCustomization DirectoryCustomization

	var dirPrivate directoryCustomization
	if err := json.Unmarshal(data, &dirPrivate); err != nil {
		return err
	}

	dir := DirectoryCustomization(dirPrivate)
	if uid, ok := dir.User.(float64); ok {
		// check if uid can be converted to int64
		if uid != float64(int64(uid)) {
			return fmt.Errorf("invalid user %f: must be an integer", uid)
		}
		dir.User = int64(uid)
	}
	if gid, ok := dir.Group.(float64); ok {
		// check if gid can be converted to int64
		if gid != float64(int64(gid)) {
			return fmt.Errorf("invalid group %f: must be an integer", gid)
		}
		dir.Group = int64(gid)
	}
	// try converting to fsnode.Directory to validate all values
	_, err := dir.ToFsNodeDirectory()
	if err != nil {
		return err
	}

	*d = dir
	return nil
}

// ToFsNodeDirectory converts the DirectoryCustomization to an fsnode.Directory
func (d DirectoryCustomization) ToFsNodeDirectory() (*fsnode.Directory, error) {
	var mode *os.FileMode
	if d.Mode != "" {
		err := validateModeString(d.Mode)
		if err != nil {
			return nil, err
		}
		modeNum, err := strconv.ParseUint(d.Mode, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid mode %s: %v", d.Mode, err)
		}
		mode = common.ToPtr(os.FileMode(modeNum))
	}

	return fsnode.NewDirectory(d.Path, mode, d.User, d.Group, d.EnsureParents)
}

// FileCustomization represents a file to be created in the image
type FileCustomization struct {
	// Absolute path to the file
	Path string `json:"path" toml:"path"`
	// Owner of the directory specified as a string (user name), int64 (UID) or nil
	User interface{} `json:"user,omitempty" toml:"user,omitempty"`
	// Owner of the directory specified as a string (group name), int64 (UID) or nil
	Group interface{} `json:"group,omitempty" toml:"group,omitempty"`
	// Permissions of the file specified as an octal number
	Mode string `json:"mode,omitempty" toml:"mode,omitempty"`
	// Data is the file content in plain text
	Data string `json:"data,omitempty" toml:"data,omitempty"`
}

// Custom TOML unmarshalling for FileCustomization with validation
func (f *FileCustomization) UnmarshalTOML(data interface{}) error {
	var file FileCustomization

	dataMap, _ := data.(map[string]interface{})

	switch path := dataMap["path"].(type) {
	case string:
		file.Path = path
	default:
		return fmt.Errorf("UnmarshalTOML: path must be a string")
	}

	switch user := dataMap["user"].(type) {
	case string:
		file.User = user
	case int64:
		file.User = user
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: user must be a string or an integer")
	}

	switch group := dataMap["group"].(type) {
	case string:
		file.Group = group
	case int64:
		file.Group = group
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: group must be a string or an integer")
	}

	switch mode := dataMap["mode"].(type) {
	case string:
		file.Mode = mode
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: mode must be a string")
	}

	switch data := dataMap["data"].(type) {
	case string:
		file.Data = data
	case nil:
		break
	default:
		return fmt.Errorf("UnmarshalTOML: data must be a string")
	}

	// try converting to fsnode.File to validate all values
	_, err := file.ToFsNodeFile()
	if err != nil {
		return err
	}

	*f = file
	return nil
}

// Custom JSON unmarshalling for FileCustomization with validation
func (f *FileCustomization) UnmarshalJSON(data []byte) error {
	type fileCustomization FileCustomization

	var filePrivate fileCustomization
	if err := json.Unmarshal(data, &filePrivate); err != nil {
		return err
	}

	file := FileCustomization(filePrivate)
	if uid, ok := file.User.(float64); ok {
		// check if uid can be converted to int64
		if uid != float64(int64(uid)) {
			return fmt.Errorf("invalid user %f: must be an integer", uid)
		}
		file.User = int64(uid)
	}
	if gid, ok := file.Group.(float64); ok {
		// check if gid can be converted to int64
		if gid != float64(int64(gid)) {
			return fmt.Errorf("invalid group %f: must be an integer", gid)
		}
		file.Group = int64(gid)
	}
	// try converting to fsnode.File to validate all values
	_, err := file.ToFsNodeFile()
	if err != nil {
		return err
	}

	*f = file
	return nil
}

// ToFsNodeFile converts the FileCustomization to an fsnode.File
func (f FileCustomization) ToFsNodeFile() (*fsnode.File, error) {
	var data []byte
	if f.Data != "" {
		data = []byte(f.Data)
	}

	var mode *os.FileMode
	if f.Mode != "" {
		err := validateModeString(f.Mode)
		if err != nil {
			return nil, err
		}
		modeNum, err := strconv.ParseUint(f.Mode, 8, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid mode %s: %v", f.Mode, err)
		}
		mode = common.ToPtr(os.FileMode(modeNum))
	}

	return fsnode.NewFile(f.Path, data, mode, f.User, f.Group)
}
