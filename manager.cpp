#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, API * api)
{
  Notification * exported = exportNotification(api, notification);
  onNotificationWrapper(exported, api);
}

void startManager(API * api)
{
  OpenZWave::Manager::Create();
  OpenZWave::Manager::Get()->AddWatcher( OnNotification, api );
}

void stopManager(API * api)
{
  OpenZWave::Manager::Get()->RemoveWatcher(OnNotification, api);
  OpenZWave::Manager::Destroy();
}

bool addDriver(char * device)
{
  return OpenZWave::Manager::Get()->AddDriver(device);
}

bool removeDriver(char * device)
{
  return OpenZWave::Manager::Get()->RemoveDriver(device);
}
