#include "schema.h"


int main() {
	struct capn capn;
	struct CodeGeneratorRequest_ptr root;
	struct CodeGeneratorRequest req;
	int i;

	if (capn_init_fp(&capn, stdin)) {
		fprintf(stderr, "failed to read schema on input\n");
		return -1;
	}

	root.p = capn_root(&capn);
	read_CodeGeneratorRequest(&root, &req);

	for (i = 0; i < req.nodes.size; i++) {
		struct Node_ptr p;
		struct Node n;
		p.p = capn_getp(&req.nodes, i);
		read_Node(&p, &n);

		fprintf(stderr, "node %s id:%#llx scope:%#llx type:%d\n",
				n.displayName.str, n.id, n.scopeId, n.body_tag);
	}

	return 0;
}
