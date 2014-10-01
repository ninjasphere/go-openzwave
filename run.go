package openzwave

// #cgo LDFLAGS: -lopenzwave -Lgo/src/github.com/ninjasphere/go-openzwave/openzwave
// #cgo CPPFLAGS: -Iopenzwave/cpp/src/platform -Iopenzwave/cpp/src -Iopenzwave/cpp/src/value_classes
// #include "api.h"
import "C"

import (
	"os"
	"os/signal"
	"reflect"
	"time"
	"unsafe"
)

const (
	EXIT_QUIT_FAILED       = 127 // the event loop did not exit
	EXIT_INTERRUPTED       = 126 // something interrupted the current process
	EXIT_INTERRUPTED_AGAIN = 125 // something interrupted the current process (twice)
	EXIT_INTERRUPT_FAILED  = 124 // something interrupted the current process, but something took too long to clean up
	EXIT_NODE_REMOVED      = 123
)

var defaultEventLoop = func(api API) int {
	for {
		select {
		case quitNow := <-api.QuitSignal():
			api.Logger().Debugf("terminating event loop in response to quit.\n")
			return quitNow
		}
	}
}

var defaultEventCallback = func(api API, event Event) {
	api.Logger().Debugf("received event %v - %v\n", reflect.TypeOf(event), event)
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
func (self *api) Run() int {

	// lock the options object, now we are done configuring it

	C.endOptions()

	// allocate various channels we need

	signals := make(chan os.Signal, 1)        // used to receive OS signals
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
		self.shutdownDriver <- EXIT_INTERRUPTED
		// and the device monitor loop
		quitDeviceMonitor <- Signal{}

		// but, just in case this doesn't happen, set up an abort timer.
		time.AfterFunc(time.Second*5, func() {
			self.logger.Errorf("timed out while waiting for event loop to quit - aborting now\n")
			exit <- EXIT_INTERRUPT_FAILED
		})

		// the user is impatient - just die now
		signal = <-signals
		self.logger.Errorf("received 2nd %v signal - aborting now\n", signal)
		exit <- EXIT_INTERRUPTED_AGAIN
	}()

	//
	// This goroutine does the following
	//    starts the manager
	//    starts a device monitoroing loop which
	//       waits for the device to be available
	// 	 starts a device removal goroutine which raises a shutdownDriver signal when removal of the device is detected
	//   	 starts the driver
	//	 starts a go routine that that waits until a shutdownDriver is signaled, then initiates the removal of the driver and quit of the event loop
	//	 runs the event loop
	//
	// It does not exit until either an OS Interrupt or Kill signal is received or driver removal or event loop blocks for some reason.
	//
	// If the device is removed, the monitoring go routine will send a signal via the shutdownDriver channel. The intent is to allow the
	// event loop to exit and have the driver removed.
	//
	// The driver removal goroutine waits for the shutdownDriver signal, then attempts to remove the driver. If this completes successfully
	// it propagates a quit signal to the event loop. It also sets up an abort timer which will exit the process if either
	// the driver removal or quit signal propagation blocks for some reason.
	//
	// If an OS signal is received, the main go routine will send signals to the shutdownDriver and to the quitDeviceMonitor channels.
	// It then waits for another signal, for the outer loop to exit or for a 5 second timeout. When one of these occurs, the
	// process will exit.
	//

	go func() {
		cSelf := unsafe.Pointer(self) // a reference to self

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
					self.shutdownDriver <- 0
				}()

				C.addDriver(cDevice)

				go func() {
					// wait until something (OS signal handler or device existence monitor) decides we need to terminate
					rc := <-self.shutdownDriver

					// we start an abort timer, because if the driver blocks, we need to restart the driver process
					// to guarantee successful operation.
					abortTimer := time.AfterFunc(5*time.Second, func() {
						self.logger.Errorf("failed to remove driver - exiting driver process\n")
						exit <- EXIT_QUIT_FAILED
					})

					// try to remove the driver
					if C.removeDriver(cDevice) {
						self.quitEventLoop <- rc
						abortTimer.Stop() // if we get to here in a timely fashion we can stop the abort timer
					} else {
						// this is unexpected, if we get to here, let the abort timer do its thing
						self.logger.Errorf("removeDriver call failed - waiting for abort\n")
					}
				}()

				rc := self.loop(self) // run the event loop
				if rc != 0 {
					done = true
					exit <- rc
					return
				}
			}
		}

		exit <- EXIT_INTERRUPTED
	}()

	return <-exit
}

func (self *api) Shutdown(exit int) {
	select {
	case self.shutdownDriver <- exit:
		break
	default:
	}
}

//export onNotificationWrapper
func onNotificationWrapper(cNotification *C.Notification, context unsafe.Pointer) {
	// marshal from C to Go
	self := (*api)(context)
	goNotification := newGoNotification(cNotification)
	if self.callback != nil {
		self.callback(self, goNotification)
	}

	// forward the notification to the network
	self.getNetwork(goNotification.GetNode().GetHomeId()).notify(self, goNotification)

	// release the notification
	goNotification.free()
}
