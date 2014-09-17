#include <stdio.h>

#include "Manager.h"
#include "wrapper.hpp"

int init(char * device)
{
	OpenZWave::Manager::Create();
	OpenZWave::Manager::Get()->AddDriver(device);
	return 0;
}
