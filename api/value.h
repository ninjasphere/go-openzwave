typedef struct ValueID {
	  uint64_t  id;
	  uint8_t   valueType;
	  uint8_t   commandClassId;
	  uint8_t   instance;
	  uint8_t   index;

} ValueID;

extern void freeValueID(ValueID *);

typedef struct Value {
	ValueID   valueId;
	char    * value;
	char    * label;
	char    * units;
	char    * help;
	int32_t   min;
	int32_t   max;
	bool	  isSet;
} Value;

extern void freeValue(Value *);

#ifdef __cplusplus
extern ValueID * exportValueID(API *, OpenZWave::ValueID const &);
extern Value * exportValue(API *, OpenZWave::ValueID const &);
#endif
