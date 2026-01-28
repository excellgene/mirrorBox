//go:build darwin || cgo

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
setActivationPolicy() {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    return 0;
}
*/
import "C"

// hides the application icon from the macOS dock.
// This makes the app behave as a "menu bar only" application.
func SetActivationPolicy() {
	C.setActivationPolicy()
}