#ifndef ZWAVE_TYPES_H
#define ZWAVE_TYPES_H

#ifdef __cplusplus
extern "C" {
#endif
// Basic types
typedef signed char                     int8;
typedef unsigned char           uint8;

typedef signed short            int16;
typedef unsigned short          uint16;

typedef signed int                      int32;
typedef unsigned int            uint32;

#ifdef _MSC_VER
typedef signed __int64          int64;
typedef unsigned __int64        uint64;
#endif

#ifdef __GNUC__
typedef signed long long        int64;
typedef unsigned long long  uint64;
#endif
#ifdef __cplusplus
}
#endif
#endif
