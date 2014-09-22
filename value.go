package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"


type Value struct {
	inC *C.Value
}

type ValueID struct {
	inC *C.ValueID
}

