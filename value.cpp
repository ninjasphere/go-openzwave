#include "api.h"

static Value * newValue(uint32 homeId, OpenZWave::ValueID const &valueId)
{
  Value * tmp = (Value *)malloc(sizeof(Value));
  *tmp = (struct Value){0};
  tmp->homeId = homeId;

  tmp->valueId.id = valueId.GetId();
  tmp->valueId.valueType = valueId.GetType();
  tmp->valueId.commandClassId = valueId.GetCommandClassId();
  tmp->valueId.instance = valueId.GetInstance();
  tmp->valueId.index = valueId.GetIndex();
  
  return tmp;
}

void freeValue(Value * valueObj)
{

  if (valueObj->value) {
    free(valueObj->value);
  }
  if (valueObj->label) {
    free(valueObj->label);
  }
  if (valueObj->units) {
    free(valueObj->units);
  }
  if (valueObj->help) {
    free(valueObj->help);
  }
  free(valueObj);
}

Value * exportValue(API * api, uint32 homeId, OpenZWave::ValueID const &valueId)
{
  Value * const tmp = newValue(homeId, valueId);

  std::string value;

  OpenZWave::Manager * const zwManager = OpenZWave::Manager::Get();

  if (zwManager->GetValueAsString(valueId, &value)) {
    tmp->value = strdup(value.c_str());
  } else {
    tmp->value = strdup("");
  }

  tmp->label = strdup(zwManager->GetValueLabel(valueId).c_str());
  tmp->help = strdup(zwManager->GetValueHelp(valueId).c_str());
  tmp->units = strdup(zwManager->GetValueUnits(valueId).c_str());
  tmp->min = zwManager->GetValueMin(valueId);
  tmp->max = zwManager->GetValueMax(valueId);
  tmp->isSet = zwManager->IsValueSet(valueId);
  
  return tmp;
}

bool setUint8Value(uint32 homeId, uint64 id, uint8 value)
{
	OpenZWave::ValueID valueId = OpenZWave::ValueID(homeId, id);
	OpenZWave::Manager::Get()->SetChangeVerified(valueId, true);
	return OpenZWave::Manager::Get()->SetValue(valueId, value);
}

bool getUint8Value(uint32 homeId, uint64 id, uint8 *value)
{
  return OpenZWave::Manager::Get()->GetValueAsByte(OpenZWave::ValueID(homeId, id), value);
}

bool setBoolValue(uint32 homeId, uint64 id, bool value)
{
  return OpenZWave::Manager::Get()->SetValue(OpenZWave::ValueID(homeId, id), value);
}

bool getBoolValue(uint32 homeId, uint64 id, bool *value)
{
  return OpenZWave::Manager::Get()->GetValueAsBool(OpenZWave::ValueID(homeId, id), value);
}

bool  setFloatValue(uint32 homeId, uint64 id, float value)
{
	return OpenZWave::Manager::Get()->SetValue(OpenZWave::ValueID(homeId, id), value);
}

bool  getFloatValue(uint32 homeId, uint64 id, float *value)
{
	  return OpenZWave::Manager::Get()->GetValueAsFloat(OpenZWave::ValueID(homeId, id), value);
}

bool  setIntValue(uint32 homeId, uint64 id, int value)
{
	return OpenZWave::Manager::Get()->SetValue(OpenZWave::ValueID(homeId, id), value);
}

bool  getIntValue(uint32 homeId, uint64 id, int *value)
{
	  return OpenZWave::Manager::Get()->GetValueAsInt(OpenZWave::ValueID(homeId, id), value);
}

bool  setStringValue(uint32 homeId, uint64 id, char * value)
{
	bool result = OpenZWave::Manager::Get()->SetValue(OpenZWave::ValueID(homeId, id), std::string(value));
	free(value);
	return result;
}

bool  getStringValue(uint32 homeId, uint64 id, char ** value)
{
	  std::string tmp;
	  if (OpenZWave::Manager::Get()->GetValueAsString(OpenZWave::ValueID(homeId, id), &tmp)) {
		*value = strdup(tmp.c_str());
		return true;
	  } else {
		*value = NULL;
		return false;
	  }
}


bool refreshValue(uint32 homeId, uint64 id)
{
  return OpenZWave::Manager::Get()->RefreshValue(OpenZWave::ValueID(homeId, id));
}

bool  setPollingState(uint32 homeId, uint64 id, bool state)
{
  if (state) {
    return OpenZWave::Manager::Get()->EnablePoll(OpenZWave::ValueID(homeId, id));
  } else {
    return OpenZWave::Manager::Get()->DisablePoll(OpenZWave::ValueID(homeId, id));
  }
}
