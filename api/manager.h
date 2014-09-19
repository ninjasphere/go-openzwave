typedef void * Manager;
extern Manager createManager();
extern void addDriver(Manager , char *);
extern void addWatcher(Manager, void *);
