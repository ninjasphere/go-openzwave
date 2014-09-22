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
	inC *C.ValueID
}

func (self ValueID) String() string {
	return fmt.Sprintf(
		"ValueID["+
			"type=%s, "+
			"commandClassId=%s, "+
			"instance=%d, "+
			"index=%d]",
		VT.ToEnum(int(self.inC.valueType)),
		CC.ToEnum(int(self.inC.commandClassId)),
		uint(self.inC.instance),
		uint(self.inC.index))
}

type Value struct {
	inC *C.Value
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
		C.GoString(self.inC.value),
		C.GoString(self.inC.label),
		C.GoString(self.inC.units),
		C.GoString(self.inC.help),
		self.inC.min,
		self.inC.max,
		self.inC.isSet)
}
