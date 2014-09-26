#include "api.h"

static Value * newValue()
{
  Value * tmp = (Value *)malloc(sizeof(Value));
  *tmp = (struct Value){0};
  tmp->goRef = newGoValue(tmp);
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

Value * exportValue(API * api, OpenZWave::ValueID const &valueId)
{
  Value * const tmp = newValue();

  std::string value;

  OpenZWave::Manager * const zwManager = OpenZWave::Manager::Get();

  if (zwManager->GetValueAsString(valueId, &value)) {
    tmp->value = strdup(value.c_str());
  } else {
    tmp->value = strdup("");
  }


  tmp->valueId.id = valueId.GetId();
  tmp->valueId.valueType = valueId.GetType();
  tmp->valueId.commandClassId = valueId.GetCommandClassId();
  tmp->valueId.instance = valueId.GetInstance();
  tmp->valueId.index = valueId.GetIndex();
  
  tmp->label = strdup(zwManager->GetValueLabel(valueId).c_str());
  tmp->help = strdup(zwManager->GetValueHelp(valueId).c_str());
  tmp->units = strdup(zwManager->GetValueUnits(valueId).c_str());
  tmp->min = zwManager->GetValueMin(valueId);
  tmp->max = zwManager->GetValueMax(valueId);
  tmp->isSet = zwManager->IsValueSet(valueId);
  
  return tmp;
}

bool setUint8Value(uint32_t homeId, uint64_t id, uint8_t value)
{
  return OpenZWave::Manager::Get()->SetValue(OpenZWave::ValueID(homeId, id), value);
}
