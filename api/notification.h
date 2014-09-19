typedef struct {
  uint32_t  flags;
  uint8_t   notificationType;
  uint8_t   notificationCode;
  NodeId    nodeId;
  Node      *node;
  Value		*value;
} Notification;

extern Notification * newNotification(uint8_t notificationType);
extern void freeNotification(Notification *);
