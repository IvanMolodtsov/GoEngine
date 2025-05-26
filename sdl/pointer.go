package sdl

import "unsafe"

func CastPointer[T, P any](pointer *P) *T {
	return (*T)(unsafe.Pointer(pointer))
}

func DereferencePointer[T, P any](pointer *P) T {
	return *(*T)(unsafe.Pointer(pointer))
}
