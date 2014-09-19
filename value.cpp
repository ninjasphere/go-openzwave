#include <stdlib.h>
#include "Manager.h"
#include "Value.h"

#include "api.h"

Value * newValue(uint8_t valueType, uint64_t valueId)
{
  Value * tmp = (Value *)malloc(sizeof(Value));
  *tmp = (struct Value){0};
  tmp->valueType = valueType;
  tmp->valueId = valueId;
  return tmp;
}

void freeValue(Value * value)
{
  free(value);
}
