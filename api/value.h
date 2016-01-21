typedef void GoValue;

typedef struct ValueID {
  uint64  id;
  uint8   valueType;
  uint8   commandClassId;
  uint8   instance;
  uint8   index;

} ValueID;

typedef struct Value {
  uint32  homeId;
  ValueID   valueId;
  char    * value;
  char    * label;
  char    * units;
  char    * help;
  int32_t   min;
  int32_t   max;
  bool	    isSet;
} Value;

extern void  freeValue(Value *);
extern bool  setUint8Value(uint32 homeId, uint64 id, uint8 value);
extern bool  getUint8Value(uint32 homeId, uint64 id, uint8 *value);
extern bool  setBoolValue(uint32 homeId, uint64 id, bool value);
extern bool  getBoolValue(uint32 homeId, uint64 id, bool *value);
extern bool  setFloatValue(uint32 homeId, uint64 id, float value);
extern bool  getFloatValue(uint32 homeId, uint64 id, float *value);
extern bool  setIntValue(uint32 homeId, uint64 id, int value);
extern bool  getIntValue(uint32 homeId, uint64 id, int *value);
extern bool  setStringValue(uint32 homeId, uint64 id, char * value);
extern bool  getStringValue(uint32 homeId, uint64 id, char ** value);
extern bool  refreshValue(uint32 homeId, uint64 id);
extern bool  setPollingState(uint32 homeId, uint64 id, bool state);

#ifdef __cplusplus
extern Value * exportValue(API *, uint32 homeId, OpenZWave::ValueID const &);
#endif
