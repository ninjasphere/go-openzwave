#include "Manager.h"
#include "Notification.h"

#include "api.h"

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
	free(notification);
}
