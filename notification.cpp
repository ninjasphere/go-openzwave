
#include "api.h"
// exports a C++ notification into a copy that can be consumed by the Go Layer. The consumer
// must propagate the responsibility to free the returned object to its callers.
Notification * exportNotification(OpenZWave::Notification const* notification)
{
	  Notification * result = newNotification(notification->GetType());
	  result->nodeId = (struct NodeId) {notification->GetHomeId(), notification->GetNodeId() };
	  result->notificationCode =
	    notification->GetType() == OpenZWave::Notification::Type_Notification
	    ? notification->GetNotification()
	    : -1;
	  OpenZWave::ValueID const & valueId = notification->GetValueID();
	  result->valueId = newValueID(valueId.GetType(), valueId.GetId());
	  return result;
}


Notification * newNotification(uint8_t notificationType)
{
	Notification * tmp = (Notification *)malloc(sizeof(Notification));
	*tmp = (Notification){0};
	tmp->notificationType = notificationType;
	return tmp;
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
