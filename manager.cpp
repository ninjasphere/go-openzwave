#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, void* context)
{
  Manager * manager = asManager(context);
  Notification * exported = exportNotification(manager, notification);
  onNotificationWrapper(exported, context);
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

Manager * startManager(void * context)
{
  Manager * manager = newManager(OpenZWave::Manager::Create());
  manager->manager->AddWatcher( OnNotification, context );
  return manager;
}

void stopManager(Manager * manager, void *context)
{
  manager->manager->RemoveWatcher(OnNotification, context);
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
