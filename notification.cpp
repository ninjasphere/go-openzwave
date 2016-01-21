
#include "api.h"
static Notification * newNotification(uint8 notificationType)
{
  Notification * tmp = (Notification *)malloc(sizeof(Notification));
  *tmp = (Notification){0};
  tmp->notificationType = notificationType;
  return tmp;
}

// exports a C++ notification into a copy that can be consumed by the Go Layer. The consumer
// must propagate the responsibility to free the returned object to its callers.
Notification * exportNotification(API * api, OpenZWave::Notification const* notification)
{
  Notification * result = newNotification(notification->GetType());
  result->node = exportNode(api, notification->GetHomeId(), notification->GetNodeId());
  OpenZWave::Notification::NotificationType nt = notification->GetType();
  result->notificationCode = -1;

  switch (nt) {
  case OpenZWave::Notification::Type_ValueAdded:
  case OpenZWave::Notification::Type_ValueRemoved:
  case OpenZWave::Notification::Type_ValueChanged:
  case OpenZWave::Notification::Type_ValueRefreshed:
    result->value = exportValue(api, notification->GetHomeId(), notification->GetValueID());
    break;
  case OpenZWave::Notification::Type_Notification:
    result->notificationCode = notification->GetNotification();
    break;
  default:
    break;
    // do nothing
  }
  return result;
}


void freeNotification(Notification * notification)
{
  if (notification->node) {
    freeNode(notification->node);
  }
  if (notification->value) {
    freeValue(notification->value);
  }
  free(notification);
}
