#include "api.h"

bool addDriver(char * device)
{
  return OpenZWave::Manager::Get()->AddDriver(device);
}

bool removeDriver(char * device)
{
  return OpenZWave::Manager::Get()->RemoveDriver(device);
}
