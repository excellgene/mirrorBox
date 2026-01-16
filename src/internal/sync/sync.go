package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"excellgene.com/symbaSync/internal/infra/smb"
)

type SyncResult struct {
	FilesCreated int
	FilesUpdated int
	FilesDeleted int
	BytesCopied  int64
	Errors       []error
}

type Syncer struct {
	smbClient smb.Client
}

// NewSyncer creates a new syncer with the given SMB client.
func NewSyncer(smbClient smb.Client) *Syncer {
	return &Syncer{
		smbClient: smbClient,
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

	localPath := filepath.Join(sourcePath, diff.Path)
	remotePath := filepath.Join(destPath, diff.Path)

	// If it's a directory, just create it
	if diff.Source.IsDir {
		return s.smbClient.MkdirAll(ctx, remotePath)
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(remotePath)
	if err := s.smbClient.MkdirAll(ctx, parentDir); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	// Open local file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer file.Close()

	// Upload to SMB
	if err := s.smbClient.Upload(ctx, remotePath, file, diff.Source.Size); err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	return nil
}

// update handles updating an existing file at destination.
func (s *Syncer) update(ctx context.Context, diff FileDiff, sourcePath, destPath string) error {
	// For now, update is same as create (overwrite)
	return s.create(ctx, diff, sourcePath, destPath)
}

// delete handles removing a file or directory from destination.
func (s *Syncer) delete(ctx context.Context, diff FileDiff, destPath string) error {
	remotePath := filepath.Join(destPath, diff.Path)
	return s.smbClient.Delete(ctx, remotePath)
}
