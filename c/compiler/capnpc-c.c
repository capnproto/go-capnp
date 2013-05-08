#include "schema.h"

struct scope {
	struct capn_tree hdr;
	struct scope *enclosing;
	capn_text name;
	uint64_t id;
};

struct scope *find_scope(struct scope *s, uint64_t id) {
	while (s && s->id != id) {
		s = (struct scope*) s->hdr.link[s->id < id];
	}
	return s;
}

struct capn_tree *insert_scope(struct capn_tree *root, struct scope *s) {
	struct capn_tree **x = &root;
	while (*x) {
		s->hdr.parent = *x;
		x = &(*x)->link[((struct scope*)*x)->id < s->id];
	}
	*x = &s->hdr;
	return capn_tree_insert(root, &s->hdr);
}

int main() {
	struct capn capn;
	CodeGeneratorRequest_ptr root;
	struct CodeGeneratorRequest req;
	int i, j;

	if (capn_init_fp(&capn, stdin, 0)) {
		fprintf(stderr, "failed to read schema from stdin\n");
		return -1;
	}

	root.p = capn_get_root(&capn);
	read_CodeGeneratorRequest(&req, root);

	for (i = 0; i < req.nodes.p.size; i++) {
		struct Node N;
		struct FileNode F;
		struct StructNode S;
		struct EnumNode E;
		struct InterfaceNode I;
		struct ConstNode C;
		struct AnnotationNode A;

		get_Node(&N, req.nodes, i);
		fprintf(stderr, "node %s id:%#llx scope:%#llx type:%d\n",
				N.displayName.str, N.id, N.scopeId, N.body_tag);

		switch (N.body_tag) {
		case Node_fileNode:
			read_FileNode(&F, N.body.fileNode);
			for (j = 0; j < F.imports.p.size; j++) {
				struct FileNode_Import fi;
				get_FileNode_Import(&fi, F.imports, j);
				fprintf(stderr, "\timport %#llx %s\n", fi.id, fi.name.str);
			}
			break;
		case Node_structNode:
			read_StructNode(&S, N.body.structNode);
			fprintf(stderr, "\tstruct %d %d %d\n",
					S.dataSectionWordSize, S.pointerSectionSize, S.preferredListEncoding);
			for (j = 0; j < S.members.p.size; j++) {
			}
			break;
		case Node_enumNode:
			read_EnumNode(&E, N.body.enumNode);
			for (j = 0; j < E.enumerants.p.size; j++) {
				struct EnumNode_Enumerant ee;
				get_EnumNode_Enumerant(&ee, E.enumerants, j);
				fprintf(stderr, "\tenum %d %s %d\n", j, ee.name.str, ee.codeOrder);
			}
			break;
		case Node_interfaceNode:
			read_InterfaceNode(&I, N.body.interfaceNode);
			break;
		case Node_constNode:
			read_ConstNode(&C, N.body.constNode);
			break;
		case Node_annotationNode:
			read_AnnotationNode(&A, N.body.annotationNode);
			break;
		}
	}

	return 0;
}
