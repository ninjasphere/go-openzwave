typedef struct ValueID {
	  uint64_t  valueId;
	  uint8_t   valueType;
} ValueID;

extern ValueID * newValueID(uint8_t valueType, uint64_t valueId);
extern void freeValueID(ValueID *);

typedef struct Value {
	ValueID valueId;
	char * value;
	char * label;
	char * units;
	char * help;
} Value;

extern Value * newValue();
extern void freeValue(Value *);
