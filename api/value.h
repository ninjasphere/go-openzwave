typedef struct Value {
	  uint64_t  valueId;
	  uint8_t   valueType;
} Value;

extern Value * newValue(uint8_t valueType, uint64_t valueId);
extern void freeValue(Value *);
