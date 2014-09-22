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

import (
	"os"
	"os/signal"
	"time"
	"unsafe"
)

// A value-less type that is used to represent signals generated by the API, particularly quit signals used to
// ask an EventLoop to quit.
type Signal struct{}

type api struct {
	options       C.Options // an opaque reference to C++ Options object
	notifications chan *Notification
	device        string
	quit          chan Signal
	manager       *C.Manager
	logger        Logger
}

//
// Begin the construction of the API by returning a Configurator
//
// configPath is the name of the directory containing openzwave configuration files.
//
// userPath is the name of the directory containing user specific openzwave configuration files.
//
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
	//defer C.free(unsafe.Pointer(cOverrides))
	return api{
		C.startOptions(cConfigPath, cUserPath, cOverrides),
		make(chan *Notification),
		defaultDriverName,
		make(chan Signal, 0),
		nil,
		defaultLogger{}}
}

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
	//Configure the logger implementation
	SetLogger(Logger) Configurator
	// Add an integer option.
	AddIntOption(option string, value int) Configurator
	// Add a boolean ``option.
	AddBoolOption(option string, value bool) Configurator
	// Set the device name used by the driver.
	SetDeviceName(device string) Configurator
	// Conclude the configuration and start running the supplied event loop. The
	// body of the EventLoop has access to an API reference.
	Run(loop EventLoop) int
}

//
// The API interface is available to implementors of the EventLoop type when the
// Configurator.Run() method is called.
//
type API interface {
	// notifications are received on this channel
	Notifications() chan *Notification

	// The EventLoop should return from the function when a signal is received on this channel
	QuitSignal() chan Signal

	// the API logger
	Logger() Logger
}

// Clients of the API should provide Configuration.Run() with an implementation of this type to
// handle the notifications generated by the API.
type EventLoop func(API)

// configure the C++ Options object with an integer value
func (self api) AddIntOption(option string, value int) Configurator {
	var cOption *C.char = C.CString(option)
	//defer C.free(unsafe.Pointer(cOption))

	C.addIntOption(self.options, cOption, C.int(value))
	return self
}

// configure the C++ Options object with a boolean value
func (self api) AddBoolOption(option string, value bool) Configurator {
	var cOption *C.char = C.CString(option)

	//defer C.free(unsafe.Pointer(cOption))
	C.addBoolOption(self.options, cOption, C._Bool(value))
	return self
}

// set the device name
func (self api) SetDeviceName(device string) Configurator {
	if device != "" {
		self.device = device
	}
	return self
}

func (self api) SetLogger(logger Logger) Configurator {
	self.logger = logger
	return self
}

//
// Run the supplied event loop
//
// The intent of the complexity is to gracefully handle device insertion and removal events and to
// deal with unexpected (but observed) lockups during the driver removal processing.
//
// The function will only return if a signal is received or if there was an unexpected
// lockup during driver removal processing. The exit code identifies which path
// caused the exit to occur.
//
func (self api) Run(loop EventLoop) int {

	// lock the options object, now we are done configuring it

	C.endOptions(self.options)

	// allocate various channels we need

	signals := make(chan os.Signal, 1)        // used to receive OS signals
	startQuit := make(chan Signal, 2)         // used to indicate we need to quit the event loop
	quitDeviceMonitor := make(chan Signal, 1) // used to indicate to outer loop that it should exit
	exit := make(chan int, 1)                 // used to indicate we are ready to exit

	// indicate that we want to wait for these signals

	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		// Block until a signal is received.

		signal := <-signals
		// once we receive a signal, exit of the process is inevitable
		self.logger.Infof("received %v signal - commencing shutdown\n", signal)

		// try a graceful shutdown of the event loop
		startQuit <- Signal{}
		// and the device monitor loop
		quitDeviceMonitor <- Signal{}

		// but, just in case this doesn't happen, set up an abort timer.
		time.AfterFunc(time.Second*5, func() {
			self.logger.Errorf("timed out while waiting for event loop to quit - aborting now\n")
			exit <- 1
		})

		// the user is impatient - just die now
		signal = <-signals
		self.logger.Errorf("received 2nd %v signal - aborting now\n", signal)
		exit <- 2
	}()

	//
	// This goroutine does the following
	//    starts the manager
	//    starts a device monitoroing loop which
	//       waits for the device to be available
	// 	 starts a device removal goroutine which raises a startQuit signal when removal of the device is detected
	//   	 starts the driver
	//	 starts a go routine that that waits until a startQuit is signaled, then initiates the removal of the driver and quit of the event loop
	//	 runs the event loop
	//
	// It does not exit until either an OS Interrupt or Kill signal is received or driver removal or event loop blocks for some reason.
	//
	// If the device is removed, the monitoring go routine will send a signal via the startQuit channel. The intent is to allow the
	// event loop to exit and have the driver removed.
	//
	// The driver removal goroutine waits for the startQuit signal, then attempts to remove the driver. If this completes successfully
	// it propagates a quit signal to the event loop. It also sets up an abort timer which will exit the process if either
	// the driver removal or quit signal propagation blocks for some reason.
	//
	// If an OS signal is received, the main go routine will send signals to the startQuit and to the quitDeviceMonitor channels.
	// It then waits for another signal, for the outer loop to exit or for a 5 second timeout. When one of these occurs, the
	// process will exit.
	//

	go func() {
		cSelf := unsafe.Pointer(&self) // a reference to self

		C.startManager(cSelf) // start the manager
		defer C.stopManager(cSelf)

		cDevice := C.CString(self.device) // allocate a C string for device
		defer C.free(unsafe.Pointer(cDevice))

		// a function which returns true if the device exists
		deviceExists := func() bool {
			if _, err := os.Stat(self.device); err == nil {
				return true
			} else {
				if os.IsNotExist(err) {
					return false
				} else {
					return true
				}
			}
		}

		// waits until the state matches the desired state.
		pollUntilDeviceExistsStateEquals := func(comparand bool) {
			for deviceExists() != comparand {
				time.Sleep(time.Second)
			}
		}

		// there is one iteration of this loop for each device insertion/removal cycle
		done := false
		for !done {
			select {
			case doneSignal := <-quitDeviceMonitor: // we received a signal, allow us to quit
				_ = doneSignal
				done = true
			default:
				// one iteration of a device insert/removal cycle

				// wait until device present
				self.logger.Infof("waiting until %s is available\n", self.device)
				pollUntilDeviceExistsStateEquals(true)

				go func() {

					// wait until device absent
					pollUntilDeviceExistsStateEquals(false)
					self.logger.Infof("device %s removed\n", self.device)

					// start the removal of the driver
					startQuit <- Signal{}
				}()

				C.addDriver(self.manager, cDevice)

				go func() {
					// wait until something (OS signal handler or device existence monitor) decides we need to terminate
					<-startQuit

					// we start an abort timer, because if the driver blocks, we need to restart the driver process
					// to guarantee successful operation.
					abortTimer := time.AfterFunc(5*time.Second, func() {
						self.logger.Errorf("failed to remove driver - exiting driver process\n")
						exit <- 3
					})

					// try to remove the driver
					if C.removeDriver(self.manager, cDevice) {
						self.quit <- Signal{}
						abortTimer.Stop() // if we get to here in a timely fashion we can stop the abort timer
					} else {
						// this is unexpected, if we get to here, let the abort timer do its thing
						self.logger.Errorf("removeDriver call failed - waiting for abort\n")
					}
				}()

				loop(self) // run the event loop
			}
		}

		exit <- 4
	}()

	return <-exit
}

func (self api) Notifications() chan *Notification {
	return self.notifications
}

func (self api) QuitSignal() chan Signal {
	return self.quit
}

func (self api) Logger() Logger {
	return self.logger
}

//export onNotificationWrapper
func onNotificationWrapper(notification *C.Notification, context unsafe.Pointer) {
	self := (*api)(context)
	self.notifications <- (*Notification)(notification.goRef)
}

//export asManager
func asManager(context unsafe.Pointer) *C.Manager {
	self := (*api)(context)
	return self.manager
}

//export setManager
func setManager(context unsafe.Pointer, manager *C.Manager) {
	api := (*api)(context)
	api.manager = manager
}
