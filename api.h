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
#include "api/notification.h"

typedef void * Options;
typedef void * Manager;
extern int TRUE;

typedef void * const Context;

extern int FALSE;
extern int LogLevel_Detail;
extern int LogLevel_Debug;
extern int LogLevel_Info;
extern int LogLevel_Error;
extern Options createOptions(char *, char *);
extern void addIntOption(Options, char *, int );
extern void addBoolOption(Options, char *, int);
extern Manager lockOptions(Options );
extern void addDriver(Manager , char *);
extern void addWatcher(Manager, void *);
#ifdef __cplusplus
}
#endif

#ifdef __cplusplus
#include "Notification.h"
// C++ only parts
void OnNotification (OpenZWave::Notification const* _notification, void* _context);
#endif

#endif
