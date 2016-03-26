#ifndef API_H
#define API_H
//
// api.h
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//
// See api.cpp for additional information.
//

#include <stdint.h>
#include <stdlib.h>
#include <stdbool.h>
#include "zwave_types.h"

#ifdef __cplusplus
#include <iostream>

//
// __cplusplus will be true only in code that is part of the implementation (C++ code)
//
// In these cases, we include the implementation C++ headers so that individual modules
// do not have to.
//

#include "Value.h"
#include "Manager.h"
#include "Notification.h"
#include "Options.h"
#include "platform/Log.h"

//
// The following ifdef magic is absolutely required to ensure that #cgo that doesn't understand how to parse C++ headers doesn't
// see the extern "C" declaration during the C parse of these headers. Failure to include this magic will result in link error
//
// See also http://stackoverflow.com/questions/1713214/how-to-use-c-in-go
//
extern "C" {
#else
#endif

#include "api/shareable.h"
typedef void API;
#include "api/manager.h"
#include "api/node.h"
#include "api/value.h"
#include "api/notification.h"
#include "api/options.h"


#ifdef __cplusplus
#include "_cgo_export.h"
}

#endif

#endif
