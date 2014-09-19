#include "Manager.h"
#include "Notification.h"

#include "api.h"

extern "C" {
#include "_cgo_export.h"
}

static void OnNotification (OpenZWave::Notification const* notification, void* context)
{
  Notification * result = newNotification(notification->GetType());
  result->nodeId = (struct NodeId) {notification->GetHomeId(), notification->GetNodeId() };
  OnNotificationWrapper(result, context);
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
