Manager * startManager(API * api)
{
  Manager * manager = newManager(OpenZWave::Manager::Create());
  manager->manager->AddWatcher( OnNotification, api );
  return manager;
}

void stopManager(API * api)
{
  Manager * manager = asManager(api);
  manager->manager->RemoveWatcher(OnNotification, api);
  freeManager(manager);
}

