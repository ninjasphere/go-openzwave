package main

import (
	"flag"
	"os"

	openzwave "github.com/ninjasphere/go-openzwave"
	"github.com/ninjasphere/go-openzwave/LOG_LEVEL"
)

func main() {
	var (
		configDir    string
		logFileName  string
		save         bool
		debug        bool
		console      bool
		help         bool
		monitor      bool
		logLevel     int
		pollInterval int
		device       string
		logging      bool = false
	)

	flag.BoolVar(&monitor, "monitor", false, "Run the monitor")
	flag.StringVar(&configDir, "configDir", "../go-openzwave/openzwave/config", "Location of openzwave configuration directory")
	flag.BoolVar(&save, "save", false, "Save the configuration")
	flag.BoolVar(&debug, "debug", false, "Enable debugging")
	flag.BoolVar(&console, "console", false, "Enable console output")
	flag.IntVar(&pollInterval, "pollInterval", 10, "The polling interval")
	flag.StringVar(&logFileName, "logFileName", "zwcli.log", "Log file name")
	flag.BoolVar(&help, "help", false, "Print this help")
	flag.StringVar(&device, "device", "", "Device name /dev/ttyUSB0 on Linux, /dev/cu.SLAB_USBtoUART on OSX")
	flag.Parse()

	if help || !monitor {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if debug {
		logLevel = LOG_LEVEL.DEBUG
		logging = true
	} else {
		logLevel = LOG_LEVEL.NONE
		logFileName = "/dev/null"
	}

	callback := func(api openzwave.API, notification openzwave.Notification) {
		api.Logger().Infof("%v\n", notification)
	}

	os.Exit(openzwave.
		BuildAPI(configDir, "", "").
		AddBoolOption("SaveConfiguration", save).
		AddBoolOption("logging", logging).
		AddStringOption("LogFileName", logFileName, false).
		AddBoolOption("ConsoleOutput", console).
		AddBoolOption("NotifyTransactions", true).
		AddIntOption("SaveLogLevel", logLevel).
		AddIntOption("QueueLogLevel", logLevel).
		AddIntOption("DumpTrigger", logLevel).
		AddIntOption("PollInterval", pollInterval).
		AddBoolOption("IntervalBetweenPolls", true).
		AddBoolOption("ValidateValueChanges", true).
		SetDeviceName(device).
		SetNotificationCallback(callback).
		Run())
}
