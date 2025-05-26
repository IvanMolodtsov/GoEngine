package sdl

// #cgo LDFLAGS: -lSDL3
// #include <SDL3/SDL.h>
import "C"

type Window = C.SDL_Window

type WindowFlags = C.SDL_WindowFlags

type WindowID = C.SDL_WindowID

const (
	WindowFullscreen        WindowFlags = 0x0000000000000001
	WindowOpenGL            WindowFlags = 0x0000000000000002
	WindowOccluded          WindowFlags = 0x0000000000000004
	WindowHidden            WindowFlags = 0x0000000000000008
	WindowBorderless        WindowFlags = 0x0000000000000010
	WindowResizable         WindowFlags = 0x0000000000000020
	WindowMinimized         WindowFlags = 0x0000000000000040
	WindowMaximized         WindowFlags = 0x0000000000000080
	WindowMouseGrabbed      WindowFlags = 0x0000000000000100
	WindowInputFocus        WindowFlags = 0x0000000000000200
	WindowMouseFocus        WindowFlags = 0x0000000000000400
	WindowExternal          WindowFlags = 0x0000000000000800
	WindowModal             WindowFlags = 0x0000000000001000
	WindowHighPixelDensity  WindowFlags = 0x0000000000002000
	WindowMouseCapture      WindowFlags = 0x0000000000004000
	WindowMouseRelativeMode WindowFlags = 0x0000000000008000
	WindowAlwaysOnTop       WindowFlags = 0x0000000000010000
	WindowUtility           WindowFlags = 0x0000000000020000
	WindowTooltip           WindowFlags = 0x0000000000040000
	WindowPopupMenu         WindowFlags = 0x0000000000080000
	WindowKeyboardGrabbed   WindowFlags = 0x0000000000100000
	WindowVulkan            WindowFlags = 0x0000000010000000
	WindowMetal             WindowFlags = 0x0000000020000000
	WindowTransparent       WindowFlags = 0x0000000040000000
	WindowNotFocusable      WindowFlags = 0x0000000080000000
)

func CreateWindow(name string, width, height int32, flags WindowFlags) (*Window, error) {
	var w *Window = C.SDL_CreateWindow(C.CString(name), C.int(width), C.int(height), flags)
	if w == nil {
		return nil, GetError()
	}
	return w, nil
}

func DestroyWindow(window *Window) {
	C.SDL_DestroyWindow(window)
}
