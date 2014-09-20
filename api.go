package openzwave

//
// Provides a facade for the C++ API that exposes just enough of the underlying C++
// API to be useful to implementing the Ninja Zwave driver
//
// The functions in this module are responsible for marshalling to and from the C functions
// declared in api.hpp and api.cpp.
//

//
// The following #cgo directives assumed that 'go' is a symbolic link that references the gopath that contains the current directory, e.g. ../../../..
//
// All 'go' packages that have this package as a dependency should include such a go link and will then inherit the library built in this package.
//

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
//
// #include "api.h"
import "C"

import "fmt"
import "os"
import "os/signal"
import "unsafe"

import "github.com/ninjasphere/go-openzwave/NT"
import "github.com/ninjasphere/go-openzwave/VT"
import "github.com/ninjasphere/go-openzwave/CODE"

type api struct {
	options       C.Options // an opaque reference to C++ Options object
	notifications chan Notification
	device        string
}

//
// The Phase0 -> Phase2 interface represent 3 different states in the evolution of the api from
// creation, through configuration, through use. 
//
// Each phase includes at least one method that allows transition to the next phase.
//
// Use of strong typing like this helps guide the consumer of the api package
// to construct a valid build sequence.
//

type Phase0 interface {
	AddIntOption(option string, value int) Phase0
	AddBoolOption(option string, value bool) Phase0
	SetDriver(device string) Phase0
	Run(loop EventLoop) int
}

type EventLoop func (chan Notification);

type Notification struct {
	notification *C.Notification
}

func (self Notification) String() string {
	return fmt.Sprintf(
		"Notification[\n"+
			"notificationType=%s,\n"+
			"notificationCode=%s,\n"+
			"homeId=0x%08x,\n"+
			"nodeId=0x%02x,\n"+
			"valueType=%s,\n"+
			"valueId=0x%08x]\n",
		NT.ToEnum(int(self.notification.notificationType)),
		CODE.ToEnum(int(self.notification.notificationCode)),
		self.notification.nodeId.homeId,
		self.notification.nodeId.nodeId,
		VT.ToEnum(int(self.notification.valueId.valueType)),
		self.notification.valueId.valueId)
}


// allocate the control block used to track the state of the API
func API(configPath string, userPath string, overrides string) Phase0 {
	var (
		cConfigPath *C.char = C.CString(configPath)
		cUserPath *C.char = C.CString(userPath)
		cOverrides *C.char = C.CString(overrides)
	)
	//defer C.free(unsafe.Pointer(cConfigPath))
	//defer C.free(unsafe.Pointer(cUserPath))
	//defer C.free(unsafe.Pointer(cOverrides))
	return api{
	       C.startOptions(cConfigPath, cUserPath, cOverrides), 
	       make(chan Notification),
	       defaultDriverName}
}

// configure the C++ Options object with an integer value
func (self api) AddIntOption(option string, value int) Phase0 {
	var cOption *C.char = C.CString(option)
	//defer C.free(unsafe.Pointer(cOption))

	C.addIntOption(self.options, cOption, C.int(value))
	return self
}

// configure the C++ Options object with a boolean value
func (self api) AddBoolOption(option string, value bool) Phase0 {
	var cOption *C.char = C.CString(option)

	//defer C.free(unsafe.Pointer(cOption))
	C.addBoolOption(self.options, cOption, C._Bool(value))
	return self
}

// add a driver.
func (self api) SetDriver(device string) Phase0 {
	if device != "" {
	       self.device = device
	}
	return self
}

func (self api) Run(loop EventLoop) int {

	C.endOptions(self.options)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	go func() {
		cDevice := C.CString(self.device)
		//defer C.free(unsafe.Pointer(cDevice))

		manager := C.startManager(cDevice, unsafe.Pointer(&self));
		defer C.stopManager(manager);

		loop(self.notifications);
	}()
	
	// Block until a signal is received.
	signal := <-signals
	fmt.Printf("received signal: %v", signal)
	return 1
}

//export OnNotificationWrapper
func OnNotificationWrapper(notification *C.Notification, context unsafe.Pointer) {
	(*api)(context).notifications <- Notification{notification}
}
