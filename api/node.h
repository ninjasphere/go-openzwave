typedef void GoNode;

typedef struct NodeId {
	uint32_t homeId;
	uint8_t nodeId;
} NodeId;

typedef struct Node {
	struct NodeId nodeId;
	uint8_t basicType;
	uint8_t genericType;
	uint8_t specificType;
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
