//
// api.cpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//

#include <stdio.h>

#include "Manager.h"
#include "Options.h"
#include "platform/Log.h"
#include "api.h"

extern "C" {
#include "_cgo_export.h"
}

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

extern void addWatcher(Manager manager, void * context) 
{
 ((OpenZWave::Manager *)manager)->AddWatcher( OnNotification, context );
}

void OnNotification (OpenZWave::Notification const* notification, void* context)
{
  OpenZWave::ValueID id = notification->GetValueID();
  switch (notification->GetType()) {
  case OpenZWave::Notification::Type_Notification:
    switch (notification->GetNotification()) {
    case OpenZWave::Notification::Code_MsgComplete:
    case OpenZWave::Notification::Code_Timeout:
    case OpenZWave::Notification::Code_NoOperation:
    case OpenZWave::Notification::Code_Awake:
    case OpenZWave::Notification::Code_Sleep:
    case OpenZWave::Notification::Code_Dead:
    default:
      break;
    }
    break;
  case OpenZWave::Notification::Type_ValueAdded:
  case OpenZWave::Notification::Type_ValueRemoved:
  case OpenZWave::Notification::Type_ValueChanged:
  case OpenZWave::Notification::Type_ValueRefreshed:
  case OpenZWave::Notification::Type_Group:
  case OpenZWave::Notification::Type_NodeNew:
  case OpenZWave::Notification::Type_NodeAdded:
  case OpenZWave::Notification::Type_NodeRemoved:
  case OpenZWave::Notification::Type_NodeProtocolInfo:
  case OpenZWave::Notification::Type_NodeNaming:
  case OpenZWave::Notification::Type_NodeEvent:
  case OpenZWave::Notification::Type_PollingDisabled:
  case OpenZWave::Notification::Type_PollingEnabled:
  case OpenZWave::Notification::Type_SceneEvent:
  case OpenZWave::Notification::Type_CreateButton:
  case OpenZWave::Notification::Type_DeleteButton:
  case OpenZWave::Notification::Type_ButtonOn:
  case OpenZWave::Notification::Type_ButtonOff:
  case OpenZWave::Notification::Type_DriverReady:
  case OpenZWave::Notification::Type_DriverFailed:
  case OpenZWave::Notification::Type_DriverReset:
  case OpenZWave::Notification::Type_EssentialNodeQueriesComplete:
  case OpenZWave::Notification::Type_NodeQueriesComplete:
  case OpenZWave::Notification::Type_AwakeNodesQueried:
  case OpenZWave::Notification::Type_AllNodesQueriedSomeDead:
  case OpenZWave::Notification::Type_AllNodesQueried:
  default:
    break;
  }
  Notification * notificationT = (Notification *)malloc(sizeof(Notification));
  OnNotificationWrapper(notificationT, context);
}
