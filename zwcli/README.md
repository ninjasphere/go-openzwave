NAME
====
zwcli - a command line API for the Go API defined in go-openzwave

DESCRIPTION
===========
The purpose of this Go program is to provide a tool to aid the development and testing of go-openzwave by providing a way of exercising 
the API independently of the Ninja driver (driver-go-zwave). The tool might be useful in its own right for interrogating the zwave network.

BUILDING
========
Run `make` from the command line at least once. Once built, further builds can be created with `go install`.

Note that due to cross-compilation issues with cgo, it is necessary to build the linux/arm target on a native linux/arm host.

