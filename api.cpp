//
// api.cpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//

#include "api.h"

int TRUE = 1;
int FALSE = 0;

int LogLevel_Detail = (int)(OpenZWave::LogLevel_Detail);
int LogLevel_Debug = (int)(OpenZWave::LogLevel_Debug);
int LogLevel_Error = (int)(OpenZWave::LogLevel_Error);
int LogLevel_Info = (int)(OpenZWave::LogLevel_Info);



