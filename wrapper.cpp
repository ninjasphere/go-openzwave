//
// wrapper.cpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//

#include <stdio.h>

#include "Manager.h"
#include "Options.h"
#include "platform/Log.h"
#include "wrapper.hpp"
#include "watcher.hpp"

int TRUE = 1;
int FALSE = 0;
int LogLevel_Detail = (int)(OpenZWave::LogLevel_Detail);
int LogLevel_Debug = (int)(OpenZWave::LogLevel_Debug);
int LogLevel_Error = (int)(OpenZWave::LogLevel_Error);
int LogLevel_Info = (int)(OpenZWave::LogLevel_Info);

extern Options createOptions(char * config, char * log)
{
  OpenZWave::Options::Create(config, log, "");
  return (Options*) OpenZWave::Options::Get();
}

extern void addIntOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionInt(option, value);
}

extern void addBoolOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionBool(option, value == TRUE ? true : false);
}

extern Manager lockOptions(Options options)
{
  ((OpenZWave::Options *)options)->Lock();
  return (Manager *)OpenZWave::Manager::Create();
}

extern void addDriver(Manager manager, char * device)
{
  ((OpenZWave::Manager *)manager)->AddDriver(device);
}

extern void addWatcher(Manager manager) 
{
  ((OpenZWave::Manager *)manager)->AddWatcher( OpenZWave::OnNotification, NULL );
}
