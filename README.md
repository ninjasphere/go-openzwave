go-openzwave
============
A minimal 'Go' binding to the C++ openzwave library.

Note that only the subset of functionality required to support Ninja's current requirements is exposed.

Building
========
To build run `make` in the project directory. Once the non-Go dependencies are built `go install` will suffice to rebuild it.

This build relies on the cgo tool to provide access to the openzwave C++ library. Cross-compilation and cgo cannot create linux/arm targets so to build
the linux/arm target you need to run the build natively on a linux/arm host.

Abstractions
============
There is a set of C functions and structure declarations for each abstraction in C++ that the Go Layer needs access to. The type specific declarations are located in a file called api/{type}.h.

Each declaration file has two sections - a C section which will be consumed by the 'Go' cgo tool and the implementation and a C++ section which will be consumed only by the implementation.
The C++ section must be guarded by an #ifdef __cplusplus directive in order to prevent it being parsed by the C parser, in particular the C parser that the #cgo tool will invoke.

References to implementation types should only appear in the C++ section of each type declaration.

Any includes for OpenWave C++ header files should be moved into the "implementation dependencies" section of api.h and guarded with #ifdef __cplusplus directives.

As a rule, .cpp files should only need to include "api.h"

Memory Management
=================
The C structures must not refer to any C++ types or objects. As a rule, there should be a triplet of functions `new{Type}()`, `free{Type}()` and `export{Type}()` for each abstraction.

The `export{Type}` must perform a deep copy of all required parts of the C++ representation to the C representation that can be used by `Go`. It should call the `new{Type}` function of its
own type and the `export{Type}` functions of related types. The object returned by the export function must not contain any references to memory that are also retained by the C or C++ layer.

The `Go` consumer must guarantee to call the `free{Type}` function when an exported object is no longer required. The `free{Type}` function should deallocate any allocations performed by the corresponding `export{Type}` function during the original export of the object, usually by call free(void *) on directly allocated memory and `free{Type}` on memory allocated by the `export{RelatedType}` functions
of related types.

Code Generation
===============
This module uses bash scripts to analyse the OpenZWave++ code and generate typesafe Go enumerations for the relevant constants. These generation steps run as part of the 'deps' step of the project Makefile. Note that use of code generation for the sub-packages means that `go get -d` will fail for these projects because of a chicken and egg problem. This can be fixed by manually checking out this project and then running `make deps` before running `go get -d` in the consuming project.

Dependencies
============
openzwave - 1.4 - https://github.com/OpenZWave/open-zwave/tree/v1.4

Files
=====
* api.h - a two-part (C and C++) header file. Should be the only include required by the implementation. Implementation types are restricted to the C++ part of the file.
* api/{type}.h - a two-part (C and C++) header that contains the shared C abstraction and type specific C++ functions used by the implementation of a specific C abstraction.
* {type}.cpp - a C++ implementation of the functions declared in api/{type}.h
* api.go - 'Go' glue-code required to marhsall to and from the functions declared in api/*.h
* scripts/Generate{PREFIX}.sh - code generators for various type safe enumerations extracted from the OpenZWave source.
