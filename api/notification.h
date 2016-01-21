typedef void GoNotification;

typedef struct Notification {
  uint8          notificationType;
  uint8          notificationCode;
  Node           * node; //owned
  Value          * value; // owned
} Notification;

extern void freeNotification(Notification *);

#ifdef __cplusplus
extern Notification * exportNotification(API * manager, OpenZWave::Notification const* notification);
#endif
