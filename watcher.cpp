#include "Notification.h"
#include "watcher.hpp"
namespace OpenZWave {
void OnNotification (Notification const* _notification, void* _context)
{
  ValueID id = _notification->GetValueID();
  switch (_notification->GetType()) {
  case Notification::Type_Notification:
    switch (_notification->GetNotification()) {
    case Notification::Code_MsgComplete:
    case Notification::Code_Timeout:
    case Notification::Code_NoOperation:
    case Notification::Code_Awake:
    case Notification::Code_Sleep:
    case Notification::Code_Dead:
    default:
      break;
    }
    break;
  case Notification::Type_ValueAdded:
  case Notification::Type_ValueRemoved:
  case Notification::Type_ValueChanged:
  case Notification::Type_ValueRefreshed:
  case Notification::Type_Group:
  case Notification::Type_NodeNew:
  case Notification::Type_NodeAdded:
  case Notification::Type_NodeRemoved:
  case Notification::Type_NodeProtocolInfo:
  case Notification::Type_NodeNaming:
  case Notification::Type_NodeEvent:
  case Notification::Type_PollingDisabled:
  case Notification::Type_PollingEnabled:
  case Notification::Type_SceneEvent:
  case Notification::Type_CreateButton:
  case Notification::Type_DeleteButton:
  case Notification::Type_ButtonOn:
  case Notification::Type_ButtonOff:
  case Notification::Type_DriverReady:
  case Notification::Type_DriverFailed:
  case Notification::Type_DriverReset:
  case Notification::Type_EssentialNodeQueriesComplete:
  case Notification::Type_NodeQueriesComplete:
  case Notification::Type_AwakeNodesQueried:
  case Notification::Type_AllNodesQueriedSomeDead:
  case Notification::Type_AllNodesQueried:
  default:
    break;
  }
}
}
