#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, void* context)
{
  Notification * exported = exportNotification(notification);
  onNotificationWrapper(exported, context);
}

Manager startManager(void * context)
{
  OpenZWave::Manager * manager = OpenZWave::Manager::Create();
  manager->AddWatcher( OnNotification, context );
  return (struct Manager) { manager };
}

void stopManager(Manager manager, void *context)
{
  manager.manager->RemoveWatcher(OnNotification, context);
  OpenZWave::Manager::Destroy();
}

bool addDriver(Manager manager, char * device)
{
  return manager.manager->AddDriver(device);
}

bool removeDriver(Manager manager, char * device)
{
  return manager.manager->RemoveDriver(device);
}
