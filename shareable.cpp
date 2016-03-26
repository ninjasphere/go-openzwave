#include "api.h"

shareable * newShareable(int i) {
	shareable * r = (shareable *) malloc(sizeof(shareable));
	r->sharedIndex = i;
	return 	r;
}
