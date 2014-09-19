//
// api.cpp
//
// Provides marshalling between C and C++ abstractions. There is typically one C function in this module for each Go function in api.go
//

#include <stdio.h>

#include "Manager.h"
#include "Options.h"
#include "platform/Log.h"
#include "api.h"

extern "C" {
#include "_cgo_export.h"
}

int TRUE = 1;
int FALSE = 0;
int LogLevel_Detail = (int)(OpenZWave::LogLevel_Detail);
int LogLevel_Debug = (int)(OpenZWave::LogLevel_Debug);
int LogLevel_Error = (int)(OpenZWave::LogLevel_Error);
int LogLevel_Info = (int)(OpenZWave::LogLevel_Info);

extern Options createOptions(char * config, char * log)
{
  OpenZWave::Options::Create(config, log, "");
  return (Options*) OpenZWave::Options::Get();
}

extern void addIntOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionInt(option, value);
}

extern void addBoolOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionBool(option, value == TRUE ? true : false);
}

extern Manager lockOptions(Options options)
{
  ((OpenZWave::Options *)options)->Lock();
  return (Manager *)OpenZWave::Manager::Create();
}

extern void addDriver(Manager manager, char * device)
{
  ((OpenZWave::Manager *)manager)->AddDriver(device);
}

extern void addWatcher(Manager manager, void * context) 
{
 ((OpenZWave::Manager *)manager)->AddWatcher( OnNotification, context );
}

void OnNotification (OpenZWave::Notification const* notification, void* context)
{
  Notification * result = newNotification(notification->GetType());
  result->nodeId = (struct NodeId) {notification->GetHomeId(), notification->GetNodeId() };
  OnNotificationWrapper(result, context);
}

extern Notification * newNotification(uint8_t notificationType)
{
	Notification * tmp = (Notification *)malloc(sizeof(Notification));
	*tmp = (Notification){0};
	tmp->notificationType = notificationType;
	return tmp;
}

extern void freeNotification(Notification * notification)
{
	if (notification->node) {
		freeNode(notification->node);
	}
	free(notification);
}

extern Node * newNode(uint32_t homeId, uint8_t nodeId)
{
	Node * tmp = (Node *)malloc(sizeof(Node));
	*tmp = (Node){0};
	tmp->nodeId.homeId = homeId;
	tmp->nodeId.nodeId = nodeId;
	return tmp;
}

extern void freeNode(Node * node)
{
	free(node);
};

