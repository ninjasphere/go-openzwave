#include "api.h"

static ValueID * newValueID(uint8_t valueType, uint64_t id)
{
  ValueID * tmp = (ValueID *)malloc(sizeof(ValueID));
  *tmp = (struct ValueID){0};
  tmp->valueType = valueType;
  tmp->id = id;
  return tmp;
}

void freeValueID(ValueID * value)
{
  free(value);
}

ValueID * exportValueID(Manager * manager, OpenZWave::ValueID const & src)
{
	  ValueID * const target = newValueID(src.GetType(), src.GetId());
	  target->commandClassId = src.GetCommandClassId();
	  target->instance = src.GetInstance();
	  target->index = src.GetIndex();
	  return target;
}


static Value * newValue()
{
  Value * tmp = (Value *)malloc(sizeof(Value));
  *tmp = (struct Value){0};
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

Value * exportValue(Manager * manager, OpenZWave::ValueID const &valueId)
{
	Value * const tmp = newValue();
	*tmp = (struct Value){0};

	std::string value;

	if (manager->manager->GetValueAsString(valueId, &value)) {
		tmp->value = strdup(value.c_str());
	} else {
		tmp->value = strdup("");
	}

	OpenZWave::Manager * const zwManager = manager->manager;

	tmp->label = strdup(zwManager->GetValueLabel(valueId).c_str());
	tmp->help = strdup(zwManager->GetValueHelp(valueId).c_str());
	tmp->units = strdup(zwManager->GetValueUnits(valueId).c_str());
	tmp->min = zwManager->GetValueMin(valueId);
	tmp->max = zwManager->GetValueMax(valueId);

	return tmp;
}
