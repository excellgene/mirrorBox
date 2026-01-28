# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MirroxBox is a cross-platform folder synchronization tool built with Go and Fyne (GUI framework). It runs as a system tray application that periodically syncs configured folder pairs from source to destination.

## Build & Run Commands

### Development

Run the application (from repository root):
```bash
# macOS
go run ./src/cmd/app/main.go

# Windows
go run .\src\cmd\app\main.go
```

Build the application (working directory must be `src/`):
```bash
cd src
go build ./cmd/app
```

### Production Build

Package for distribution (working directory must be `src/cmd/app/`):
```bash
cd src/cmd/app

# macOS
fyne package -name mirrorBox -os darwin -icon ../../internal/ui/icons/appicon.png --metadata LSUIElement=true

# Windows
fyne package -name mirrorBox -os windows -icon ../../internal/ui/icons/appicon.png
```

**Note**: The dock icon is automatically hidden on macOS through `platform.HideDockIcon()` which uses CGO to call Objective-C's `NSApplicationActivationPolicyAccessory`. This makes the app appear only in the system tray, not in the dock.

## Architecture

### Core Components

The application follows a layered architecture with clear separation of concerns:

**Entry Point** (`src/cmd/app/main.go`):
- Initializes config store from `~/.config/mirrorbox/config.json`
- Creates app state, job factory, and dispatcher
- Sets up system tray, status window, and settings window
- Coordinates communication via two event loops:
  - `handleTrayEvents`: Processes user actions (sync now, settings, status, quit)
  - `handleDispatcherEvents`: Updates UI when jobs complete

**Configuration** (`internal/config/`):
- `Config`: Holds check interval, start-at-boot flag, and folder pairs
- `Store`: Persists config as JSON to disk
- Each folder pair has SourcePath, DestinationPath, and Enabled flag

**Application State** (`internal/app/`):
- `State`: Thread-safe registry of all sync jobs (map[string]*Job)
- `Dispatcher`: Orchestrates job execution
  - Runs jobs on-demand or via scheduler
  - Emits `JobEvent` on completion (includes status, result, error)
  - Uses context for cancellation and 30-minute timeouts per job
- `JobFactory`: Creates Job instances from Config

**Sync Engine** (`internal/sync/`):
- `Job`: Represents a sync task (name, source, dest, status, last run)
  - Execution flow: Walk source → Walk dest → Compute diff → Apply sync
- `Differ`: Compares file lists and produces list of actions (create/update/delete)
  - Uses mod time + size to detect changes
  - `DeleteExtraFiles` flag controls whether files only at dest are removed
- `Syncer`: Applies diff operations to make dest match source
- `SyncResult`: Statistics about files created/updated/deleted and bytes copied

**Filesystem Abstraction** (`internal/sync/fs/`):
- `Walker`: Interface for traversing directories (currently `LocalWalker` only)
- `Copier`: Interface for copying files (currently `LocalCopier` only)
- This abstraction allows future support for SMB/network destinations

**UI** (`internal/ui/`):
- `StatusWindow`: Shows job execution status and history
- `SettingsWindow`: Configure folders, interval, and start-at-boot

**System Tray** (`internal/tray/`):
- `Tray`: Fyne-based system tray integration
- `Menu`: System tray menu (Sync Now, Settings, Status, Quit)
- Emits `Event` enum for user actions
- Icon adapts to system theme (dark/light)

**Platform-Specific** (`internal/platform/`):
- `dock_darwin.go`: macOS-specific code using CGO/Objective-C to hide dock icon
- `dock_other.go`: Stub implementation for non-macOS platforms
- Uses `NSApplicationActivationPolicyAccessory` to make app tray-only

### Key Design Patterns

**Event-Driven Communication**: Components communicate via channels:
- Tray emits user action events
- Dispatcher emits job completion events
- Main goroutines bridge these events to UI updates via `fyne.Do()`

**Concurrency Model**:
- Each sync job runs in its own goroutine
- Dispatcher uses WaitGroup to track active jobs
- Context with timeout prevents jobs from hanging indefinitely

**Dependency Injection**: Jobs receive walker/differ/syncer dependencies, making them testable and allowing different implementations (local vs network filesystems)

## Configuration

Config location: `~/.config/mirrorbox/config.json`

Example structure:
```json
{
  "check_interval": 300000000000,
  "start_at_boot": false,
  "folders": [
    {
      "SourcePath": "/path/to/source",
      "DestinationPath": "/path/to/dest",
      "Enabled": true
    }
  ]
}
```

## Go Module

Module name: `excellgene.com/symbaSync`
Go version: 1.22
Primary dependency: Fyne v2 (`fyne.io/fyne/v2`)
