//
// api.cpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//

#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, API * api)
{
  Notification * exported = exportNotification(api, notification);
  onNotificationWrapper(exported, api);
}


Manager * startManager(API * api)
{
  Manager * manager = newManager(OpenZWave::Manager::Create());
  manager->manager->AddWatcher( OnNotification, api );
  return manager;
}

void stopManager(API * api)
{
  Manager * manager = asManager(api);
  manager->manager->RemoveWatcher(OnNotification, api);
  freeManager(manager);
}

