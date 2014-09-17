package openzwave

//
// Provides a facade for the C++ API that exposes just enough of the underlying C++
// API to be useful to implementing the Ninja Zwave driver
//
// The functions in this module are responsible for marshalling to and from the C functions
// declared in wrapper.hpp and wrapper.cpp.
//

//
// #cgo LDFLAGS: -lopenzwave -L../go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include <stdlib.h>
// #include "wrapper.hpp"
import "C"

import "unsafe"

var LogLevel_Detail int = int(C.LogLevel_Detail)
var LogLevel_Error int = int(C.LogLevel_Error)
var LogLevel_Debug int = int(C.LogLevel_Debug)
var LogLevel_Info int = int(C.LogLevel_Info)

type API struct {
	options C.Options // an opaque reference to C++ Options object
	manager C.Manager // an opaque reference to C++ Manager opject
}

// allocate the control block used to track the state of the API
func NewAPI() *API {
	return &API{nil, nil}
}

// create and stash the C++ Options object
func (self *API) CreateOptions(configPath string, logPath string) *API {
	var cConfigPath *C.char = C.CString(configPath)
	var cLogPath *C.char = C.CString(logPath)
	defer C.free(unsafe.Pointer(cConfigPath))
	defer C.free(unsafe.Pointer(cLogPath))
	self.options = C.createOptions(cConfigPath, cLogPath)
	return self
}

// configure the C++ Options object with an integer value
func (self *API) AddIntOption(option string, value int) *API {
	var cOption *C.char = C.CString(option)
	defer C.free(unsafe.Pointer(cOption))

	C.addIntOption(self.options, cOption, C.int(value))
	return self
}

// configure the C++ Options object with a boolean value
func (self *API) AddBoolOption(option string, value bool) *API {
	var cOption *C.char = C.CString(option)
	var cBool C.int

	defer C.free(unsafe.Pointer(cOption))
	if value {
		cBool = C.TRUE
	} else {
		cBool = C.FALSE
	}
	C.addBoolOption(self.options, cOption, cBool)
	return self
}

// lock the options object and allocate the manager object.
func (self *API) LockOptions() *API {
	self.manager = C.lockOptions(self.options)
	return self
}

// add a driver.
func (self *API) AddDriver(device string) *API {
	var cDevice *C.char = C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	C.addDriver(self.manager, cDevice)
	return self
}

// add a watcher
func (self *API) AddWatcher(channel chan *interface{}) *API {
	C.addWatcher(self.manager)
	return self
}
