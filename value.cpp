#include "api.h"

ValueID * newValueID(uint8_t valueType, uint64_t id)
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
	  return target;
}


Value * newValue()
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
