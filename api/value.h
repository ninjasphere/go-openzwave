typedef struct ValueID {
	  uint64_t  valueId;
	  uint8_t   valueType;
} ValueID;

extern ValueID * newValueID(uint8_t valueType, uint64_t valueId);
extern void freeValueID(ValueID *);
