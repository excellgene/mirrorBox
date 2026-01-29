package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"excellgene.com/mirrorBox/internal/sync/fs"
)

type SyncResult struct {
	FilesCreated int
	FilesUpdated int
	FilesDeleted int
	BytesCopied  int64
	Errors       []error
}

type Syncer struct {
	copier fs.Copier
}

// NewSyncer creates a new syncer with a file copier.
func NewSyncer(copier fs.Copier) *Syncer {
	return &Syncer{
		copier: copier,
	}
}

// Sync applies the diff to make destination match source.
// ctx allows cancellation of long-running operations.
func (s *Syncer) Sync(ctx context.Context, diff *DiffResult, sourcePath, destPath string) (*SyncResult, error) {
	result := &SyncResult{}

	for _, fileDiff := range diff.Diffs {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		var err error
		switch fileDiff.Action {
		case ActionCreate:
			err = s.create(ctx, fileDiff, sourcePath, destPath)
			if err == nil {
				result.FilesCreated++
				if fileDiff.Source != nil {
					result.BytesCopied += fileDiff.Source.Size
				}
			}

		case ActionUpdate:
			err = s.update(ctx, fileDiff, sourcePath, destPath)
			if err == nil {
				result.FilesUpdated++
				if fileDiff.Source != nil {
					result.BytesCopied += fileDiff.Source.Size
				}
			}

		case ActionDelete:
			err = s.delete(ctx, fileDiff, destPath)
			if err == nil {
				result.FilesDeleted++
			}
		}

		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", fileDiff.Path, err))
		}
	}

	return result, nil
}

// create handles creating a new file or directory at destination.
func (s *Syncer) create(ctx context.Context, diff FileDiff, sourcePath, destPath string) error {
	if diff.Source == nil {
		return fmt.Errorf("no source file info")
	}

	srcPath := filepath.Join(sourcePath, diff.Path)
	dstPath := filepath.Join(destPath, diff.Path)

	if diff.Source.IsDir {
		return os.MkdirAll(dstPath, 0755)
	}

	if err := s.copier.Copy(srcPath, dstPath); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}

// update handles updating an existing file at destination.
func (s *Syncer) update(ctx context.Context, diff FileDiff, sourcePath, destPath string) error {
	return s.create(ctx, diff, sourcePath, destPath)
}

func (s *Syncer) delete(ctx context.Context, diff FileDiff, destPath string) error {
	targetPath := filepath.Join(destPath, diff.Path)
	return os.RemoveAll(targetPath)
}
