
typedef struct Manager {
	OPENZWAVE_MANAGER * manager;
} Manager;
extern Manager * startManager(API * context);
extern void stopManager(Manager * manager, API * context);
extern bool addDriver(Manager * manager, char * device);
extern bool removeDriver(Manager * manager, char * device);
