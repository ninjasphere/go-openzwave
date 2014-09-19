typedef struct {
  uint32_t  flags;
  uint8_t   notificationType;
  uint8_t   notificationCode;
  NodeId    nodeId;
  Node      *node;
  uint64_t  valueId;
  uint8_t   valueType;
} Notification;

extern Notification * newNotification(uint8_t notificationType);
extern void freeNotification(Notification *);
