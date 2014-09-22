package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import (
	"fmt"

	"github.com/ninjasphere/go-openzwave/CC"
	"github.com/ninjasphere/go-openzwave/VT"
)

type ValueID struct {
	cRef *C.ValueID
}

func (self ValueID) String() string {
	return fmt.Sprintf(
		"ValueID["+
			"type=%s, "+
			"commandClassId=%s, "+
			"instance=%d, "+
			"index=%d]",
		VT.ToEnum(int(self.cRef.valueType)),
		CC.ToEnum(int(self.cRef.commandClassId)),
		uint(self.cRef.instance),
		uint(self.cRef.index))
}

type Value struct {
	cRef *C.Value
}

func (self Value) String() string {
	return fmt.Sprintf(
		"Value["+
			"value=%s, "+
			"label=%s, "+
			"units=%s, "+
			"help=%s, "+
			"min=%d "+
			"max=%d "+
			"isSet=%v]",
		C.GoString(self.cRef.value),
		C.GoString(self.cRef.label),
		C.GoString(self.cRef.units),
		C.GoString(self.cRef.help),
		self.cRef.min,
		self.cRef.max,
		self.cRef.isSet)
}
