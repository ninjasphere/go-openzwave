go-openzwave
============
A minimal 'Go' binding to the C++ openzwave library.

Note that only the subset of functionality required to support Ninja's current requirements is exposed.

Data
====
Data shared between the C++ and Go libraries is represented with pointers to C structures, defined in a C header file. Allocators and releasers are implemented in the C library.

Building
========
To build run `make` in the project directory. Once the dependency is build `go install` will suffice to rebuild it.

This build relies on the cgo tool to provide access to the openzwave C++ library. Cross-compilation and cgo cannot create linux/arm targets so to build
the linux/arm target you need to run the build natively on a linux/arm host.

Dependencies
============
openzwave - 1.0.791 - https://code.google.com/p/open-zwave/

Files
=====
api.h - C language header file
api/notification.h - bindings for notifications
api/node.h - bindings for nodes
