package sdl

// #cgo LDFLAGS: -L. -lSDL3
// #include <SDL3/SDL.h>
import "C"
import "unsafe"

type EventType = C.SDL_EventType

const (
	EventFirst                      EventType = 0x0
	EventQuit                       EventType = 0x100
	EventTerminating                EventType = 0x101
	EventLowMemory                  EventType = 0x102
	EventWillEnterBackground        EventType = 0x103
	EventDidEnterBackground         EventType = 0x104
	EventWillEnterForeground        EventType = 0x105
	EventDidEnterForeground         EventType = 0x106
	EventLocaleChanged              EventType = 0x107
	EventSystemThemeChanged         EventType = 0x108
	EventDisplayOrientation         EventType = 0x151
	EventDisplayFirst               EventType = 0x151
	EventDisplayAdded               EventType = 0x152
	EventDisplayRemoved             EventType = 0x153
	EventDisplayMoved               EventType = 0x154
	EventDisplayDesktopModeChanged  EventType = 0x155
	EventDisplayCurrentModeChanged  EventType = 0x156
	EventDisplayContentScaleChanged EventType = 0x157
	EventDisplayLast                EventType = 0x157
	EventWindowShown                EventType = 0x202
	EventWindowFirst                EventType = 0x202
	EventWindowHidden               EventType = 0x203
	EventWindowExposed              EventType = 0x204
	EventWindowMoved                EventType = 0x205
	EventWindowResized              EventType = 0x206
	EventWindowPixelSizeChanged     EventType = 0x207
	EventWindowMetalViewResized     EventType = 0x208
	EventWindowMinimized            EventType = 0x209
	EventWindowMaximized            EventType = 0x20A
	EventWindowRestored             EventType = 0x20B
	EventWindowMouseEnter           EventType = 0x20C
	EventWindowMouseLeave           EventType = 0x20D
	EventWindowFocusGained          EventType = 0x20E
	EventWindowFocusLost            EventType = 0x20F
	EventWindowCloseRequested       EventType = 0x210
	EventWindowHitTest              EventType = 0x211
	EventWindowICCProfChanged       EventType = 0x212
	EventWindowDisplayChanged       EventType = 0x213
	EventWindowDisplayScaleChanged  EventType = 0x214
	EventWindowSafeAreaChanged      EventType = 0x215
	EventWindowOccluded             EventType = 0x216
	EventWindowEnterFullscreen      EventType = 0x217
	EventWindowLeaveFullscreen      EventType = 0x218
	EventWindowDestroyed            EventType = 0x219
	EventWindowHDRStateChanged      EventType = 0x21A
	EventWindowLast                 EventType = 0x21A
	EventKeyDown                    EventType = 0x300
	EventKeyUp                      EventType = 0x301
	EventTextEditing                EventType = 0x302
	EventTextInput                  EventType = 0x303
	EventKeymapChanged              EventType = 0x304
	EventKeyboardAdded              EventType = 0x305
	EventKeyboardRemoved            EventType = 0x306
	EventTextEditingCandidates      EventType = 0x307
	EventMouseMotion                EventType = 0x400
	EventMouseButtonDown            EventType = 0x401
	EventMouseButtonUp              EventType = 0x402
	EventMouseWheel                 EventType = 0x403
	EventMouseAdded                 EventType = 0x404
	EventMouseRemoved               EventType = 0x405
	EventJoystickAxisMotion         EventType = 0x600
	EventJoystickBallMotion         EventType = 0x601
	EventJoystickHatMotion          EventType = 0x602
	EventJoystickButtonDown         EventType = 0x603
	EventJoystickButtonUp           EventType = 0x604
	EventJoystickAdded              EventType = 0x605
	EventJoystickRemoved            EventType = 0x606
	EventJoystickBatteryUpdated     EventType = 0x607
	EventJoystickUpdateComplete     EventType = 0x608
	EventGamepadAxisMotion          EventType = 0x650
	EventGamepadButtonDown          EventType = 0x651
	EventGamepadButtonUp            EventType = 0x652
	EventGamepadAdded               EventType = 0x653
	EventGamepadRemoved             EventType = 0x654
	EventGamepadRemapped            EventType = 0x655
	EventGamepadTouchpadDown        EventType = 0x656
	EventGamepadTouchpadMotion      EventType = 0x657
	EventGamepadTouchpadUp          EventType = 0x658
	EventGamepadSensorUpdate        EventType = 0x659
	EventGamepadUpdateComplete      EventType = 0x65A
	EventGamepadSteamHandleUpdated  EventType = 0x65B
	EventFingerDown                 EventType = 0x700
	EventFingerUp                   EventType = 0x701
	EventFingerMotion               EventType = 0x702
	EventFingerCanceled             EventType = 0x703
	EventClipboardUpdate            EventType = 0x900
	EventDropFile                   EventType = 0x1000
	EventDropText                   EventType = 0x1001
	EventDropBegin                  EventType = 0x1002
	EventDropComplete               EventType = 0x1003
	EventDropPosition               EventType = 0x1004
	EventAudioDeviceAdded           EventType = 0x1100
	EventAudioDeviceRemoved         EventType = 0x1101
	EventAudioDeviceFormatChanged   EventType = 0x1102
	EventSensorUpdate               EventType = 0x1200
	EventPenProximityIn             EventType = 0x1300
	EventPenProximityOut            EventType = 0x1301
	EventPenDown                    EventType = 0x1302
	EventPenUp                      EventType = 0x1303
	EventPenButtonDown              EventType = 0x1304
	EventPenButtonUp                EventType = 0x1305
	EventPenMotion                  EventType = 0x1306
	EventPenAxis                    EventType = 0x1307
	EventCameraDeviceAdded          EventType = 0x1400
	EventCameraDeviceRemoved        EventType = 0x1401
	EventCameraDeviceApproved       EventType = 0x1402
	EventCameraDeviceDenied         EventType = 0x1403
	EventRenderTargetsReset         EventType = 0x2000
	EventRenderDeviceReset          EventType = 0x2001
	EventRenderDeviceLost           EventType = 0x2002
	EventPrivate0                   EventType = 0x4000
	EventPrivate1                   EventType = 0x4001
	EventPrivate2                   EventType = 0x4002
	EventPrivate3                   EventType = 0x4003
	EventPollSentinel               EventType = 0x7F00
	EventUser                       EventType = 0x8000
	EventLast                       EventType = 0xFFFF
	EventEnumPadding                EventType = 0x7FFFFFFF
)

type Event [128]byte

func (e *Event) cPtr() *C.SDL_Event {
	return (*C.SDL_Event)(unsafe.Pointer(e))
}

func PollEvent(event *Event) bool {
	return bool(C.SDL_PollEvent(event.cPtr()))
}

// CommonEvent fields are shared by every event.
type CommonEvent struct {
	Type     EventType
	Reserved uint32
	// Timestamp in nanoseconds, populated using [GetTicksNS].
	Timestamp uint64
}

func (e *Event) Type() EventType {
	return *(*EventType)(unsafe.Pointer(e))
}

func (e *Event) Key() KeyboardEvent {
	return DereferencePointer[KeyboardEvent](e)
}

type KeyboardEvent struct {
	CommonEvent
	WindowID WindowID
	Which    KeyboardID
	Scancode Scancode
	Key      Keycode
	Mod      Keymod
	Raw      uint16
	Down     bool
	Repeat   bool
}

func (e *Event) MouseMotion() MouseMotionEvent {
	return DereferencePointer[MouseMotionEvent](e)
}

type MouseMotionEvent struct {
	CommonEvent
	WindowID WindowID
	Which    MouseID
	State    MouseButtonFlags
	X        float32
	Y        float32
	XRel     float32
	YRel     float32
}
