
#include "api.h"
static Notification * newNotification(uint8_t notificationType)
{
	Notification * tmp = (Notification *)malloc(sizeof(Notification));
	*tmp = (Notification){0};
	tmp->notificationType = notificationType;
	return tmp;
}

// exports a C++ notification into a copy that can be consumed by the Go Layer. The consumer
// must propagate the responsibility to free the returned object to its callers.
Notification * exportNotification(Manager * manager, OpenZWave::Notification const* notification)
{
	  Notification * result = newNotification(notification->GetType());
	  result->nodeId = (struct NodeId) {notification->GetHomeId(), notification->GetNodeId() };
	  result->notificationCode =
	    notification->GetType() == OpenZWave::Notification::Type_Notification
	    ? notification->GetNotification()
	    : -1;
	  result->valueId = exportValueID(manager, notification->GetValueID());
	  return result;
}


void freeNotification(Notification * notification)
{
	if (notification->node) {
		freeNode(notification->node);
	}
	if (notification->valueId) {
		freeValueID(notification->valueId);
	}
	free(notification);
}
