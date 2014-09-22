#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, API * api)
{
  Manager * manager = asManager(api);
  Notification * exported = exportNotification(manager, notification);
  onNotificationWrapper(exported, api);
}

static Manager * newManager(OpenZWave::Manager * cppObj) {
  Manager * manager = (Manager *)malloc((sizeof *manager));
  *manager = (struct Manager){ cppObj };
  return manager;
}

static void freeManager(Manager * manager)
{
  OpenZWave::Manager::Destroy();
  free(manager);
}

Manager * startManager(API * api)
{
  Manager * manager = newManager(OpenZWave::Manager::Create());
  manager->manager->AddWatcher( OnNotification, api );
  return manager;
}

void stopManager(Manager * manager, API * api)
{
  manager->manager->RemoveWatcher(OnNotification, api);
  freeManager(manager);
}

bool addDriver(Manager * manager, char * device)
{
  return manager->manager->AddDriver(device);
}

bool removeDriver(Manager * manager, char * device)
{
  return manager->manager->RemoveDriver(device);
}
