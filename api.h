#ifndef API_H 
#define API_H
//
// api.h
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//
// See api.cpp for additional information.
//

//
// The following ifdef magic is absolutely required to ensure that #cgo that doesn't understand how to parse C++ headers doesn't
// see the extern "C" declaration during the C parse of these headers. Failure to include this magic will result in link error
//
// See also http://stackoverflow.com/questions/1713214/how-to-use-c-in-go
//
#ifdef __cplusplus
extern "C" {
#endif 
#include <stdint.h>

#include "api/node.h"
#include "api/value.h"
#include "api/notification.h"
#include "api/manager.h"
#include "api/options.h"


extern int TRUE;
extern int FALSE;

extern int FALSE;
extern int LogLevel_Detail;
extern int LogLevel_Debug;
extern int LogLevel_Info;
extern int LogLevel_Error;
#ifdef __cplusplus
}
#endif

#endif
