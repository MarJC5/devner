//go:build !darwin

package sysenv

// EnsureDevPath is a no-op on non-darwin platforms. The launchd PATH
// stripping behaviour is macOS-specific; on Linux and Windows the
// inherited environment already contains the directories we need.
func EnsureDevPath() {}
