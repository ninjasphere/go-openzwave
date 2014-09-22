#include "api.h"

Manager * newManager(OpenZWave::Manager * cppObj) {
  Manager * manager = (Manager *)malloc((sizeof *manager));
  *manager = (struct Manager){ cppObj };
  return manager;
}

void freeManager(Manager * manager)
{
  OpenZWave::Manager::Destroy();
  free(manager);
}

bool addDriver(Manager * manager, char * device)
{
  return manager->manager->AddDriver(device);
}

bool removeDriver(Manager * manager, char * device)
{
  return manager->manager->RemoveDriver(device);
}
