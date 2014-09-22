typedef struct Manager {
	OPENZWAVE_MANAGER * manager;
} Manager;
extern Manager * newManager(OPENZWAVE_MANAGER * cppObj);
extern void freeManager(Manager * manager);
extern bool addDriver(Manager * manager, char * device);
extern bool removeDriver(Manager * manager, char * device);
