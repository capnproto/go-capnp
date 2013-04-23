/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"

#include <stdlib.h>
#include <string.h>

#define STRUCT_PTR 0
#define LIST_PTR 1
#define FAR_PTR 2
#define DOUBLE_FAR_PTR 6

struct capn {
	struct capn_segment **seg;
	int seg_nr, seg_alloc;
	capn_lookup_t lookup;
	capn_create_t create;
	void *lookup_user;
	void *create_user;
};

struct capn_segment {
	struct capn *capn;
	uint32_t id;
	char *start, *end, *used;
	void (*free)(void*);
};

struct capn *capn_new(void) {
	return calloc(1, sizeof(struct capn));
}

#ifndef min
static int min(int a, int b) { return (a < b) ? a : b; }
#endif

void capn_free(struct capn* c) {
	if (c) {
		int i;
		for (i = 0; i < c->seg_nr; i++) {
			struct capn_segment *s = c->seg[i];
			if (s->free) {
				s->free(s->start);
			}
			free(s);
		}
		free(c->seg);
		free(c);
	}
}

void capn_set_lookup(struct capn* c, capn_lookup_t lookup, void *user) {
	c->lookup = lookup;
	c->lookup_user = user;
}

void capn_set_create(struct capn* c, capn_create_t create, void *user) {
	c->create = create;
	c->create_user = user;
}

/* Finds the first position where inserting id would not violate the
 * ordering. */
static int find_segment(struct capn *c, uint32_t id) {
	int b = 0, sz = c->seg_nr;
	while (sz > 0) {
		int half = sz >> 1;
		int m = b + half;
		if (c->seg[m]->id < id) {
			b = m + 1;
			sz = sz - half - 1;
		} else {
			sz = half;
		}
	}
	return b;
}

void capn_add(struct capn *c, struct capn_segment *s) {
	if (c->seg_nr <= c->seg_alloc) {
		c->seg_alloc = (c->seg_alloc * 2) + 16;
		c->seg = realloc(c->seg, c->seg_alloc * sizeof(*c->seg));
		memset(c->seg, 0, (c->seg_alloc-c->seg_nr)*sizeof(*c->seg));
	}

	s->capn = c;
	s->id = c->seg_nr ? c->seg[c->seg_nr-1]->id : 0;
	c->seg[c->seg_nr++] = s;
}

static struct capn_segment *lookup_far(struct capn_segment *s, uint32_t id) {
	struct capn *c = s->capn;
	int i;

	if (s->id == id)
		return s;

	i = find_segment(c, id);

	if (i < c->seg_nr && c->seg[i]->id == id)
		return c->seg[i];

	if (!c->lookup)
		return NULL;

	s = c->lookup(c->lookup_user, id);
	if (!s)
		return NULL;

	capn_add(c, s);
	memmove(c->seg+i+1, c->seg+i, sizeof(s) * (c->seg_nr-i-1));
	s->id = id;
	c->seg[i] = s;

	return s;
}

static struct capn_segment *new_data(struct capn_segment *s, int bytes) {
	int i;
	struct capn *c = s->capn;

	if (s->used + bytes <= s->end)
		return s;

	for (i = 0; i < c->seg_nr; i++) {
		s = c->seg[i];
		if (s->used + bytes <= s->end) {
			return s;
		}
	}

	if (!c->create)
		return NULL;

	s = c->create(c->create_user, bytes);
	if (!s)
		return NULL;

	capn_add(c, s);
	return s;
}

struct capn_segment *capn_new_segment(void *data, int len, int cap, void (*free)(void*)) {
	struct capn_segment *s = calloc(1, sizeof(*s));
	s->free = free;
	s->start = data;
	s->end = s->start + cap;
	s->used = s->start + len;
	return s;
}


static int is_ptr_valid(struct capn_ptr *p) {
	/* TODO: be careful of wraparound if start < p->data < used but the
	 * size overflows the type */
	if (p->data < p->seg->start)
		return 0;

	switch (p->type) {
	case CAPN_STRUCT:
		return p->data + 8*(p->u.s.datasz + p->u.s.ptrsz) <= p->seg->used;
	case CAPN_1BIT_LIST:
		return p->data + (p->elements + 7) / 8 <= p->seg->used;
	case CAPN_1BYTE_LIST:
		return p->data + p->elements <= p->seg->used;
	case CAPN_2BYTE_LIST:
		return p->data + 2*p->elements <= p->seg->used;
	case CAPN_4BYTE_LIST:
		return p->data + 4*p->elements <= p->seg->used;
	case CAPN_8BYTE_LIST:
	case CAPN_POINTER_LIST:
		return p->data + 8*p->elements <= p->seg->used;
	case CAPN_COMPOSITE_LIST:
		return p->data + 8*p->elements*(p->u.s.datasz + p->u.s.ptrsz) <= p->seg->used;
	case CAPN_FAR:
		return p->data + 8 <= p->seg->used;
	case CAPN_DOUBLE_FAR:
		return p->data + 16 <= p->seg->used;
	default:
		return 0;
	}
}

static struct capn_ptr read_ptr(struct capn_segment *s, char *d) {
	uint64_t val = capn_from_le_64(*(uint64_t*) d);
	struct capn_ptr p;

	if (val == 0) {
		memset(&p, 0, sizeof(p));
		return p;
	}

	switch (val & 3) {
	case STRUCT_PTR:
		p.seg = s;
		p.type = CAPN_STRUCT;
		p.elements = 1;
		p.data = d + 8 + 8 * (((int)(int32_t)(uint32_t) val) >> 2);
		p.u.s.datasz = (uint16_t)(val >> 32);
		p.u.s.ptrsz = (uint16_t)(val >> 48);
		if (!is_ptr_valid(&p))
			goto err;
		return p;

	case FAR_PTR:
		p.seg = s;
		p.type = (val & 4) ? CAPN_DOUBLE_FAR : CAPN_FAR;
		p.elements = ((uint32_t) val) >> 3;
		p.u.id = (uint32_t)(val >> 32);
		p.data = NULL;
		return p;

	case LIST_PTR:
		p.seg = s;
		p.type = 8 + ((val >> 32) & 7);
		p.elements = (int)(val >> 35);
		p.u.id = 0;
		p.data = d + 8 + 8 * (((int)(int32_t)(uint32_t) val) >> 2);
		if (p.type == CAPN_COMPOSITE_LIST) {
			uint64_t tag;
			int sz = p.elements;
			if (p.data < s->start || p.data+8 > s->used)
				goto err;
			tag = capn_from_le_64(*(uint64_t*) p.data);
			if ((tag & 3) != STRUCT_PTR)
				goto err;
			p.u.s.datasz = (uint16_t)(tag >> 32);
			p.u.s.ptrsz = (uint16_t)(tag >> 48);
			p.elements = ((int32_t)(uint32_t)tag) >> 2;
			if (p.elements * (p.u.s.datasz + p.u.s.ptrsz) != sz)
				goto err;
			p.data += 8;
		}
		if (!is_ptr_valid(&p))
			goto err;
		return p;
	}

err:
	memset(&p, 0, sizeof(p));
	return p;
}

static void write_struct_ptr(char *p, int offset, int datasz, int ptrsz) {
	*(uint64_t*)p = capn_to_le_64((uint64_t)STRUCT_PTR
		| (uint64_t)(uint32_t)(int32_t)(offset << 2)
		| (((uint64_t)datasz) << 32)
		| (((uint64_t)ptrsz) << 48));
}

static void write_list_ptr(char *p, int offset, enum CAPN_TYPE type, int elems) {
	*(uint64_t*)p = capn_to_le_64((uint64_t)LIST_PTR
		| (((uint64_t)(uint32_t)(int32_t)(offset << 2)))
		| ((uint64_t)(type >> 1) << 32)
		| (((uint64_t)elems) << 35));
}

static void write_far_ptr(char *p, enum CAPN_TYPE type, struct capn_segment *seg, char *data) {
	*(uint64_t*)p = capn_to_le_64((uint64_t)type
		| ((uint64_t)(data - seg->start) << 3)
		| (((uint64_t)seg->id) << 32));
}

static void write_ptr_tag(char *d, struct capn_ptr *p, int offset) {
	switch (p->type) {
	case CAPN_STRUCT:
		write_struct_ptr(d, offset, p->u.s.datasz, p->u.s.ptrsz);
		break;
	case CAPN_COMPOSITE_LIST:
		write_list_ptr(d, offset, CAPN_COMPOSITE_LIST, p->elements*(p->u.s.datasz+p->u.s.ptrsz));
		break;
	default:
		write_list_ptr(d, offset, p->type, p->elements);
		break;
	}
}

static int write_ptr(struct capn_segment *s, char *d, struct capn_ptr *p) {
	struct capn_ptr far = *p;

	if (capn_deref_ptr(p)) {
		*(uint64_t*) d = 0;
		return 0;
	}

	if (p->seg == s) {
		write_ptr_tag(d, p, (p->data - d)/8);
		return 0;

	} else if (p->seg->capn == s->capn) {
		struct capn_segment *n;

		if (far.type == CAPN_FAR || far.type == CAPN_DOUBLE_FAR) {
			write_far_ptr(d, far.type, far.seg, far.seg->start+far.elements);
			return 0;

		} else if (p->seg->used + 8 <= p->seg->end) {
			write_ptr_tag(p->seg->used, p, (p->data - p->seg->used)/8);
			write_far_ptr(d, CAPN_FAR, p->seg, p->seg->used);
			p->seg->used += 8;
			return 0;

		} else if ((n = new_data(s, 16)) != NULL) {
			write_far_ptr(n->used, CAPN_FAR, p->seg, p->data);
			write_ptr_tag(n->used+8, p, 0);
			write_far_ptr(d, CAPN_DOUBLE_FAR, n, n->used);
			n->used += 16;
			return 0;

		} else {
			return -1;
		}

	} else {
		struct capn_ptr copy;

		switch (p->type) {
		case CAPN_STRUCT:
			copy = capn_new_struct(s, p->u.s.datasz, p->u.s.ptrsz, 0);
			break;
		case CAPN_COMPOSITE_LIST:
			copy = capn_new_composite(s, p->elements, p->u.s.datasz, p->u.s.ptrsz, 0);
			break;
		default:
			copy = capn_new_list(s, p->type, p->elements, 0);
			break;
		}

		return capn_copy(&copy, p) || write_ptr(s, d, &copy);
	}
}

int capn_deref_ptr(struct capn_ptr *p) {
	struct capn_ptr tag;

	switch (p->type) {
	case CAPN_FAR:
		p->seg = lookup_far(p->seg, p->elements);
		if (!p->seg)
			goto err;

		p->data = p->seg->start + 8*p->elements;
		if (!is_ptr_valid(p))
			goto err;

		*p = read_ptr(p->seg, p->data);
		if (p->type < CAPN_STRUCT)
			goto err;
		return 0;

	case CAPN_DOUBLE_FAR:
		p->seg = lookup_far(p->seg, p->elements);
		if (!p->seg)
			goto err;

		p->data = p->seg->start + 8*p->elements;
		if (!is_ptr_valid(p))
			goto err;

		tag = read_ptr(p->seg, p->data+8);
		if (tag.type < CAPN_STRUCT || tag.data != p->data+16)
			goto err;

		*p = read_ptr(p->seg, p->data);
		if (p->type != CAPN_FAR)
			goto err;

		p->seg = lookup_far(p->seg, p->elements);
		if (!p->seg)
			goto err;

		p->data = p->seg->start + 8*p->elements;
		p->type = tag.type;
		p->elements = tag.elements;
		p->u.s.datasz = tag.u.s.datasz;
		p->u.s.ptrsz = tag.u.s.ptrsz;
		if (!is_ptr_valid(p))
			goto err;
		return 0;

	case CAPN_NULL:
		return -1;

	default:
		return 0;
	}

err:
	memset(p, 0, sizeof(*p));
	return -1;
}

int capn_copy(struct capn_ptr *to, struct capn_ptr *from) {
	int i, elems, datasz, ptrsz;
	int esz = 1;

	if (capn_deref_ptr(to) || capn_deref_ptr(from))
		return -1;

	if (!(from->type == to->type || (from->type >= CAPN_POINTER_LIST && to->type >= CAPN_POINTER_LIST)))
		return -1;

	elems = min(from->elements, to->elements);
	datasz = min(from->u.s.datasz, to->u.s.datasz);
	ptrsz = min(from->u.s.ptrsz, to->u.s.ptrsz);

	switch (from->type) {
	case CAPN_STRUCT:
		memcpy(to->data, from->data, 8*datasz);
		memset(to->data + 8*datasz, 0, 8*(to->u.s.datasz-datasz));
		for (i = 0; i < ptrsz; i++) {
			struct capn_ptr ptr = read_ptr(from->seg, from->data + 8*(from->u.s.datasz + i));
			if (write_ptr(to->seg, to->data + 8*(to->u.s.datasz + i), &ptr))
				return -1;
		}
		memset(to->data + 8*(to->u.s.datasz+ptrsz), 0, 8*(to->u.s.ptrsz-ptrsz));
		return 0;

	case CAPN_1BIT_LIST:
		memset(to->data, 0, (to->elements+7)/8);
		memcpy(to->data, from->data, (elems+7)/8);
		return 0;

	case CAPN_8BYTE_LIST:
		esz *= 2;
	case CAPN_4BYTE_LIST:
		esz *= 2;
	case CAPN_2BYTE_LIST:
		esz *= 2;
	case CAPN_1BYTE_LIST:
		memcpy(to->data, from->data, esz*elems);
		memset(to->data + esz*elems, 0, esz*(to->elements-elems));
		return 0;

	case CAPN_POINTER_LIST:
	case CAPN_COMPOSITE_LIST:
		for (i = 0; i < elems; i++) {
			struct capn_ptr ptr = capn_read_ptr(from, i);
			if (capn_write_ptr(to, i, &ptr))
				return -1;
		}
		return 0;

	default:
		return -1;
	}
}

char *capn_to_string(struct capn_ptr *p, int *psz) {
	if (capn_deref_ptr(p) || p->type != CAPN_1BYTE_LIST)
		goto err;

	if (p->elements < 1 || p->data[p->elements - 1] != 0)
		goto err;

	if (psz)
		*psz = p->elements - 1;
	return p->data;

err:
	if (psz)
		*psz = 0;
	return NULL;
}

void capn_read_struct(struct capn_ptr *p, void *s, int datasz, int ptrsz) {
	struct capn_ptr *ptr = (struct capn_ptr*) ((char*) s + datasz);
	char *data = s;
	int i, sz;

	if (capn_deref_ptr(p) || p->type != CAPN_STRUCT)
		return;

	sz = min(p->u.s.datasz, datasz);
	memcpy(data, p->data, sz*8);
	memset(data + 8*sz, 0, 8*(datasz-sz));

	sz = min(p->u.s.ptrsz, ptrsz);
	for (i = 0; i < sz; i++) {
		ptr[i] = read_ptr(p->seg, p->data + 8*(p->u.s.datasz + i));
	}
	memset(ptr + sz, 0, sizeof(*ptr)*(ptrsz-sz));
}

void capn_write_struct(struct capn_ptr *p, const void *data, int datasz, int ptrsz) {
	struct capn_ptr *ptr = (struct capn_ptr*) ((char*) data + datasz);
	int i;

	if (capn_deref_ptr(p) || p->type != CAPN_STRUCT)
		return;

	memcpy(p->data, data, min(p->u.s.datasz, datasz)*8);

	for (i = 0; i < min(p->u.s.ptrsz, ptrsz); i++) {
		write_ptr(p->seg, p->data + 8*(p->u.s.datasz+i), &ptr[i]);
	}
}

static void read_1(char *to, char *from, int i, int esz) {
	switch (esz) {
	case 1:
		if (*(uint8_t*) from & 1) {
			((uint8_t*) to)[i/8] |= 1 << (i%8);
		} else {
			((uint8_t*) to)[i/8] &= ~(1 << (i%8));
		}
		break;
	case 8:
		((uint8_t*) to)[i] = *(uint8_t*) from;
		break;
	case 16:
		((uint16_t*) to)[i] = capn_from_le_16(*(uint16_t*) from);
		break;
	case 32:
		((uint32_t*) to)[i] = capn_from_le_32(*(uint32_t*) from);
		break;
	case 64:
		((uint64_t*) to)[i] = capn_from_le_64(*(uint64_t*) from);
		break;
	}
}

static void read_list(struct capn_ptr *p, int off, void *data, int sz, enum CAPN_TYPE etyp, int esz) {
	char *to = (char*) data - off*esz/8;
	int i;

	if (capn_deref_ptr(p) || off >= p->elements)
		return;

	if (off+sz > p->elements) {
		memset(to + p->elements*esz/8, 0, (off+sz-p->elements)*esz/8);
		sz = p->elements - off;
	}

	switch (p->type) {
	case CAPN_1BIT_LIST:
	case CAPN_1BYTE_LIST:
		if (p->type != etyp)
			goto err;

		memcpy(data, p->data + off, sz);
		return;

	case CAPN_POINTER_LIST:
		for (i = off; i < off+sz; i++) {
			struct capn_ptr mbr = read_ptr(p->seg, p->data + 8*i);

			if (mbr.type == CAPN_STRUCT && mbr.u.s.datasz) {
				read_1(to, mbr.data, i, esz);
			} else {
				static char zero[8];
				read_1(to, zero, i, esz);
			}
		}
		return;

	case CAPN_COMPOSITE_LIST:
		if (!p->u.s.datasz)
			goto err;

		int csz = 8*(p->u.s.datasz + p->u.s.ptrsz);

		for (i = off; i < off+sz; i++) {
			read_1(to, p->data + csz*i, i, esz);
		}
		return;

	default:
		if (p->type != etyp)
			goto err;

		for (i = off; i < off+sz; i++) {
			read_1(to, p->data + i*esz/8, i, esz);
		}
		return;
	}

err:
	memset(to, 0, sz*esz/8);
}

void capn_read_1(struct capn_ptr *p, int off, uint8_t *data, int sz) {
	read_list(p, off, data, sz, CAPN_1BIT_LIST, 1);
}
void capn_read_8(struct capn_ptr *p, int off, uint8_t *data, int sz) {
	read_list(p, off, data, sz, CAPN_1BYTE_LIST, 8);
}
void capn_read_16(struct capn_ptr *p, int off, uint16_t *data, int sz) {
	read_list(p, off, data, sz, CAPN_2BYTE_LIST, 16);
}
void capn_read_32(struct capn_ptr *p, int off, uint32_t *data, int sz) {
	read_list(p, off, data, sz, CAPN_4BYTE_LIST, 32);
}
void capn_read_64(struct capn_ptr *p, int off, uint64_t *data, int sz) {
	read_list(p, off, data, sz, CAPN_8BYTE_LIST, 64);
}

static void write_1(char *to, const char *from, int i, int esz) {
	switch (esz) {
	case 1:
		if (((uint8_t*)from)[i/8] & (1 << (i%8))) {
			*(uint8_t*) to |= 1;
		} else {
			*(uint8_t*) to &= ~1;
		}
	case 8:
		*(uint8_t*) to = ((uint8_t*)from)[i];
		break;
	case 16:
		*(uint16_t*) to = capn_to_le_16(((uint16_t*)from)[i]);
		break;
	case 32:
		*(uint32_t*) to = capn_to_le_32(((uint32_t*)from)[i]);
		break;
	case 64:
		*(uint64_t*) to = capn_to_le_64(((uint64_t*)from)[i]);
		break;
	}
}

static void write_list(struct capn_ptr *p, int off, const void *data, int sz, enum CAPN_TYPE etyp, int esz) {
	const char *from = (const char*) data - off*esz/8;
	int i;

	if (capn_deref_ptr(p) || off >= p->elements)
		return;

	if (off+sz > p->elements) {
		sz = p->elements-off;
	}

	switch (p->type) {
	case CAPN_1BIT_LIST:
	case CAPN_1BYTE_LIST:
		if (p->type == etyp) {
			memcpy(p->data + off, data, sz);
		}
		break;

	case CAPN_POINTER_LIST:
		for (i = off; i < off+sz; i++) {
			struct capn_ptr mbr = read_ptr(p->seg, p->data + 8*i);

			if (mbr.type == CAPN_STRUCT && mbr.u.s.datasz) {
				write_1(mbr.data, from, i, esz);
			}
		}
		break;

	case CAPN_COMPOSITE_LIST:
		if (p->u.s.datasz) {
			int csz = 8*(p->u.s.datasz + p->u.s.ptrsz);

			for (i = off; i < off+sz; i++) {
				write_1(p->data + csz*i, from, i, esz);
			}
		}
		break;

	default:
		if (p->type == etyp) {
			for (i = off; i < off+sz; i++) {
				write_1(p->data + i*esz/8, from, i, esz);
			}
		}
		break;
	}

}

void capn_write_1(struct capn_ptr *p, int off, const void *data, int sz) {
	write_list(p, off, data, sz, CAPN_1BIT_LIST, 1);
}
void capn_write_8(struct capn_ptr *p, int off, const uint8_t *data, int sz) {
	write_list(p, off, data, sz, CAPN_1BYTE_LIST, 8);
}
void capn_write_16(struct capn_ptr *p, int off, const uint16_t *data, int sz) {
	write_list(p, off, data, sz, CAPN_2BYTE_LIST, 16);
}
void capn_write_32(struct capn_ptr *p, int off, const uint32_t *data, int sz) {
	write_list(p, off, data, sz, CAPN_4BYTE_LIST, 32);
}
void capn_write_64(struct capn_ptr *p, int off, const uint64_t *data, int sz) {
	write_list(p, off, data, sz, CAPN_8BYTE_LIST, 64);
}


struct capn_ptr capn_read_ptr(struct capn_ptr *p, int off) {
	struct capn_ptr ret;

	if (capn_deref_ptr(p) || off >= p->elements)
		goto err;

	switch (p->type) {
	case CAPN_POINTER_LIST:
		return read_ptr(p->seg, p->data + 8*off);

	case CAPN_COMPOSITE_LIST:
		ret.seg = p->seg;
		ret.type = CAPN_STRUCT;
		ret.elements = 0;
		ret.u.s.datasz = p->u.s.datasz;
		ret.u.s.ptrsz = p->u.s.ptrsz;
		ret.data = p->data + off*8*(p->u.s.datasz + p->u.s.ptrsz);
		return ret;

	default:
		goto err;
	}

err:
	memset(&ret, 0, sizeof(ret));
	return ret;
}

int capn_write_ptr(struct capn_ptr *p, int off, struct capn_ptr *ptr) {
	struct capn_ptr mbr;

	if (capn_deref_ptr(p) || off >= p->elements)
		return -1;

	switch (p->type) {
	case CAPN_POINTER_LIST:
		return write_ptr(p->seg, p->data + 8*off, ptr);

	case CAPN_COMPOSITE_LIST:
		mbr = capn_read_ptr(p, off);
		return capn_copy(&mbr, ptr);

	default:
		return -1;
	}
}

static char *new_value(struct capn_ptr *p, int bytes, int want_tag) {

	if (!want_tag && p->seg->used + bytes <= p->seg->end) {
		p->data = p->seg->used;
		p->seg->used += bytes;
		return p->data;

	} else if ((want_tag && p->seg->used + bytes + 8 <= p->seg->end)
	|| (p->seg = new_data(p->seg, bytes + 8)) != NULL) {
		/* tag requested or forced to since we had to switch
		 * segments */

		char *data = p->seg->used + 8;
		write_ptr_tag(p->seg->used, p, 0);
		p->type = CAPN_FAR;
		p->elements = p->seg->used - p->seg->start;
		p->u.id = p->seg->id;
		p->data = NULL;
		p->seg->used += bytes + 8;
		return data;

	} else {
		memset(p, 0, sizeof(*p));
		return NULL;
	}
}

struct capn_ptr capn_new_struct(struct capn_segment *seg, int datasz, int ptrs, int want_tag) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_STRUCT;
	p.u.s.datasz = datasz;
	p.u.s.ptrsz = ptrs;
	new_value(&p, 8*(datasz + ptrs), want_tag);
	return p;
}

struct capn_ptr capn_new_string(struct capn_segment *seg, const char *str, int sz, int want_tag) {
	char *data;
	struct capn_ptr p;
	if (sz < 0)
		sz = strlen(str);

	p.seg = seg;
	p.type = CAPN_1BYTE_LIST;
	p.elements = sz + 1;
	p.u.id = 0;
	data = new_value(&p, sz + 1, want_tag);
	if (data) {
		memcpy(data, str, sz);
	}

	return p;
}

struct capn_ptr capn_new_list(struct capn_segment *seg, enum CAPN_TYPE type, int sz, int want_tag) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = type;
	p.elements = sz;
	p.u.id = 0;

	switch (type) {
	case CAPN_1BIT_LIST:
		new_value(&p, (sz+7)/8, want_tag);
		break;
	case CAPN_1BYTE_LIST:
		new_value(&p, sz, want_tag);
		break;
	case CAPN_2BYTE_LIST:
		new_value(&p, sz*2, want_tag);
		break;
	case CAPN_4BYTE_LIST:
		new_value(&p, sz*4, want_tag);
		break;
	case CAPN_8BYTE_LIST:
	case CAPN_POINTER_LIST:
		new_value(&p, sz*8, want_tag);
		break;
	default:
		memset(&p, 0, sizeof(p));
		break;
	}

	return p;
}

struct capn_ptr capn_new_composite(struct capn_segment *seg, int elems, int datasz, int ptrs, int want_tag) {
	char *data;
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_COMPOSITE_LIST;
	p.elements = elems;
	p.u.s.datasz = datasz;
	p.u.s.ptrsz = ptrs;

	data = new_value(&p, 8 + 8*elems*(datasz + ptrs), want_tag);
	if (data != NULL) {
		write_struct_ptr(data, elems, datasz, ptrs);

		/* new_value may have returned a far pointer */
		if (p.data == data) {
			p.data += 8;
		}
	}

	return p;
}
