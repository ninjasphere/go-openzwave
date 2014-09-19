#include <stdlib.h>
#include "Manager.h"
#include "Value.h"

#include "api.h"

ValueID * newValueID(uint8_t valueType, uint64_t valueId)
{
  ValueID * tmp = (ValueID *)malloc(sizeof(ValueID));
  *tmp = (struct ValueID){0};
  tmp->valueType = valueType;
  tmp->valueId = valueId;
  return tmp;
}

void freeValueID(ValueID * value)
{
  free(value);
}
