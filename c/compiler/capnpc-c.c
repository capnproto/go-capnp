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
static FILE *HDR;
static struct str SRC;
static struct capn g_valcapn;
static struct capn_segment g_valseg;
static int g_valc;
static int g_val0used, g_nullused;

static struct capn_tree *g_node_tree;

struct node *find_node(uint64_t id) {
	struct node *s = (struct node*) g_node_tree;
	while (s && s->n.id != id) {
		s = (struct node*) s->hdr.link[s->n.id < id];
	}
	if (s == NULL) {
		fprintf(stderr, "cant find node with id 0x%x%x\n", (uint32_t) (id >> 32), (uint32_t) id);
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

struct value {
	struct Type t;
	const char *tname;
	struct str buf;
	struct Value v;
	capn_ptr ptr;
	int64_t num;
};

static void decode_value(struct value* v, Type_ptr type, Value_ptr value, const char *symbol) {
	struct Type lt;
	memset(v, 0, sizeof(*v));
	read_Type(&v->t, type);
	read_Value(&v->v, value);

	switch (v->t.body_tag) {
	case Type_voidType:
		v->tname = "void";
		break;
	case Type_boolType:
		v->tname = "unsigned";
		break;
	case Type_int8Type:
		v->tname = "int8_t";
		break;
	case Type_int16Type:
		v->tname = "int16_t";
		break;
	case Type_int32Type:
		v->tname = "int32_t";
		break;
	case Type_int64Type:
		v->tname = "int64_t";
		break;
	case Type_uint8Type:
		v->tname = "uint8_t";
		break;
	case Type_uint16Type:
		v->tname = "uint16_t";
		break;
	case Type_uint32Type:
		v->tname = "uint32_t";
		break;
	case Type_uint64Type:
		v->tname = "uint64_t";
		break;
	case Type_float32Type:
		v->tname = "float";
		break;
	case Type_float64Type:
		v->tname = "double";
		break;
	case Type_textType:
		v->tname = "capn_text";
		break;
	case Type_dataType:
		v->tname = "capn_data";
		break;
	case Type_enumType:
		v->tname = strf(&v->buf, "enum %s", find_node(v->t.body.enumType)->name.str);
		break;
	case Type_structType:
	case Type_interfaceType:
		v->tname = strf(&v->buf, "%s_ptr", find_node(v->t.body.structType)->name.str);
		break;
	case Type_objectType:
		v->tname = "capn_ptr";
		break;
	case Type_listType:
		read_Type(&lt, v->t.body.listType);

		switch (lt.body_tag) {
		case Type_voidType:
			v->tname = "capn_ptr";
			break;
		case Type_boolType:
			v->tname = "capn_list1";
			break;
		case Type_int8Type:
		case Type_uint8Type:
			v->tname = "capn_list8";
			break;
		case Type_int16Type:
		case Type_uint16Type:
		case Type_enumType:
			v->tname = "capn_list16";
			break;
		case Type_int32Type:
		case Type_uint32Type:
		case Type_float32Type:
			v->tname = "capn_list32";
			break;
		case Type_int64Type:
		case Type_uint64Type:
		case Type_float64Type:
			v->tname = "capn_list64";
			break;
		case Type_textType:
		case Type_dataType:
		case Type_objectType:
		case Type_listType:
			v->tname = "capn_ptr";
			break;
		case Type_structType:
		case Type_interfaceType:
			v->tname = strf(&v->buf, "%s_list", find_node(lt.body.structType)->name.str);
			break;
		}
	}

	switch (v->v.body_tag) {
	case Value_boolValue:
		v->num = v->v.body.boolValue;
		break;
	case Value_int8Value:
		v->num = v->v.body.int8Value;
		break;
	case Value_uint8Value:
		v->num = v->v.body.uint8Value;
		break;
	case Value_int16Value:
		v->num = v->v.body.int16Value;
		break;
	case Value_uint16Value:
	case Value_enumValue:
		v->num = v->v.body.uint16Value;
		break;
	case Value_int32Value:
		v->num = v->v.body.int32Value;
		break;
	case Value_uint32Value:
	case Value_float32Value:
		v->num = v->v.body.uint32Value;
		break;
	case Value_int64Value:
		v->num = v->v.body.int64Value;
		break;
	case Value_float64Value:
	case Value_uint64Value:
		v->num = v->v.body.uint64Value;
		break;
	case Value_textValue:
		if (v->v.body.textValue.len) {
			const char *scope = "";
			capn_ptr p = capn_root(&g_valcapn);
			if (capn_set_text(p, 0, v->v.body.textValue)) {
				fprintf(stderr, "failed to copy text\n");
				exit(2);
			}
			p = capn_getp(p, 0);
			if (!p.type)
				break;

			v->ptr = p;

			if (!symbol) {
				static struct str buf = STR_INIT;
				v->num = ++g_valc;
				symbol = strf(&buf, "capn_val%d", (int) v->num);
				scope = "static ";
			} else {
				v->num = 1;
			}

			str_addf(&SRC, "%scapn_text %s = {%d,(char*)&capn_buf[%d],(struct capn_segment*)&capn_seg};\n",
					scope, symbol, p.len-1, (int) (p.data-p.seg->data-8));
		}
		break;

	case Value_dataValue:
	case Value_structValue:
	case Value_objectValue:
	case Value_listValue:
		if (v->v.body.objectValue.type) {
			const char *scope = "";
			capn_ptr p = capn_root(&g_valcapn);
			if (capn_setp(p, 0, v->v.body.objectValue)) {
				fprintf(stderr, "failed to copy object\n");
				exit(2);
			}
			p = capn_getp(p, 0);
			if (!p.type)
				break;

			v->ptr = p;

			if (!symbol) {
				static struct str buf = STR_INIT;
				v->num = ++g_valc;
				symbol = strf(&buf, "capn_val%d", (int) v->num);
				scope = "static ";
			} else {
				v->num = 1;
			}

			str_addf(&SRC, "%s%s %s = {", scope, v->tname, symbol);
			if (strcmp(v->tname, "capn_ptr"))
				str_addf(&SRC, "{");

			str_addf(&SRC, "%d,%d,%d,%d,%d,(char*)&capn_buf[%d],(struct capn_segment*)&capn_seg",
					p.type, p.has_ptr_tag,
					p.datasz, p.ptrsz,
					p.len, (int) (p.data-p.seg->data-8));

			if (strcmp(v->tname, "capn_ptr"))
				str_addf(&SRC, "}");

			str_addf(&SRC, "};\n");
		}
		break;

	case Value_interfaceValue:
	case Value_voidValue:
		break;
	}
}

static void define_const(struct node *n) {
	struct value v;
	decode_value(&v, n->c.type, n->c.value, n->name.str);

	switch (v.v.body_tag) {
	case Value_boolValue:
	case Value_int8Value:
	case Value_int16Value:
	case Value_int32Value:
		fprintf(HDR, "extern %s %s;\n", v.tname, n->name.str);
		str_addf(&SRC, "%s %s = %d;\n", v.tname, n->name.str, (int) v.num);
		break;

	case Value_uint8Value:
	case Value_uint16Value:
	case Value_uint32Value:
		fprintf(HDR, "extern %s %s;\n", v.tname, n->name.str);
		str_addf(&SRC, "%s %s = %uu;\n", v.tname, n->name.str, (uint32_t) v.num);
		break;

	case Value_enumValue:
		fprintf(HDR, "extern %s %s;\n", v.tname, n->name.str);
		str_addf(&SRC, "%s %s = (%s) %uu;\n", v.tname, n->name.str, v.tname, (uint32_t) v.num);
		break;

	case Value_int64Value:
	case Value_uint64Value:
		fprintf(HDR, "extern %s %s;\n", v.tname, n->name.str);
		str_addf(&SRC, "%s %s = ((uint64_t) %#xu << 32) | %#xu;\n", v.tname, n->name.str,
				(uint32_t) (v.num >> 32), (uint32_t) v.num);
		break;

	case Value_float32Value:
		fprintf(HDR, "extern union capn_conv_f32 %s;\n", n->name.str);
		str_addf(&SRC, "union capn_conv_f32 %s = {%#xu};\n", n->name.str, (uint32_t) v.num);
		break;

	case Value_float64Value:
		fprintf(HDR, "extern union capn_conv_f64 %s;\n", n->name.str);
		str_addf(&SRC, "union capn_conv_f64 %s = {((uint64_t) %#xu << 32) | %#xu};\n",
				n->name.str, (uint32_t) (v.num >> 32), (uint32_t) v.num);
		break;

	case Value_textValue:
	case Value_dataValue:
	case Value_structValue:
	case Value_objectValue:
	case Value_listValue:
		fprintf(HDR, "extern %s %s;\n", v.tname, n->name.str);
		if (!v.num) {
			str_addf(&SRC, "%s %s;\n", v.tname, n->name.str);
		}
		break;

	case Value_interfaceValue:
	case Value_voidValue:
		break;
	}

	str_release(&v.buf);
}

struct member {
	unsigned int is_valid : 1;
	struct StructNode_Member m;
	struct StructNode_Field f;
	struct StructNode_Union u;
	struct value v;
	struct member *mbrs;
	int idx;
};

static struct member *decode_member(struct member *mbrs, StructNode_Member_list l, int i) {
	struct member m;
	memset(&m, 0, sizeof(m));
	m.is_valid = 1;
	m.idx = i;
	get_StructNode_Member(&m.m, l, i);

	if (m.m.codeOrder >= l.p.len) {
		fprintf(stderr, "unexpectedly large code order %d >= %d\n", m.m.codeOrder, l.p.len);
		exit(3);
	}

	if (m.m.body_tag == StructNode_Member_fieldMember) {
		read_StructNode_Field(&m.f, m.m.body.fieldMember);
		decode_value(&m.v, m.f.type, m.f.defaultValue, NULL);
	}

	memcpy(&mbrs[m.m.codeOrder], &m, sizeof(m));
	return &mbrs[m.m.codeOrder];
}

static const char *xor_member(struct member *m) {
	static struct str buf = STR_INIT;

	if (m->v.num) {
		switch (m->v.v.body_tag) {
		case Value_int8Value:
		case Value_int16Value:
		case Value_int32Value:
			return strf(&buf, " ^ %d", (int32_t) m->v.num);

		case Value_uint8Value:
		case Value_uint16Value:
		case Value_uint32Value:
		case Value_enumValue:
			return strf(&buf, " ^ %uu", (uint32_t) m->v.num);

		case Value_float32Value:
			return strf(&buf, " ^ %#xu", (uint32_t) m->v.num);

		case Value_int64Value:
		case Value_uint64Value:
		case Value_float64Value:
			return strf(&buf, " ^ ((uint64_t) %#xu << 32) ^ %#xu",
					(uint32_t) (m->v.num >> 32), (uint32_t) m->v.num);

		default:
			return "";
		}
	} else {
		return "";
	}
}

static const char *ptr_member(struct member *m, const char *var) {
	static struct str buf = STR_INIT;
	if (!strcmp(m->v.tname, "capn_ptr")) {
		return var;
	} else if (var[0] == '*') {
		return strf(&buf, "%s->p", var+1);
	} else {
		return strf(&buf, "%s.p", var);
	}
}

static void set_member(struct member *m, const char *ptr, const char *tab, const char *var) {
	const char *xor = xor_member(m);
	const char *pvar = ptr_member(m, var);

	if (m->v.t.body_tag == Type_voidType)
		return;

	str_add(&SRC, tab, -1);

	switch (m->v.t.body_tag) {
	case Type_voidType:
		break;
	case Type_boolType:
		str_addf(&SRC, "err = err || capn_write1(%s, %d, %s != %d);\n", ptr, m->f.offset, var, (int) m->v.num);
		break;
	case Type_int8Type:
		str_addf(&SRC, "err = err || capn_write8(%s, %d, (uint8_t) %s%s);\n", ptr, m->f.offset, var, xor);
		break;
	case Type_int16Type:
	case Type_enumType:
		str_addf(&SRC, "err = err || capn_write16(%s, %d, (uint16_t) %s%s);\n", ptr, 2*m->f.offset, var, xor);
		break;
	case Type_int32Type:
		str_addf(&SRC, "err = err || capn_write32(%s, %d, (uint32_t) %s%s);\n", ptr, 4*m->f.offset, var, xor);
		break;
	case Type_int64Type:
		str_addf(&SRC, "err = err || capn_write64(%s, %d, (uint64_t) %s%s);\n", ptr, 8*m->f.offset, var, xor);
		break;
	case Type_uint8Type:
		str_addf(&SRC, "err = err || capn_write8(%s, %d, %s%s);\n", ptr, m->f.offset, var, xor);
		break;
	case Type_uint16Type:
		str_addf(&SRC, "err = err || capn_write16(%s, %d, %s%s);\n", ptr, 2*m->f.offset, var, xor);
		break;
	case Type_uint32Type:
		str_addf(&SRC, "err = err || capn_write32(%s, %d, %s%s);\n", ptr, 4*m->f.offset, var, xor);
		break;
	case Type_float32Type:
		str_addf(&SRC, "err = err || capn_write32(%s, %d, capn_from_f32(%s)%s);\n", ptr, 4*m->f.offset, var, xor);
		break;
	case Type_uint64Type:
		str_addf(&SRC, "err = err || capn_write64(%s, %d, %s%s);\n", ptr, 8*m->f.offset, var, xor);
		break;
	case Type_float64Type:
		str_addf(&SRC, "err = err || capn_write64(%s, %d, capn_from_f64(%s)%s);\n", ptr, 8*m->f.offset, var, xor);
		break;
	case Type_textType:
		if (m->v.num) {
			g_val0used = 1;
			str_addf(&SRC, "err = err || capn_set_text(%s, %d, (%s.str != capn_val%d.str) ? %s : capn_val0);\n",
					ptr, m->f.offset, var, (int)m->v.num, var);
		} else {
			str_addf(&SRC, "err = err || capn_set_text(%s, %d, %s);\n",
					ptr, m->f.offset, var);
		}
		break;
	case Type_dataType:
	case Type_structType:
	case Type_interfaceType:
	case Type_listType:
	case Type_objectType:
		if (!m->v.num) {
			str_addf(&SRC, "err = err || capn_setp(%s, %d, %s);\n",
					ptr, m->f.offset, pvar);
		} else if (!strcmp(m->v.tname, "capn_ptr")) {
			g_nullused = 1;
			str_addf(&SRC, "err = err || capn_setp(%s, %d, (%s.data != capn_val%d.data) ? %s : capn_null);\n",
					ptr, m->f.offset, pvar, (int)m->v.num, pvar);
		} else {
			g_nullused = 1;
			str_addf(&SRC, "err = err || capn_setp(%s, %d, (%s.data != capn_val%d.p.data) ? %s : capn_null);\n",
					ptr, m->f.offset, pvar, (int)m->v.num, pvar);
		}
		break;
	}
}

static void get_member(struct member *m, const char *ptr, const char *tab, const char *var) {
	const char *xor = xor_member(m);
	const char *pvar = ptr_member(m, var);

	if (m->v.t.body_tag == Type_voidType)
		return;

	str_add(&SRC, tab, -1);

	switch (m->v.t.body_tag) {
	case Type_voidType:
		return;
	case Type_boolType:
		str_addf(&SRC, "%s = (capn_read8(%s, %d) & %d) != %d;\n",
				var, ptr, m->f.offset/8, 1 << (m->f.offset%8), (int)m->v.num);
		return;
	case Type_int8Type:
		str_addf(&SRC, "%s = (int8_t) capn_read8(%s, %d)%s;\n", var, ptr, m->f.offset, xor);
		return;
	case Type_int16Type:
		str_addf(&SRC, "%s = (int16_t) capn_read16(%s, %d)%s;\n", var, ptr, 2*m->f.offset, xor);
		return;
	case Type_int32Type:
		str_addf(&SRC, "%s = (int32_t) capn_read32(%s, %d)%s;\n", var, ptr, 4*m->f.offset, xor);
		return;
	case Type_int64Type:
		str_addf(&SRC, "%s = (int64_t) capn_read64(%s, %d)%s;\n", var, ptr, 8*m->f.offset, xor);
		return;
	case Type_uint8Type:
		str_addf(&SRC, "%s = capn_read8(%s, %d)%s;\n", var, ptr, m->f.offset, xor);
		return;
	case Type_uint16Type:
		str_addf(&SRC, "%s = capn_read16(%s, %d)%s;\n", var, ptr, 2*m->f.offset, xor);
		return;
	case Type_uint32Type:
		str_addf(&SRC, "%s = capn_read32(%s, %d)%s;\n", var, ptr, 4*m->f.offset, xor);
		return;
	case Type_uint64Type:
		str_addf(&SRC, "%s = capn_read64(%s, %d)%s;\n", var, ptr, 8*m->f.offset, xor);
		return;
	case Type_float32Type:
		str_addf(&SRC, "%s = capn_to_f32(capn_read32(%s, %d)%s);\n", var, ptr, 4*m->f.offset, xor);
		return;
	case Type_float64Type:
		str_addf(&SRC, "%s = capn_to_f64(capn_read64(%s, %d)%s);\n", var, ptr, 8*m->f.offset, xor);
		return;
	case Type_enumType:
		str_addf(&SRC, "%s = (%s) capn_read16(%s, %d)%s;\n", var, m->v.tname, ptr, 2*m->f.offset, xor);
		return;
	case Type_textType:
		if (!m->v.num)
			g_val0used = 1;
		str_addf(&SRC, "%s = capn_get_text(%s, %d, capn_val%d);\n", var, ptr, m->f.offset, (int)m->v.num);
		return;

	case Type_dataType:
		str_addf(&SRC, "%s = capn_get_data(%s, %d);\n", var, ptr, m->f.offset);
		break;
	case Type_structType:
	case Type_interfaceType:
	case Type_objectType:
	case Type_listType:
		str_addf(&SRC, "%s = capn_getp(%s, %d);\n", pvar, ptr, m->f.offset);
		break;
	}

	if (m->v.num) {
		str_addf(&SRC, "%sif (!%s.type) {\n", tab, pvar);
		str_addf(&SRC, "%s\t%s = capn_val%d;\n", tab, var, (int)m->v.num);
		str_addf(&SRC, "%s}\n", tab);
	}
}

static void union_block(struct member *m, struct member *u, int set) {
	static struct str buf = STR_INIT;
	if (set) {
		set_member(u, "p.p", "\t\t", strf(&buf, "s->%s.%s", m->m.name.str, u->m.name.str));
	} else {
		get_member(u, "p.p", "\t\t", strf(&buf, "s->%s.%s", m->m.name.str, u->m.name.str));
	}

	str_addf(&SRC, "\t\tbreak;\n");
}

static void union_cases(struct node *n, struct member *m, int set, int mask) {
	struct member *u = NULL;
	int j;

	for (j = 0; j < m->u.members.p.len; j++) {
		if (!m->mbrs[j].v.num && (mask & (1 << m->mbrs[j].v.t.body_tag))) {
			u = &m->mbrs[j];
			str_addf(&SRC, "\tcase %s_%s:\n", n->name.str, u->m.name.str);
		}
	}

	if (u)
		union_block(m, u, set);
}

static void do_union(struct node *n, struct member *m, int set) {
	int j;

	if (set) {
		str_addf(&SRC, "\terr = err || capn_write16(p.p, %d, s->%s_tag);\n",
				2*m->u.discriminantOffset, m->m.name.str);
	} else {
		str_addf(&SRC, "\ts->%s_tag = (enum %s_%s) capn_read16(p.p, %d);\n",
				m->m.name.str, n->name.str, m->m.name.str, 2*m->u.discriminantOffset);
	}

	str_addf(&SRC, "\n\tswitch (s->%s_tag) {\n", m->m.name.str);

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

	/* when we have defaults we have to emit each case seperately */
	for (j = 0; j < m->u.members.p.len; j++) {
		struct member *u = &m->mbrs[j];
		if (u->v.num) {
			str_addf(&SRC, "\tcase %s_%s:\n", n->name.str, u->m.name.str);
			union_block(m, u, set);
		}
	}

	str_addf(&SRC, "\t}\n");
}

static void print_member(struct member *m, const char *tab) {
	switch (m->v.t.body_tag) {
	case Type_voidType:
		break;
	case Type_boolType:
		fprintf(HDR, "%s%s %s:1;\n", tab, m->v.tname, m->m.name.str);
		break;
	default:
		fprintf(HDR, "%s%s %s;\n", tab, m->v.tname, m->m.name.str);
		break;
	}
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
			print_member(m, "\t");
			break;
		case StructNode_Member_unionMember:
			fprintf(HDR, "\tenum %s_%s %s_tag;\n", n->name.str, m->m.name.str, m->m.name.str);
			fprintf(HDR, "\tunion {\n");
			for (j = 0; j < m->u.members.p.len; j++) {
				print_member(&m->mbrs[j], "\t\t");
			}
			fprintf(HDR, "\t} %s;\n", m->m.name.str);
			break;
		}
	}
	fprintf(HDR, "};\n");

	str_addf(&SRC, "\n%s_ptr new_%s(struct capn_segment *s) {\n", n->name.str, n->name.str);
	str_addf(&SRC, "\t%s_ptr p;\n", n->name.str);
	str_addf(&SRC, "\tp.p = capn_new_struct(s, %d, %d);\n", 8*n->s.dataSectionWordSize, n->s.pointerSectionSize);
	str_addf(&SRC, "\treturn p;\n");
	str_addf(&SRC, "}\n");

	str_addf(&SRC, "%s_list new_%s_list(struct capn_segment *s, int len) {\n", n->name.str, n->name.str);
	str_addf(&SRC, "\t%s_list p;\n", n->name.str);
	str_addf(&SRC, "\tp.p = capn_new_list(s, len, %d, %d);\n", 8*n->s.dataSectionWordSize, n->s.pointerSectionSize);
	str_addf(&SRC, "\treturn p;\n");
	str_addf(&SRC, "}\n");

	str_addf(&SRC, "void read_%s(struct %s *s, %s_ptr p) {\n", n->name.str, n->name.str, n->name.str);
	for (i = 0; i < mlen; i++) {
		struct member *m = &mbrs[i];
		if (!m->is_valid) continue;

		switch (m->m.body_tag) {
		case StructNode_Member_fieldMember:
			get_member(m, "p.p", "\t", strf(&buf, "s->%s", m->m.name.str));
			break;
		case StructNode_Member_unionMember:
			do_union(n, m, 0);
			break;
		}
	}
	str_addf(&SRC, "}\n");

	str_addf(&SRC, "int write_%s(const struct %s *s, %s_ptr p) {\n", n->name.str, n->name.str, n->name.str);
	str_addf(&SRC, "\tint err = 0;\n");
	for (i = 0; i < mlen; i++) {
		struct member *m = &mbrs[i];
		if (!m->is_valid) continue;

		switch (m->m.body_tag) {
		case StructNode_Member_fieldMember:
			set_member(m, "p.p", "\t", strf(&buf, "s->%s", m->m.name.str));
			break;
		case StructNode_Member_unionMember:
			do_union(n, m, 1);
			break;
		}
	}
	str_addf(&SRC, "\treturn err;\n}\n");

	str_addf(&SRC, "void get_%s(struct %s *s, %s_list l, int i) {\n", n->name.str, n->name.str, n->name.str);
	str_addf(&SRC, "\t%s_ptr p;\n", n->name.str);
	str_addf(&SRC, "\tp.p = capn_getp(l.p, i);\n");
	str_addf(&SRC, "\tread_%s(s, p);\n", n->name.str);
	str_addf(&SRC, "}\n");

	str_addf(&SRC, "int set_%s(const struct %s *s, %s_list l, int i) {\n", n->name.str, n->name.str, n->name.str);
	str_addf(&SRC, "\t%s_ptr p;\n", n->name.str);
	str_addf(&SRC, "\tp.p = capn_getp(l.p, i);\n");
	str_addf(&SRC, "\treturn write_%s(s, p);\n", n->name.str);
	str_addf(&SRC, "}\n");
}

static int find_offset(struct str *v, int inc, uint64_t mask) {
	int i, j;
	union {uint64_t u; char c[8];} umask;
	umask.u = capn_flip64(mask);
	for (i = 0; i < v->len*8; i += inc) {
		for (j = i; j < i+inc; j++) {
			if (((uint8_t*)v->str)[j/8] & (1 << (j%8))) {
				goto loop;
			}
		}

		for (j = i; j < i+inc; j++) {
			((uint8_t*)v->str)[j/8] |= 1 << (j%8);
		}

		return j/inc;
loop:
		continue;
	}
	str_add(v, umask.c, 8);
	return (v->len-8)/inc;
}

static void define_method(struct node *iface, int ord) {
	static struct str buf = STR_INIT;
	struct member *mbrs;
	struct InterfaceNode_Method method;
	int i, ptrs = 0, datasz;

	get_InterfaceNode_Method(&method, iface->i.methods, ord);
	mbrs = calloc(method.params.p.len, sizeof(*mbrs));

	str_reset(&buf);

	for (i = 0; i < method.params.p.len; i++) {
		struct InterfaceNode_Method_Param param;
		struct member *m = &mbrs[i];

		get_InterfaceNode_Method_Param(&param, method.params, i);
		decode_value(&m->v, param.type, param.defaultValue, NULL);
		m->m.name = param.name;
		m->m.annotations = param.annotations;
		m->m.ordinal = i;

		switch (m->v.t.body_tag) {
		case Type_voidType:
			break;
		case Type_boolType:
			m->f.offset = find_offset(&buf, 1, 1);
			break;
		case Type_int8Type:
		case Type_uint8Type:
			m->f.offset = find_offset(&buf, 8, 0xFF);
			break;
		case Type_int16Type:
		case Type_uint16Type:
		case Type_enumType:
			m->f.offset = find_offset(&buf, 16, 0xFFFF);
			break;
		case Type_int32Type:
		case Type_uint32Type:
		case Type_float32Type:
			m->f.offset = find_offset(&buf, 32, 0xFFFFFFFFu);
			break;
		case Type_int64Type:
		case Type_uint64Type:
		case Type_float64Type:
			m->f.offset = find_offset(&buf, 64, ~((uint64_t) 0));
			break;
		case Type_textType:
		case Type_dataType:
		case Type_listType:
		case Type_structType:
		case Type_interfaceType:
		case Type_objectType:
			m->f.offset = ptrs++;
			break;
		}
	}

	datasz = buf.len;


	/* write function to initiate a call */

	fprintf(HDR, "\nint write_%s_%s(struct capn_msg*", iface->name.str, method.name.str);
	str_addf(&SRC, "\nint write_%s_%s(struct capn_msg *m", iface->name.str, method.name.str);
	for (i = 0; i < method.params.p.len; i++) {
		struct member *m = &mbrs[i];
		fprintf(HDR, ", %s %s", m->v.tname, m->m.name.str);
		str_addf(&SRC, ", %s a%d", m->v.tname, i);
	}
	fprintf(HDR, ");\n");
	str_addf(&SRC, ") {\n");

	str_addf(&SRC, "\tint err = 0;\n");
	str_addf(&SRC, "\tm->method = %d;\n", ord);
	str_addf(&SRC, "\tm->iface = ((uint64_t) %#xu << 32) | %#xu;\n", (uint32_t) (iface->n.id >> 32), (uint32_t) iface->n.id);

	if (datasz || ptrs) {
		str_addf(&SRC, "\tm->args = capn_new_struct(m->seg, %d, %d);\n", datasz, ptrs);
	} else {
		g_nullused = 1;
		str_addf(&SRC, "\tm->args = capn_null;\n");
	}

	for (i = 0; i < method.params.p.len; i++) {
		set_member(&mbrs[i], "m->args", "\t", strf(&buf, "a%d", i));
	}

	str_addf(&SRC, "\treturn err;\n");
	str_addf(&SRC, "}\n");


	/* read function to handle a call */

	if (datasz || ptrs) {
		fprintf(HDR, "void read_%s_%s(struct capn_msg*", iface->name.str, method.name.str);
		str_addf(&SRC, "void read_%s_%s(struct capn_msg *m", iface->name.str, method.name.str);
		for (i = 0; i < method.params.p.len; i++) {
			struct member *m = &mbrs[i];
			fprintf(HDR, ", %s *%s", m->v.tname, m->m.name.str);
			str_addf(&SRC, ", %s *a%d", m->v.tname, i);
		}
		fprintf(HDR, ");\n");
		str_addf(&SRC, ") {\n");

		for (i = 0; i < method.params.p.len; i++) {
			get_member(&mbrs[i], "m->args", "\t", strf(&buf, "*a%d", i));
		}

		str_addf(&SRC, "}\n");
	}

	free(mbrs);
}

static void declare(struct node *n, enum Node_body type, const char *format, int num) {
	fprintf(HDR, "\n");
	for (n = n->first_child; n != NULL; n = n->next_child) {
		if (n->n.body_tag == type) {
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
	int i, j;

	if (capn_init_fp(&capn, stdin, 0)) {
		fprintf(stderr, "failed to read schema from stdin\n");
		return -1;
	}

	g_valseg.data = calloc(1, capn.seglist->len);
	g_valseg.cap = capn.seglist->len;

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
		static struct str b = STR_INIT;
		struct node *s;
		char *p;
		FILE *srcf;

		g_valc = 0;
		g_valseg.len = 0;
		g_val0used = 0;
		g_nullused = 0;
		capn_init_malloc(&g_valcapn);
		capn_append_segment(&g_valcapn, &g_valseg);

		n = find_node(capn_get64(req.requestedFiles, i));

		HDR = fopen(strf(&b, "%s.h", n->n.displayName.str), "w");
		if (!HDR) {
			fprintf(stderr, "failed to open %s: %s\n", b.str, strerror(errno));
			exit(2);
		}

		str_reset(&SRC);
		srcf = fopen(strf(&b, "%s.c", n->n.displayName.str), "w");
		if (!srcf) {
			fprintf(stderr, "failed to open %s: %s\n", b.str, strerror(errno));
			exit(2);
		}

		fprintf(HDR, "#ifndef CAPN_%X%X\n", (uint32_t) (n->n.id >> 32), (uint32_t) n->n.id);
		fprintf(HDR, "#define CAPN_%X%X\n", (uint32_t) (n->n.id >> 32), (uint32_t) n->n.id);
		fprintf(HDR, "/* AUTO GENERATED - DO NOT EDIT */\n");
		fprintf(HDR, "#include <capn.h>\n");

		for (i = 0; i < n->f.imports.p.len; i++) {
			struct FileNode_Import im;
			get_FileNode_Import(&im, n->f.imports, i);
			fprintf(HDR, "#include \"%s.h\"\n", im.name.str);
		}

		fprintf(HDR, "\n#ifdef __cplusplus\nextern \"C\" {\n#endif\n");


		declare(n, Node_structNode, "struct %s;\n", 1);
		declare(n, Node_structNode, "typedef struct {capn_ptr p;} %s_ptr;\n", 1);
		declare(n, Node_structNode, "typedef struct {capn_ptr p;} %s_list;\n", 1);
		declare(n, Node_interfaceNode, "typedef struct {capn_ptr p;} %s_ptr;\n", 1);
		declare(n, Node_interfaceNode, "typedef struct {capn_ptr p;} %s_list;\n", 1);

		for (s = n->first_child; s != NULL; s = s->next_child) {
			switch (s->n.body_tag) {
			case Node_structNode:
				define_struct(s);
				break;
			case Node_enumNode:
				define_enum(s);
				break;
			case Node_constNode:
				define_const(s);
				break;
			case Node_interfaceNode:
				for (j = 0; j < s->i.methods.p.len; j++) {
					define_method(s, j);
				}
				break;
			default:
				break;
			}
		}

		declare(n, Node_structNode, "%s_ptr new_%s(struct capn_segment*);\n", 2);
		declare(n, Node_structNode, "%s_list new_%s_list(struct capn_segment*, int len);\n", 2);
		declare(n, Node_structNode, "void read_%s(struct %s*, %s_ptr);\n", 3);
		declare(n, Node_structNode, "int write_%s(const struct %s*, %s_ptr);\n", 3);
		declare(n, Node_structNode, "void get_%s(struct %s*, %s_list, int i);\n", 3);
		declare(n, Node_structNode, "int set_%s(const struct %s*, %s_list, int i);\n", 3);

		p = strrchr(n->n.displayName.str, '/');
		fprintf(srcf, "#include \"%s.h\"\n", p ? p+1 : n->n.displayName.str);
		fprintf(srcf, "/* AUTO GENERATED - DO NOT EDIT */\n");

		if (g_val0used)
			fprintf(srcf, "static const capn_text capn_val0 = {0,\"\"};\n");
		if (g_nullused)
			fprintf(srcf, "static const capn_ptr capn_null = {CAPN_NULL};\n");

		if (g_valseg.len) {
			fprintf(srcf, "static const uint8_t capn_buf[%d] = {", g_valseg.len-8);
			for (i = 8; i < g_valseg.len; i++) {
				if (i > 8)
					fprintf(srcf, ",");
				if ((i % 8) == 0)
					fprintf(srcf, "\n\t");
				fprintf(srcf, "%u", ((uint8_t*)g_valseg.data)[i]);
			}
			fprintf(srcf, "\n};\n");

			fprintf(srcf, "static const struct capn_segment capn_seg = {{0},0,0,0,(char*)&capn_buf[0],%d,%d};\n",
					g_valseg.len-8, g_valseg.len-8);
		}

		fwrite(SRC.str, 1, SRC.len, srcf);

		fprintf(HDR, "\n#ifdef __cplusplus\n}\n#endif\n#endif\n");
		fclose(HDR);
		fclose(srcf);
		capn_free(&g_valcapn);
		HDR = srcf = NULL;
	}

	return 0;
}
