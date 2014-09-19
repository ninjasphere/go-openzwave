typedef struct NodeId {
  uint32_t homeId;
  uint8_t  nodeId;
} NodeId;

typedef struct Node {
  struct NodeId nodeId;
} Node;

extern Node * newNode(uint32_t homeId, uint8_t nodeId);
extern void freeNode(Node *);
