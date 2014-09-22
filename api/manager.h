
typedef struct Manager {
	OPENZWAVE_MANAGER * manager;
} Manager;
extern Manager * startManager(void * context);
extern void stopManager(Manager * manager, void * context);
extern bool addDriver(Manager * manager, char * device);
extern bool removeDriver(Manager * manager, char * device);
