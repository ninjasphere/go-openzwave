typedef struct ValueID {
	  uint64_t  valueId;
	  uint8_t   valueType;
} ValueID;

extern void freeValueID(ValueID *);

typedef struct Value {
	ValueID valueId;
	char * value;
	char * label;
	char * units;
	char * help;
} Value;

extern void freeValue(Value *);

#ifdef __cplusplus
extern ValueID * exportValueID(OpenZWave::ValueID const &);
extern Value * newValue(OpenZWave::Value const &);
#endif
