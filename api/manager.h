
typedef struct Manager {
	OPENZWAVE_MANAGER * manager;
} Manager;
extern Manager startManager(char * device, void * context);
extern void stopManager(Manager );
