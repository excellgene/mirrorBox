package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Copier interface {
	Copy(srcPath, dstPath string) error
}

type LocalCopier struct {
	preservePerms bool
	bufferSize    int
}

func NewLocalCopier(preservePerms bool) *LocalCopier {
	return &LocalCopier{
		preservePerms: preservePerms,
		bufferSize:    32 * 1024, // 32KB buffer
	}
}

func (c *LocalCopier) Copy(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}

	if srcInfo.IsDir() {
		return os.MkdirAll(dstPath, srcInfo.Mode())
	}

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("create parent directories: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy file contents: %w", err)
	}

	if c.preservePerms {
		if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
			return fmt.Errorf("set file permissions: %w", err)
		}

		if err := os.Chtimes(dstPath, srcInfo.ModTime(), srcInfo.ModTime()); err != nil {
			return fmt.Errorf("set file times: %w", err)
		}
	}

	return nil
}

func ExecuteCopy(c Copier, srcPath, dstPath string) error {
	if exists, err := Exists(srcPath); err != nil {
		return fmt.Errorf("check source exists: %w", err)
	} else if !exists {
		return fmt.Errorf("source path does not exist: %s", srcPath)
	}

	if err := c.Copy(srcPath, dstPath); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	return nil
}
