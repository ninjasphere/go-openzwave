typedef void * Manager;
extern Manager createManager();
extern void addDriver(Manager , char *);
extern void setNotificationWatcher(Manager, void *);
