package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
// #include "api.h"
import "C"

// This interface is used to configure the API by setting various options,
// and the controller device name. When configuration is finished, call the Run method
// with an EventLoop function.
//
// To obtain a reference to this interface, call BuildAPI().
//
// The result of the Run function can be passed as a parameter to Os.Exit().
//
//
type Configurator interface {
	//Configure the logger implementation.
	SetLogger(Logger) Configurator

	//Configure the synchronous callback.
	SetCallback(Callback) Configurator

	//Configure the event loop function
	SetEventLoop(EventLoop) Configurator

	// Add an integer option.
	AddIntOption(option string, value int) Configurator

	// Add a boolean option.
	AddBoolOption(option string, value bool) Configurator

	// Add a string option.
	AddStringOption(option string, value string, append bool) Configurator

	// Set the device name used by the driver.
	SetDeviceName(device string) Configurator

	// Run the event loop forever
	Run() int
}

//
// Begin the construction of the API by returning a Configurator
//
// configPath is the name of the directory containing openzwave configuration files;
// userPath is the name of the directory containing user specific openzwave configuration files;
// overrides are command line options (uptto --) that can overide the configuration provided by the previous two options.
//
// For more information about these parameters, refer to the documentation for the C++ OpenZWave::Options class.
//
func BuildAPI(configPath string, userPath string, overrides string) Configurator {
	var (
		cConfigPath *C.char = C.CString(configPath)
		cUserPath   *C.char = C.CString(userPath)
		cOverrides  *C.char = C.CString(overrides)
	)
	//defer C.free(unsafe.Pointer(cConfigPath))
	//defer C.free(unsafe.Pointer(cUserPath))
	//defer C.free(unsafe.Pointer(cOverrides)
	C.startOptions(cConfigPath, cUserPath, cOverrides)
	return &api{
		defaultEventLoop,
		nil,
		defaultDriverName,
		make(chan Signal, 0),
		&defaultLogger{}}
}

// configure the C++ Options object with an integer value
func (self *api) AddIntOption(option string, value int) Configurator {
	var cOption *C.char = C.CString(option)
	//defer C.free(unsafe.Pointer(cOption))

	C.addIntOption(cOption, C.int(value))
	return self
}

// configure the C++ Options object with a boolean value
func (self *api) AddBoolOption(option string, value bool) Configurator {
	var cOption *C.char = C.CString(option)

	//defer C.free(unsafe.Pointer(cOption))
	C.addBoolOption(cOption, C._Bool(value))
	return self
}

// configure the C++ Options object with a string value
func (self *api) AddStringOption(option string, value string, append bool) Configurator {
	var cOption *C.char = C.CString(option)

	//defer C.free(unsafe.Pointer(cOption))
	C.addStringOption(cOption, C.CString(value), C._Bool(append))
	return self
}

// set the device name
func (self *api) SetDeviceName(device string) Configurator {
	if device != "" {
		self.device = device
	}
	return self
}

// set the logger
func (self *api) SetLogger(logger Logger) Configurator {
	self.logger = logger
	return self
}

// Clients of the API should provide Configuration.Run() with an implementation of this type to
// handle the notifications generated by the API.
//
// The implementor can return a non-zero code to indicate that the process should exit now. 0 means
// that the loop can be restarted, if required.
//
type EventLoop func(API) int

// set the event loop
func (self *api) SetEventLoop(loop EventLoop) Configurator {
	self.loop = loop
	return self
}

// A type of function that can receive notifications from the OpenZWave library when they occur.
//
// This callback is processed synchronously. This means that the implementor:
// MUST NOT block; MUST NOT hand the reference to the notification to a goroutine which lives beyond
// the callback call; MUST NOT store the reference to the notification in a structure that lives beyond
// the duration of the callback call.
type Callback func(API, Notification)

// set the synchronous call back
func (self *api) SetCallback(callback Callback) Configurator {
	self.callback = callback
	return self
}