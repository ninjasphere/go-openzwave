typedef void GoNotification;

typedef struct Notification {
  GoNotification * goRef;
  uint8_t          notificationType;
  uint8_t          notificationCode;
  NodeId           nodeId;
  Node           * node; //owned
  ValueID        * valueId; //owned
  Value          * value; // owned
} Notification;

extern void freeNotification(Notification *);

#ifdef __cplusplus
extern Notification * exportNotification(API * manager, OpenZWave::Notification const* notification);
#endif
