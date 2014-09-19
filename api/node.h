typedef struct NodeId {
  uint32_t homeId;
  uint8_t  nodeId;
} NodeId;

typedef struct Node {
  struct NodeId nodeId;
} Node;

extern void freeNode(Node *);
#ifdef __cplusplus
extern Node * newNode(OpenZWave::Node const &);
#endif
