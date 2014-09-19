#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, void* context)
{
	Notification * exported = exportNotification(notification);
	OnNotificationWrapper(exported, context);
}

void addDriver(Manager manager, char * device)
{
  ((OpenZWave::Manager *)manager)->AddDriver(device);
}

void addWatcher(Manager manager, void * context) 
{
 ((OpenZWave::Manager *)manager)->AddWatcher( OnNotification, context );
}

Manager createManager()
{
  return (Manager)OpenZWave::Manager::Create();
}
