//go:build !darwin || !cgo

package platform

// HideDockIcon is a no-op on non-macOS platforms.
// Platform-specific implementations will be added as needed.
func HideDockIcon() {
	// No-op on Windows and Linux
}