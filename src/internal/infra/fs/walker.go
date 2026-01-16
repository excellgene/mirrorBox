package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path    string
	Size    int64
	ModTime int64
	IsDir   bool
}

type Walker interface {
	// Walk traverses the directory tree and calls fn for each file/directory.
	// Paths returned are relative to the root being walked.
	Walk(fn func(FileInfo) error) error
}

type LocalWalker struct {
	root string
}

func NewLocalWalker(root string) *LocalWalker {
	return &LocalWalker{root: root}
}

func (w *LocalWalker) Walk(fn func(FileInfo) error) error {
	return filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}

		// Get relative path from root
		relPath, err := filepath.Rel(w.root, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		// Skip root directory itself
		if relPath == "." {
			return nil
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("get file info for %s: %w", path, err)
		}

		fileInfo := FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
			IsDir:   info.IsDir(),
		}

		return fn(fileInfo)
	})
}

// Exists checks if a path exists on local filesystem.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
