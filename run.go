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
func (a *api) Run() int {

	// lock the options object, now we are done configuring it

	C.endOptions()

	// allocate various channels we need

	signals := make(chan os.Signal, 1) // used to receive OS signals
	exit := make(chan int, 1)          // used to indicate we are ready to exit

	// indicate that we want to wait for these signals

	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		// Block until a signal is received.

		signal := <-signals
		// once we receive a signal, exit of the process is inevitable
		a.logger.Infof("received %v signal - commencing shutdown\n", signal)

		// try a graceful shutdown of the event loop
		a.shutdownDriver <- EXIT_INTERRUPTED
		// and the device monitor loop
		a.quitDeviceMonitor <- EXIT_INTERRUPTED

		// but, just in case this doesn't happen, set up an abort timer.
		time.AfterFunc(time.Second*5, func() {
			a.logger.Errorf("timed out while waiting for event loop to quit - aborting now\n")
			exit <- EXIT_INTERRUPT_FAILED
		})

		// the user is impatient - just die now
		signal = <-signals
		a.logger.Errorf("received 2nd %v signal - aborting now\n", signal)
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
		cSelf := a.C()

		C.startManager(cSelf) // start the manager
		defer C.stopManager(cSelf)

		cDevice := C.CString(a.device) // allocate a C string for device
		defer C.free(unsafe.Pointer(cDevice))

		// a function which returns true if the device exists
		deviceExists := func() bool {
			if _, err := os.Stat(a.device); err == nil {
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
		doneExit := 0
		for !done {
			select {
			case doneExit = <-a.quitDeviceMonitor: // we received a signal, allow us to quit
				done = true
			default:
				// one iteration of a device insert/removal cycle

				// wait until device present
				a.logger.Infof("waiting until %s is available\n", a.device)
				pollUntilDeviceExistsStateEquals(true)
				a.logger.Infof("device %s is available\n", a.device)

				go func() {

					// wait until device absent
					pollUntilDeviceExistsStateEquals(false)
					a.logger.Infof("device %s has been removed.\n", a.device)

					// start the removal of the driver
					a.shutdownDriver <- 0
				}()

				C.addDriver(cDevice)

				go func() {
					// wait until something (OS signal handler or device existence monitor) decides we need to terminate
					rc := <-a.shutdownDriver

					// we start an abort timer, because if the driver blocks, we need to restart the driver process
					// to guarantee successful operation.
					abortTimer := time.AfterFunc(5*time.Second, func() {
						a.logger.Errorf("failed to remove driver - exiting driver process\n")
						exit <- EXIT_QUIT_FAILED
					})

					// try to remove the driver
					if C.removeDriver(cDevice) {
						a.quitEventLoop <- rc
						abortTimer.Stop() // if we get to here in a timely fashion we can stop the abort timer
					} else {
						// this is unexpected, if we get to here, let the abort timer do its thing
						a.logger.Errorf("removeDriver call failed - waiting for abort\n")
					}
				}()

				rc := a.loop(a) // run the event loop

				if rc != 0 {
					done = true
					exit <- rc
					return
				}
			}
		}

		exit <- doneExit
	}()

	return <-exit
}

func (a *api) Shutdown(exit int) {

	select {
	case a.quitDeviceMonitor <- exit:
		break
	default:
	}

	select {
	case a.shutdownDriver <- exit:
		break
	default:
	}

	a.shareable.destroy()

}

//export onNotificationWrapper
func onNotificationWrapper(cNotification *C.Notification, context unsafe.Pointer) {
	// marshal from C to Go
	a := unmarshal(context).Go().(*api)
	goNotification := newGoNotification(cNotification)
	if a.callback != nil {
		a.callback(a, goNotification)
	}

	// forward the notification to the network
	a.getNetwork(goNotification.GetNode().GetHomeId()).notify(a, goNotification)

	// release the notification
	goNotification.free()
}
