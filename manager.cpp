#include "api.h"

// forwards the notification from the C++ API to the Go layer - caller must free.
static void OnNotification (OpenZWave::Notification const* notification, void* context)
{
	Notification * exported = exportNotification(notification);
	OnNotificationWrapper(exported, context);
}

extern Manager startManager(char * device, void * context)
{
	OpenZWave::Manager * manager = OpenZWave::Manager::Create();
	manager->AddWatcher( OnNotification, context );
	manager->AddDriver(device);
	return (Manager) manager;
}
