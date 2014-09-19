#include "Manager.h"
#include "Options.h"
#include "api.h"

Options createOptions(char * config, char * log)
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

Manager lockOptions(Options options)
{
  ((OpenZWave::Options *)options)->Lock();
  return (Manager *)OpenZWave::Manager::Create();
}

