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
	CommandClassId uint8
	Instance       uint8
	Index          uint8
}

type Value interface {
	Id() ValueID
	SetUint8(value uint8) bool
	GetUint8() (uint8, bool)
	SetBool(value bool) bool
	GetBool() (bool, bool)
	Refresh() bool
	SetPollingState(bool) bool
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

func (self *value) Id() ValueID {
	return ValueID{
		CommandClassId: uint8(self.cRef.valueId.commandClassId),
		Instance:       uint8(self.cRef.valueId.instance),
		Index:          uint8(self.cRef.valueId.index),
	}
}

func (self *value) SetUint8(value uint8) bool {
	return (bool)(C.setUint8Value(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), C.uint8_t(value)))
}

func (self *value) GetUint8() (uint8, bool) {
	var value C.uint8_t
	ok := (bool)(C.getUint8Value(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), (*C.uint8_t)(&value)))
	return (uint8)(value), ok
}

func (self *value) SetBool(value bool) bool {
	return (bool)(C.setBoolValue(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), C._Bool(value)))
}

func (self *value) GetBool() (bool, bool) {
	var value C._Bool
	ok := (bool)(C.getBoolValue(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), (*C._Bool)(&value)))
	return (bool)(value), ok
}

func (self *value) Refresh() bool {
	return (bool)(C.refreshValue(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id)))
}

func (self *value) free() {
	C.freeValue(self.cRef)
}

func (self *value) SetPollingState(state bool) bool {
	return (bool)(C.setPollingState(C.uint32_t(self.cRef.homeId), C.uint64_t(self.cRef.valueId.id), C._Bool(state)))
}

// for a missing value, the set operation always fails
func (self *missingValue) SetUint8(value uint8) bool {
	return false
}

// for a missing value, the get operation always fails
func (self *missingValue) GetUint8() (uint8, bool) {
	return 0, false
}

// for a missing value, the set operation always fails
func (self *missingValue) SetBool(value bool) bool {
	return false
}

// for a missing value, the get operation always fails
func (self *missingValue) GetBool() (bool, bool) {
	return false, false
}

// for a missing value, the get operation always fails
func (self *missingValue) Refresh() bool {
	return false
}

func (self *missingValue) SetPollingState(state bool) bool {
	return false
}

func (self *missingValue) Id() ValueID {
	return ValueID{0, 0, 0}
}
