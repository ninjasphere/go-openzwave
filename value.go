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

type Value interface {
	SetUint8(value uint8) bool
	GetUint8() (uint8, bool)
}

type value struct {
	cRef *C.Value
}

type missingValue struct {
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
			"min=%d, "+
			"max=%d, "+
			"isSet=%v]",
		VT.ToEnum(int(self.cRef.valueId.valueType)),
		CC.ToEnum(int(self.cRef.valueId.commandClassId)),
		uint(self.cRef.valueId.instance),
		uint(self.cRef.valueId.index),
		C.GoString(self.cRef.value),
		C.GoString(self.cRef.label),
		C.GoString(self.cRef.units),
		C.GoString(self.cRef.help),
		(int32)(self.cRef.min),
		(int32)(self.cRef.max),
		(bool)(self.cRef.isSet))
}

func newGoValue(cRef *C.Value) *value {
	return &value{cRef}
}

func (self *value) notify(api *api, nt *notification) {
	// TODO
}

func (self *value) SetUint8(value uint8) bool {
	return (bool)(C.setUint8Value(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), C.uint8_t(value)))
}

func (self *value) GetUint8() (uint8, bool) {
	var value uint8
	ok := (bool)(C.getUint8Value(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), (*C.uint8_t)(&value)))
	return value, ok
}

// for a missing value, the set operation always fails
func (self *missingValue) SetUint8(value uint8) bool {
	return false
}

// for a missing value, the get operation always fails
func (self *missingValue) GetUint8() (uint8, bool) {
	return 0, false
}
