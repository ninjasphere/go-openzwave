
#include "api.h"
static Notification * newNotification(uint8_t notificationType)
{
  Notification * tmp = (Notification *)malloc(sizeof(Notification));
  *tmp = (Notification){0};
  tmp->notificationType = notificationType;
  tmp->goRef = newGoNotification(tmp);
  return tmp;
}

// exports a C++ notification into a copy that can be consumed by the Go Layer. The consumer
// must propagate the responsibility to free the returned object to its callers.
Notification * exportNotification(API * api, OpenZWave::Notification const* notification)
{
  Notification * result = newNotification(notification->GetType());
  result->node = exportNode(api, notification->GetHomeId(), notification->GetNodeId());
  result->notificationCode =
    notification->GetType() == OpenZWave::Notification::Type_Notification
    ? notification->GetNotification()
    : -1;
  result->value = exportValue(api, notification->GetHomeId(), notification->GetValueID());
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
