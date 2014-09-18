#ifndef WRAPPER_HPP
#define WRAPPER_HPP
//
// wrapper.hpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//
// See wrapper.cpp for additional information.
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
typedef void * Options;
typedef void * Manager;
extern int TRUE;

typedef struct {
  uint32_t homeId;
  uint8_t  nodeId;
} NodeId;

typedef struct {
  NodeId nodeId;
} Node;

typedef struct {
  uint8_t   notificationType;
  uint8_t   notificationCode;
  Node     node;
} Notification;

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
#endif
