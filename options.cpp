#include "api.h"

void startOptions(char * config, char * userPath, char * overrides)
{
  OpenZWave::Options::Create(config, userPath, overrides);
}

void addIntOption(char * option, int value)
{
  OpenZWave::Options::Get()->AddOptionInt(option, value);
}

void addBoolOption(char * option, bool flag)
{
  OpenZWave::Options::Get()->AddOptionBool(option, flag);
}

void endOptions()
{
  OpenZWave::Options::Get()->Lock();
}



