package openzwave

//
// #cgo LDFLAGS: -lopenzwave -L../go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include <stdlib.h>
// #include "wrapper.hpp"
import "C"

import "unsafe"

type API struct {
	available bool // true if a zwave device is available.
}

func NewAPI(device string) *API {
	var cDevice *C.char = C.CString(device);
	defer C.free(unsafe.Pointer(cDevice));
	var rc C.int = C.init(cDevice);
	return &API{rc != 0}
}
