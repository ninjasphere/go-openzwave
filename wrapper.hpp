//
// wrapper.hpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//
// See wrapper.cpp for additional information.
//

#ifdef __cplusplus
extern "C" {
#endif
typedef void * Options;
typedef void * Manager;
extern int TRUE;
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
#ifdef __cplusplus
}
#endif
