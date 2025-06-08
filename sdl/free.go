package sdl

// #include "stdlib.h"
import "C"
import "unsafe"

func Free[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}
