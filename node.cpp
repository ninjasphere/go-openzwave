#include "api.h"

static Node * newNode(uint32 homeId, uint8 nodeId)
{
  Node * tmp = (Node *)malloc(sizeof(Node));
  *tmp = (Node){0};
  tmp->nodeId.homeId = homeId;
  tmp->nodeId.nodeId = nodeId;
  return tmp;
}

void freeNode(Node * node)
{
  if (node->nodeType) free(node->nodeType);
  if (node->manufacturerName) free(node->manufacturerName);
  if (node->productName) free(node->productName);
  if (node->location) free(node->location);
  if (node->manufacturerId) free(node->manufacturerId);
  if (node->productType) free(node->productType);
  if (node->productId) free(node->productId);
  free(node);
};

Node * exportNode(API * api, uint32 homeId, uint8 nodeId)
{
  Node * result = newNode(homeId, nodeId);
  OpenZWave::Manager * cppRef = OpenZWave::Manager::Get();

  result->basicType = cppRef->GetNodeBasic(homeId, nodeId);
  result->genericType = cppRef->GetNodeGeneric(homeId, nodeId);
  result->specificType = cppRef->GetNodeSpecific(homeId, nodeId);
  result->nodeType = strdup(cppRef->GetNodeType(homeId, nodeId).c_str());
  result->manufacturerName = strdup(cppRef->GetNodeManufacturerName(homeId, nodeId).c_str());
  result->productName = strdup(cppRef->GetNodeProductName(homeId, nodeId).c_str());
  result->nodeName = strdup(cppRef->GetNodeName(homeId, nodeId).c_str());
  result->location = strdup(cppRef->GetNodeLocation(homeId, nodeId).c_str());
  result->manufacturerId = strdup(cppRef->GetNodeManufacturerId(homeId, nodeId).c_str());
  result->productType = strdup(cppRef->GetNodeProductType(homeId, nodeId).c_str());
  result->productId = strdup(cppRef->GetNodeProductId(homeId, nodeId).c_str());
  return result;
}
