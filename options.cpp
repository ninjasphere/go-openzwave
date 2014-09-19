#include "Manager.h"
#include "Options.h"
#include "api.h"

Options startOptions(char * config, char * log)
{
  OpenZWave::Options::Create(config, log, "");
  return (Options*) OpenZWave::Options::Get();
}

void addIntOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionInt(option, value);
}

void addBoolOption(Options options, char * option, int value)
{
  ((OpenZWave::Options *)options)->AddOptionBool(option, value == TRUE ? true : false);
}

void endOptions(Options options)
{
  ((OpenZWave::Options *)options)->Lock();
}



