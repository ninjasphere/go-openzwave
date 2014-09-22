typedef struct Notification {
  uint32_t  flags;
  uint8_t   notificationType;
  uint8_t   notificationCode;
  NodeId    nodeId;
  Node      *node; //owned
  ValueID	*valueId; //owned
  Value     *value; // owned
} Notification;

extern void freeNotification(Notification *);

#ifdef __cplusplus
extern Notification * exportNotification(Manager * manager, OpenZWave::Notification const* notification);
#endif
