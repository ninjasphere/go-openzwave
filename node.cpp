#include "api.h"

static Node * newNode(uint32_t homeId, uint8_t nodeId)
{
  Node * tmp = (Node *)malloc(sizeof(Node));
  *tmp = (Node){0};
  tmp->nodeId.homeId = homeId;
  tmp->nodeId.nodeId = nodeId;
  return tmp;
}

void freeNode(Node * node)
{
  free(node);
};

