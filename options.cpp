#include "api.h"

Options startOptions(char * config, char * userPath, char * overrides)
{
  OpenZWave::Options::Create(config, userPath, overrides);
  return (Options*) OpenZWave::Options::Get();
}

void addIntOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionInt(option, value);
}

void addBoolOption(Options options, char * option, bool flag)
{
  ((OpenZWave::Options *)options)->AddOptionBool(option, flag);
}

void endOptions(Options options)
{
  ((OpenZWave::Options *)options)->Lock();
}



