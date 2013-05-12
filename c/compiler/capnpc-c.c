#include "schema.capnp.h"
#include "str.h"
#include <stdlib.h>
#include <string.h>
#include <errno.h>

struct node {
	struct capn_tree hdr;
	struct Node n;
	struct node *next;
	struct node *first_child, *next_child;
	struct str name;
	struct StructNode s;
	struct EnumNode e;
	struct InterfaceNode i;
	struct ConstNode c;
	struct FileNode f;
};

static struct node *g_files;
static FILE *HDR, *SRC;

static struct capn_tree *g_node_tree;

struct node *find_node(uint64_t id) {
	struct node *s = (struct node*) g_node_tree;
	while (s && s->n.id != id) {
		s = (struct node*) s->hdr.link[s->n.id < id];
	}
	if (s == NULL) {
		fprintf(stderr, "cant find node with id %#llx\n", id);
		exit(2);
	}
	return s;
}

static void insert_node(struct node *s) {
	struct capn_tree **x = &g_node_tree;
	while (*x) {
		s->hdr.parent = *x;
		x = &(*x)->link[((struct node*)*x)->n.id < s->n.id];
	}
	*x = &s->hdr;
	g_node_tree = capn_tree_insert(g_node_tree, &s->hdr);
}

static void resolve_names(struct str *b, struct node *n, capn_text name, struct node *file) {
	int i, sz = b->len;
	str_add(b, name.str, name.len);
	str_add(&n->name, b->str, b->len);
	str_add(b, "_", 1);

	for (i = n->n.nestedNodes.p.len-1; i >= 0; i--) {
		struct Node_NestedNode nest;
		get_Node_NestedNode(&nest, n->n.nestedNodes, i);
		resolve_names(b, find_node(nest.id), nest.name, file);
	}

	n->next_child = file->first_child;
	file->first_child = n;
	str_setlen(b, sz);
}

static void define_enum(struct node *n) {
	int i;

	fprintf(HDR, "\nenum %s {", n->name.str);
	for (i = 0; i < n->e.enumerants.p.len; i++) {
		struct EnumNode_Enumerant ee;
		get_EnumNode_Enumerant(&ee, n->e.enumerants, i);
		if (i) {
			fprintf(HDR, ",");
		}
		fprintf(HDR, "\n\t%s_%s = %d", n->name.str, ee.name.str, i);
	}
	fprintf(HDR, "\n};\n");
}

struct member {
	unsigned int is_valid : 1;
	struct StructNode_Member m;
	struct StructNode_Field f;
	struct StructNode_Union u;
	struct Type t;
	struct Value v;
	struct node *type;
	struct Type list;
	struct member *mbrs;
	int idx;
};

static void set_member(FILE *f, struct member *m, const char *tab, const char *var) {
	if (m->t.body_tag == Type_voidType)
		return;

	fputs(tab, f);

	switch (m->t.body_tag) {
	case Type_voidType:
		break;
	case Type_boolType:
		fprintf(f, "err = err || capn_write1(p.p, %d, %s);\n", m->f.offset, var);
		break;
	case Type_int8Type:
		fprintf(f, "err = err || capn_write8(p.p, %d, (uint8_t) %s);\n", m->f.offset, var);
		break;
	case Type_int16Type:
	case Type_enumType:
		fprintf(f, "err = err || capn_write16(p.p, %d, (uint16_t) %s);\n", 2*m->f.offset, var);
		break;
	case Type_int32Type:
		fprintf(f, "err = err || capn_write32(p.p, %d, (uint32_t) %s);\n", 4*m->f.offset, var);
		break;
	case Type_int64Type:
		fprintf(f, "err = err || capn_write64(p.p, %d, (uint64_t) %s);\n", 8*m->f.offset, var);
		break;
	case Type_uint8Type:
		fprintf(f, "err = err || capn_write8(p.p, %d, %s);\n", m->f.offset, var);
		break;
	case Type_uint16Type:
		fprintf(f, "err = err || capn_write16(p.p, %d, %s);\n", 2*m->f.offset, var);
		break;
	case Type_uint32Type:
		fprintf(f, "err = err || capn_write32(p.p, %d, %s);\n", 4*m->f.offset, var);
		break;
	case Type_uint64Type:
		fprintf(f, "err = err || capn_write64(p.p, %d, %s);\n", 8*m->f.offset, var);
		break;
	case Type_float32Type:
		fprintf(f, "err = err || capn_write_float(p.p, %d, %s, 0.0f);\n", 4*m->f.offset, var);
		break;
	case Type_float64Type:
		fprintf(f, "err = err || capn_write_double(p.p, %d, %s, 0.0);\n", 4*m->f.offset, var);
		break;
	case Type_textType:
		fprintf(f, "err = err || capn_set_text(p.p, %d, %s);\n", m->f.offset, var);
		break;
	case Type_dataType:
		fprintf(f, "err = err || capn_set_data(p.p, %d, %s);\n", m->f.offset, var);
		break;
	case Type_structType:
	case Type_interfaceType:
		fprintf(f, "err = err || capn_setp(p.p, %d, %s.p);\n", m->f.offset, var);
		break;
	case Type_objectType:
		fprintf(f, "err = err || capn_setp(p.p, %d, %s);\n", m->f.offset, var);
		break;
	case Type_listType:
		switch (m->list.body_tag) {
		case Type_boolType:
		case Type_int8Type:
		case Type_uint8Type:
		case Type_int16Type:
		case Type_uint16Type:
		case Type_enumType:
		case Type_int32Type:
		case Type_uint32Type:
		case Type_int64Type:
		case Type_uint64Type:
		case Type_float32Type:
		case Type_float64Type:
		case Type_structType:
		case Type_interfaceType:
			fprintf(f, "err = err || capn_setp(p.p, %d, %s.p);\n", m->f.offset, var);
			break;
		case Type_voidType:
		case Type_textType:
		case Type_dataType:
		case Type_objectType:
		case Type_listType:
			fprintf(f, "err = err || capn_setp(p.p, %d, %s);\n", m->f.offset, var);
			break;
		}
	}
}

static void get_member(FILE *f, struct member *m, const char *tab, const char *var) {
	if (m->t.body_tag == Type_voidType)
		return;

	fputs(tab, f);
	fputs(var, f);

	switch (m->t.body_tag) {
	case Type_voidType:
		break;
	case Type_boolType:
		fprintf(f, " = (capn_read8(p.p, %d) & %d) != 0;\n",
				m->f.offset/8, 1 << (m->f.offset%8));
		break;
	case Type_int8Type:
		fprintf(f, " = (int8_t) capn_read8(p.p, %d);\n", m->f.offset);
		break;
	case Type_int16Type:
		fprintf(f, " = (int16_t) capn_read16(p.p, %d);\n", 2*m->f.offset);
		break;
	case Type_int32Type:
		fprintf(f, " = (int32_t) capn_read32(p.p, %d);\n", 4*m->f.offset);
		break;
	case Type_int64Type:
		fprintf(f, " = (int64_t) capn_read64(p.p, %d);\n", 8*m->f.offset);
		break;
	case Type_uint8Type:
		fprintf(f, " = capn_read8(p.p, %d);\n", m->f.offset);
		break;
	case Type_uint16Type:
		fprintf(f, " = capn_read16(p.p, %d);\n", 2*m->f.offset);
		break;
	case Type_uint32Type:
		fprintf(f, " = capn_read32(p.p, %d);\n", 4*m->f.offset);
		break;
	case Type_uint64Type:
		fprintf(f, " = capn_read64(p.p, %d);\n", 8*m->f.offset);
		break;
	case Type_float32Type:
		fprintf(f, " = capn_read_float(p.p, %d, 0.0f);\n", 4*m->f.offset);
		break;
	case Type_float64Type:
		fprintf(f, " = capn_read_double(p.p, %d, 0.0);\n", 8*m->f.offset);
		break;
	case Type_textType:
		fprintf(f, " = capn_get_text(p.p, %d);\n", m->f.offset);
		break;
	case Type_dataType:
		fprintf(f, " = capn_get_data(p.p, %d);\n", m->f.offset);
		break;
	case Type_enumType:
		fprintf(f, " = (enum %s) capn_read16(p.p, %d);\n", m->type->name.str, 2*m->f.offset);
		break;
	case Type_structType:
	case Type_interfaceType:
		fprintf(f, ".p = capn_getp(p.p, %d);\n", m->f.offset);
		break;
	case Type_objectType:
		fprintf(f, " = capn_getp(p.p, %d);\n", m->f.offset);
		break;
	case Type_listType:
		switch (m->list.body_tag) {
		case Type_boolType:
		case Type_int8Type:
		case Type_uint8Type:
		case Type_int16Type:
		case Type_uint16Type:
		case Type_enumType:
		case Type_int32Type:
		case Type_uint32Type:
		case Type_int64Type:
		case Type_uint64Type:
		case Type_float32Type:
		case Type_float64Type:
		case Type_structType:
		case Type_interfaceType:
			fprintf(f, ".p = capn_getp(p.p, %d);\n", m->f.offset);
			break;
		case Type_voidType:
		case Type_textType:
		case Type_dataType:
		case Type_objectType:
		case Type_listType:
			fprintf(f, " = capn_getp(p.p, %d);\n", m->f.offset);
			break;
		}
	}
}

static void union_cases(struct node *n, struct member *m, int set, int mask) {
	static struct str buf = STR_INIT;
	struct member *u = NULL;
	int j;

	for (j = 0; j < m->u.members.p.len; j++) {
		if (mask & (1 << m->mbrs[j].t.body_tag)) {
			u = &m->mbrs[j];
			fprintf(SRC, "\tcase %s_%s:\n", n->name.str, u->m.name.str);
		}
	}

	if (!u)
		return;

	str_reset(&buf);
	str_addf(&buf, "s->%s.%s", m->m.name.str, u->m.name.str);

	if (u->t.body_tag == Type_voidType) {
		/* nothing to do */
	} else if (set) {
		set_member(SRC, u, "\t\t", buf.str);
	} else {
		get_member(SRC, u, "\t\t", buf.str);
	}

	fprintf(SRC, "\t\tbreak;\n");
}

static void do_union(struct node *n, struct member *m, int set) {
	if (set) {
		fprintf(SRC, "\terr = err || capn_write16(p.p, %d, s->%s_tag);\n",
				2*m->u.discriminantOffset, m->m.name.str);
	} else {
		fprintf(SRC, "\ts->%s_tag = (enum %s_%s) capn_read16(p.p, %d);\n",
				m->m.name.str, n->name.str, m->m.name.str, 2*m->u.discriminantOffset);
	}

	fprintf(SRC, "\n\tswitch (s->%s_tag) {\n", m->m.name.str);

	/* if we have a bunch of the same C type with zero defaults, we
	 * only need to emit one switch block as the layout will line up
	 * in the C union */
	union_cases(n, m, set, (1 << Type_voidType));
	union_cases(n, m, set, (1 << Type_boolType));
	union_cases(n, m, set, (1 << Type_int8Type) | (1 << Type_uint8Type));
	union_cases(n, m, set, (1 << Type_int16Type) | (1 << Type_uint16Type) | (1 << Type_enumType));
	union_cases(n, m, set, (1 << Type_int32Type) | (1 << Type_uint32Type) | (1 << Type_float32Type));
	union_cases(n, m, set, (1 << Type_int64Type) | (1 << Type_uint64Type) | (1 << Type_float64Type));
	union_cases(n, m, set, (1 << Type_textType));
	union_cases(n, m, set, (1 << Type_dataType));
	union_cases(n, m, set, (1 << Type_structType) | (1 << Type_interfaceType) | (1 << Type_objectType) | (1 << Type_listType));

	fprintf(SRC, "\t}\n");
}

static void print_member(FILE *f, struct member *m, const char *tab) {
	if (m->t.body_tag == Type_voidType)
		return;

	fputs(tab, f);

	switch (m->t.body_tag) {
	case Type_voidType:
		break;
	case Type_boolType:
		fprintf(f, "unsigned int %s : 1;\n", m->m.name.str);
		break;
	case Type_int8Type:
		fprintf(f, "int8_t %s;\n", m->m.name.str);
		break;
	case Type_int16Type:
		fprintf(f, "int16_t %s;\n", m->m.name.str);
		break;
	case Type_int32Type:
		fprintf(f, "int32_t %s;\n", m->m.name.str);
		break;
	case Type_int64Type:
		fprintf(f, "int64_t %s;\n", m->m.name.str);
		break;
	case Type_uint8Type:
		fprintf(f, "uint8_t %s;\n", m->m.name.str);
		break;
	case Type_uint16Type:
		fprintf(f, "uint16_t %s;\n", m->m.name.str);
		break;
	case Type_uint32Type:
		fprintf(f, "uint32_t %s;\n", m->m.name.str);
		break;
	case Type_uint64Type:
		fprintf(f, "uint64_t %s;\n", m->m.name.str);
		break;
	case Type_float32Type:
		fprintf(f, "float %s;\n", m->m.name.str);
		break;
	case Type_float64Type:
		fprintf(f, "double %s;\n", m->m.name.str);
		break;
	case Type_textType:
		fprintf(f, "capn_text %s;\n", m->m.name.str);
		break;
	case Type_dataType:
		fprintf(f, "capn_data %s;\n", m->m.name.str);
		break;
	case Type_enumType:
		fprintf(f, "enum %s %s;\n", m->type->name.str, m->m.name.str);
		break;
	case Type_structType:
	case Type_interfaceType:
		fprintf(f, "%s_ptr %s;\n", m->type->name.str, m->m.name.str);
		break;
	case Type_objectType:
		fprintf(f, "capn_ptr %s;\n", m->m.name.str);
		break;
	case Type_listType:
		switch (m->list.body_tag) {
		case Type_voidType:
			fprintf(f, "capn_ptr %s;\n", m->m.name.str);
			break;
		case Type_boolType:
			fprintf(f, "capn_list1 %s;\n", m->m.name.str);
			break;
		case Type_int8Type:
		case Type_uint8Type:
			fprintf(f, "capn_list8 %s;\n", m->m.name.str);
			break;
		case Type_int16Type:
		case Type_uint16Type:
		case Type_enumType:
			fprintf(f, "capn_list16 %s;\n", m->m.name.str);
			break;
		case Type_int32Type:
		case Type_uint32Type:
			fprintf(f, "capn_list32 %s;\n", m->m.name.str);
			break;
		case Type_int64Type:
		case Type_uint64Type:
			fprintf(f, "capn_list64 %s;\n", m->m.name.str);
			break;
		case Type_float32Type:
			fprintf(f, "capn_list_float %s;\n", m->m.name.str);
			break;
		case Type_float64Type:
			fprintf(f, "capn_list_double %s;\n", m->m.name.str);
			break;
		case Type_textType:
		case Type_dataType:
		case Type_objectType:
		case Type_listType:
			fprintf(f, "capn_ptr %s;\n", m->m.name.str);
			break;
		case Type_structType:
		case Type_interfaceType:
			fprintf(f, "%s_list %s;\n", m->type->name.str, m->m.name.str);
			break;
		}
	}
}

static struct member *decode_member(struct member *mbrs, StructNode_Member_list l, int i) {
	struct member m;
	m.is_valid = 1;
	m.idx = i;
	get_StructNode_Member(&m.m, l, i);

	if (m.m.codeOrder >= l.p.len) {
		fprintf(stderr, "unexpectedly large code order %d >= %d\n", m.m.codeOrder, l.p.len);
		exit(3);
	}

	if (m.m.body_tag == StructNode_Member_fieldMember) {
		read_StructNode_Field(&m.f, m.m.body.fieldMember);
		read_Type(&m.t, m.f.type);
		read_Value(&m.v, m.f.defaultValue);

		switch (m.t.body_tag) {
		case Type_enumType:
		case Type_structType:
		case Type_interfaceType:
			m.type = find_node(m.t.body.enumType);
			break;

		case Type_listType:
			read_Type(&m.list, m.t.body.listType);

			switch (m.list.body_tag) {
			case Type_enumType:
			case Type_structType:
			case Type_interfaceType:
				m.type = find_node(m.list.body.enumType);
				break;
			default:
				break;
			}

		default:
			break;
		}
	}

	memcpy(&mbrs[m.m.codeOrder], &m, sizeof(m));
	return &mbrs[m.m.codeOrder];
}

static void define_struct(struct node *n) {
	static struct str buf = STR_INIT;

	int i, j, mlen = n->s.members.p.len;
	struct member *mbrs = calloc(mlen, sizeof(*mbrs));

	/* get list of members in code order and emit union enums */
	for (i = 0; i < mlen; i++) {
		struct member *m = decode_member(mbrs, n->s.members, i);

		if (m->m.body_tag == StructNode_Member_unionMember) {
			int first, ulen;

			/* get union members in code order */
			read_StructNode_Union(&m->u, m->m.body.unionMember);
			ulen = m->u.members.p.len;
			m->mbrs = calloc(ulen, sizeof(*m->mbrs));
			for (j = 0; j < ulen; j++) {
				decode_member(m->mbrs, m->u.members, j);
			}

			/* emit union enum definition */
			first = 1;
			fprintf(HDR, "\nenum %s_%s {", n->name.str, m->m.name.str);
			for (j = 0; j < ulen; j++) {
				struct member *u = &m->mbrs[j];
				if (!u->is_valid) continue;
				if (!first) fprintf(HDR, ",");
				fprintf(HDR, "\n\t%s_%s = %d", n->name.str, u->m.name.str, u->idx);
				first = 0;
			}
			fprintf(HDR, "\n};\n");
		}
	}

	/* emit struct definition */
	fprintf(HDR, "\nstruct %s {\n", n->name.str);
	for (i = 0; i < mlen; i++) {
		struct member *m = &mbrs[i];
		if (!m->is_valid)
			continue;

		switch (m->m.body_tag) {
		case StructNode_Member_fieldMember:
			print_member(HDR, m, "\t");
			break;
		case StructNode_Member_unionMember:
			fprintf(HDR, "\tenum %s_%s %s_tag;\n", n->name.str, m->m.name.str, m->m.name.str);
			fprintf(HDR, "\tunion {\n");
			for (j = 0; j < m->u.members.p.len; j++) {
				print_member(HDR, &m->mbrs[j], "\t\t");
			}
			fprintf(HDR, "\t} %s;\n", m->m.name.str);
			break;
		}
	}
	fprintf(HDR, "};\n");

	fprintf(SRC, "\n%s_ptr new_%s(struct capn_segment *s) {\n", n->name.str, n->name.str);
	fprintf(SRC, "\t%s_ptr p = {capn_new_struct(s, %d, %d)};\n",
			n->name.str, 8*n->s.dataSectionWordSize, n->s.pointerSectionSize);
	fprintf(SRC, "\treturn p;\n");
	fprintf(SRC, "}\n");

	fprintf(SRC, "%s_list new_%s_list(struct capn_segment *s, int len) {\n", n->name.str, n->name.str);
	fprintf(SRC, "\t%s_list p = {capn_new_list(s, len, %d, %d)};\n",
			n->name.str, 8*n->s.dataSectionWordSize, n->s.pointerSectionSize);
	fprintf(SRC, "\treturn p;\n");
	fprintf(SRC, "}\n");

	fprintf(SRC, "void read_%s(struct %s *s, %s_ptr p) {\n", n->name.str, n->name.str, n->name.str);
	for (i = 0; i < mlen; i++) {
		struct member *m = &mbrs[i];
		if (!m->is_valid) continue;

		switch (m->m.body_tag) {
		case StructNode_Member_fieldMember:
			str_reset(&buf);
			str_addf(&buf, "s->%s", m->m.name.str);
			get_member(SRC, m, "\t", buf.str);
			break;
		case StructNode_Member_unionMember:
			do_union(n, m, 0);
			break;
		}
	}
	fprintf(SRC, "}\n");

	fprintf(SRC, "int write_%s(const struct %s *s, %s_ptr p) {\n", n->name.str, n->name.str, n->name.str);
	fprintf(SRC, "\tint err = 0;\n");
	for (i = 0; i < mlen; i++) {
		struct member *m = &mbrs[i];
		if (!m->is_valid) continue;

		switch (m->m.body_tag) {
		case StructNode_Member_fieldMember:
			if (m->t.body_tag != Type_voidType) {
				str_reset(&buf);
				str_addf(&buf, "s->%s", m->m.name.str);
				set_member(SRC, m, "\t", buf.str);
			}
			break;
		case StructNode_Member_unionMember:
			do_union(n, m, 1);
			break;
		}
	}
	fprintf(SRC, "\treturn err;\n}\n");

	fprintf(SRC, "void get_%s(struct %s *s, %s_list l, int i) {\n", n->name.str, n->name.str, n->name.str);
	fprintf(SRC, "\t%s_ptr p = {capn_getp(l.p, i)};\n", n->name.str);
	fprintf(SRC, "\tread_%s(s, p);\n", n->name.str);
	fprintf(SRC, "}\n");

	fprintf(SRC, "int set_%s(const struct %s *s, %s_list l, int i) {\n", n->name.str, n->name.str, n->name.str);
	fprintf(SRC, "\t%s_ptr p = {capn_getp(l.p, i)};\n", n->name.str);
	fprintf(SRC, "\treturn write_%s(s, p);\n", n->name.str);
	fprintf(SRC, "}\n");
}

static void declare_structs(struct node *n, const char *format, int num) {
	fprintf(HDR, "\n");
	for (n = n->first_child; n != NULL; n = n->next_child) {
		if (n->n.body_tag == Node_structNode) {
			switch (num) {
			case 3:
				fprintf(HDR, format, n->name.str, n->name.str, n->name.str);
				break;
			case 2:
				fprintf(HDR, format, n->name.str, n->name.str);
				break;
			case 1:
				fprintf(HDR, format, n->name.str);
				break;
			}
		}
	}
}
int main() {
	struct capn capn;
	CodeGeneratorRequest_ptr root;
	struct CodeGeneratorRequest req;
	struct node *n;
	int i;

	if (capn_init_fp(&capn, stdin, 0)) {
		fprintf(stderr, "failed to read schema from stdin\n");
		return -1;
	}

	root.p = capn_getp(capn_root(&capn), 0);
	read_CodeGeneratorRequest(&req, root);

	for (i = 0; i < req.nodes.p.len; i++) {
		n = calloc(1, sizeof(*n));
		get_Node(&n->n, req.nodes, i);
		insert_node(n);

		switch (n->n.body_tag) {
		case Node_fileNode:
			n->next = g_files;
			g_files = n;
			read_FileNode(&n->f, n->n.body.fileNode);
			break;
		case Node_structNode:
			read_StructNode(&n->s, n->n.body.structNode);
			break;
		case Node_enumNode:
			read_EnumNode(&n->e, n->n.body.enumNode);
			break;
		case Node_interfaceNode:
			read_InterfaceNode(&n->i, n->n.body.interfaceNode);
			break;
		case Node_constNode:
			read_ConstNode(&n->c, n->n.body.constNode);
			break;
		default:
			break;
		}
	}

	for (n = g_files; n != NULL; n = n->next) {
		struct str b = STR_INIT;

		for (i = n->n.nestedNodes.p.len-1; i >= 0; i--) {
			struct Node_NestedNode nest;
			get_Node_NestedNode(&nest, n->n.nestedNodes, i);
			resolve_names(&b, find_node(nest.id), nest.name, n);
		}

		str_release(&b);
	}

	for (i = 0; i < req.requestedFiles.p.len; i++) {
		struct str b = STR_INIT;
		struct node *s;
		char *p;
		n = find_node(capn_get64(req.requestedFiles, i));

		str_reset(&b);
		str_addf(&b, "%s.h", n->n.displayName.str);
		HDR = fopen(b.str, "w");
		if (!HDR) {
			fprintf(stderr, "failed to open %s: %s\n", b.str, strerror(errno));
			exit(2);
		}

		str_reset(&b);
		str_addf(&b, "%s.c", n->n.displayName.str);
		SRC = fopen(b.str, "w");
		if (!SRC) {
			fprintf(stderr, "failed to open %s: %s\n", b.str, strerror(errno));
			exit(2);
		}

		fprintf(HDR, "#ifndef CAPN_%llx\n", n->n.id);
		fprintf(HDR, "#define CAPN_%llx\n", n->n.id);
		fprintf(HDR, "/* AUTO GENERATED DO NOT EDIT*/\n");
		fprintf(HDR, "#include <capn.h>\n");

		for (i = 0; i < n->f.imports.p.len; i++) {
			struct FileNode_Import im;
			get_FileNode_Import(&im, n->f.imports, i);
			fprintf(HDR, "#include \"%s.h\"\n", im.name.str);
		}

		fprintf(HDR, "\n#ifdef __cplusplus\nextern \"C\" {\n#endif\n");

		p = strrchr(n->n.displayName.str, '/');
		fprintf(SRC, "#include \"%s.h\"\n", p ? p+1 : n->n.displayName.str);
		fprintf(SRC, "/* AUTO GENERATED DO NOT EDIT*/\n");

		declare_structs(n, "struct %s;\n", 1);
		declare_structs(n, "typedef struct {capn_ptr p;} %s_ptr;\n", 1);
		declare_structs(n, "typedef struct {capn_ptr p;} %s_list;\n", 1);
		declare_structs(n, "%s_ptr new_%s(struct capn_segment*);\n", 2);
		declare_structs(n, "%s_list new_%s_list(struct capn_segment*, int len);\n", 2);
		declare_structs(n, "void read_%s(struct %s*, %s_ptr);\n", 3);
		declare_structs(n, "int write_%s(const struct %s*, %s_ptr);\n", 3);
		declare_structs(n, "void get_%s(struct %s*, %s_list, int i);\n", 3);
		declare_structs(n, "int set_%s(const struct %s*, %s_list, int i);\n", 3);

		for (s = n->first_child; s != NULL; s = s->next_child) {
			switch (s->n.body_tag) {
			case Node_structNode:
				define_struct(s);
				break;
			case Node_enumNode:
				define_enum(s);
				break;
			default:
				break;
			}
		}

		fprintf(HDR, "\n#ifdef __cplusplus\n}\n#endif\n#endif\n");
		str_release(&b);
		fclose(HDR);
		fclose(SRC);
		HDR = SRC = NULL;
	}

	return 0;
}
