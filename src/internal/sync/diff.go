package sync

import (
	"excellgene.com/symbaSync/internal/infra/fs"
)

// Action represents what needs to be done for a file.
type Action int

const (
	ActionNone   Action = iota // File is in sync
	ActionCreate               // File needs to be created at destination
	ActionUpdate               // File exists but needs to be updated
	ActionDelete               // File exists at dest but not source (if delete enabled)
)

// FileDiff represents the difference for a single file.
type FileDiff struct {
	Path   string      // Relative path
	Action Action      // What needs to be done
	Source *fs.FileInfo // Source file info (nil if file only exists at dest)
	Dest   *fs.FileInfo // Dest file info (nil if file only exists at source)
}

// DiffResult contains all differences found between source and destination.
type DiffResult struct {
	Diffs []FileDiff
}

// Differ compares source and destination filesystems.
// Responsibility: Pure comparison logic. No I/O.
type Differ struct {
	// DeleteExtraFiles determines if files at dest but not source should be deleted
	DeleteExtraFiles bool
}

// NewDiffer creates a new differ with default settings.
func NewDiffer() *Differ {
	return &Differ{
		DeleteExtraFiles: false, // Safe default: don't delete
	}
}

// Diff compares source and destination file lists.
// Returns list of actions needed to sync dest to match source.
func (d *Differ) Diff(source, dest []fs.FileInfo) *DiffResult {
	sourceMap := make(map[string]fs.FileInfo)
	destMap := make(map[string]fs.FileInfo)

	// Build lookup maps
	for _, f := range source {
		sourceMap[f.Path] = f
	}
	for _, f := range dest {
		destMap[f.Path] = f
	}

	var diffs []FileDiff

	// Check each source file
	for path, srcFile := range sourceMap {
		destFile, existsAtDest := destMap[path]

		if !existsAtDest {
			// File exists at source but not dest -> create
			diffs = append(diffs, FileDiff{
				Path:   path,
				Action: ActionCreate,
				Source: &srcFile,
				Dest:   nil,
			})
		} else if d.needsUpdate(srcFile, destFile) {
			// File exists at both but needs update
			diffs = append(diffs, FileDiff{
				Path:   path,
				Action: ActionUpdate,
				Source: &srcFile,
				Dest:   &destFile,
			})
		}
	}

	// Check for files at dest but not source
	if d.DeleteExtraFiles {
		for path, destFile := range destMap {
			if _, existsAtSource := sourceMap[path]; !existsAtSource {
				diffs = append(diffs, FileDiff{
					Path:   path,
					Action: ActionDelete,
					Source: nil,
					Dest:   &destFile,
				})
			}
		}
	}

	return &DiffResult{Diffs: diffs}
}

// needsUpdate determines if a file needs to be updated.
// Currently uses modification time and size.
func (d *Differ) needsUpdate(source, dest fs.FileInfo) bool {
	// Skip directories
	if source.IsDir {
		return false
	}

	// If size differs, needs update
	if source.Size != dest.Size {
		return true
	}

	// If source is newer, needs update
	if source.ModTime > dest.ModTime {
		return true
	}

	return false
}
