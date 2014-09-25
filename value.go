package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ninjasphere/go-openzwave/CC"
	"github.com/ninjasphere/go-openzwave/VT"
)

type Value interface {
}

type value struct {
	cRef *C.Value
}

func (self value) String() string {
	return fmt.Sprintf(
		"Value["+
			"type=%v, "+
			"commandClassId=%v, "+
			"instance=%d, "+
			"index=%d, "+
			"value='%s', "+
			"label='%s', "+
			"units='%s', "+
			"help='%s', "+
			"min=%d "+
			"max=%d "+
			"isSet=%v]",
		VT.ToEnum(int(self.cRef.valueId.valueType)),
		CC.ToEnum(int(self.cRef.valueId.commandClassId)),
		uint(self.cRef.valueId.instance),
		uint(self.cRef.valueId.index),
		C.GoString(self.cRef.value),
		C.GoString(self.cRef.label),
		C.GoString(self.cRef.units),
		C.GoString(self.cRef.help),
		self.cRef.min,
		self.cRef.max,
		self.cRef.isSet)
}

// convert a reference from the C Value to the Go Value
func asValue(cRef *C.Value) Value {
	return Value((*value)(cRef.goRef))
}

//export newGoValue
func newGoValue(cRef *C.Value) unsafe.Pointer {
	goRef := &value{cRef}
	cRef.goRef = unsafe.Pointer(goRef)
	return cRef.goRef
}

func (self *value) Notify(api API, notification Notification) {
	// TODO
}
