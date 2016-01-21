typedef void GoNode;

typedef struct NodeId {
	uint32 homeId;
	uint8 nodeId;
} NodeId;

typedef struct Node {
	struct NodeId nodeId;
	uint8 basicType;
	uint8 genericType;
	uint8 specificType;
	char * nodeType;
	char * manufacturerName;
	char * productName;
	char * nodeName;
	char * location;
	char * manufacturerId;
	char * productType;
	char * productId;
} Node;

extern void freeNode(Node *);
#ifdef __cplusplus
extern Node * exportNode(API * api, uint32 homeId, uint8 nodeId);
#endif
