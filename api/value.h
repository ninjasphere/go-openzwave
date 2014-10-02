typedef void GoValue;

typedef struct ValueID {
  uint64_t  id;
  uint8_t   valueType;
  uint8_t   commandClassId;
  uint8_t   instance;
  uint8_t   index;

} ValueID;

typedef struct Value {
  uint32_t  homeId;
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
extern bool  setUint8Value(uint32_t homeId, uint64_t id, uint8_t value);
extern bool  getUint8Value(uint32_t homeId, uint64_t id, uint8_t *value);
extern bool  setBoolValue(uint32_t homeId, uint64_t id, bool value);
extern bool  getBoolValue(uint32_t homeId, uint64_t id, bool *value);
extern bool  refreshValue(uint32_t homeId, uint64_t id);
extern bool  setPollingState(uint32_t homeId, uint64_t id, bool state);

#ifdef __cplusplus
extern Value * exportValue(API *, uint32_t homeId, OpenZWave::ValueID const &);
#endif
