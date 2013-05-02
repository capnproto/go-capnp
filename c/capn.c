/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"

#include <stdlib.h>
#include <string.h>

#define STRUCT_PTR 0
#define LIST_PTR 1
#define FAR_PTR 2

#define VOID_LIST 0
#define BIT_1_LIST 1
#define BYTE_1_LIST 2
#define BYTE_2_LIST 3
#define BYTE_4_LIST 4
#define BYTE_8_LIST 5
#define PTR_LIST 6
#define COMPOSITE_LIST 7

#define U64(val) ((uint64_t) (val))
#define I64(val) ((int64_t) (val))
#define U32(val) ((uint32_t) (val))
#define I32(val) ((int32_t) (val))
#define U16(val) ((uint16_t) (val))
#define I16(val) ((int16_t) (val))

#ifndef min
static int min(int a, int b) { return (a < b) ? a : b; }
#endif

static struct capn_segment *lookup_segment(struct capn_segment *s, uint32_t id) {
	if (s->id == id)
		return s;
	if (!s->capn || !s->capn->lookup)
		return NULL;
	return s->capn->lookup(s->capn->user, id);
}

static uint64_t lookup_far(struct capn_segment **s, char **d, uint64_t val) {
	uint32_t off = U32(val >> 3);

	if ((*s = lookup_segment(*s, U32(val >> 32))) == NULL) {
		return 0;
	}

	if (val & 4) {
		/* Double far pointer */
		uint64_t far, tag;
		char *p = (*s)->data + off;
		if (off + 16 >= (*s)->len) {
			return 0;
		}

		far = capn_flip64(*(uint64_t*) p);
		tag = capn_flip64(*(uint64_t*) (p+8));

		/* the far tag should not be another double, and the tag
		 * should be struct/list and have no offset */
		if ((far&7) != FAR_PTR || U32(tag) > LIST_PTR) {
			return 0;
		}

		if ((*s = lookup_segment(*s, U32(far >> 32))) == NULL) {
			return 0;
		}

		*d = (*s)->data;
		return U64(U32(far >> 3) << 2) | tag;
	} else {
		if (off + 8 >= (*s)->len) {
			return 0;
		}

		*d = (*s)->data + off;
		return capn_flip64(*(uint64_t*) *d);
	}
}

static char *struct_ptr(struct capn_segment *s, char *d) {
	uint64_t val = capn_flip64(*(uint64_t*)d);
	uint16_t datasz;

	if ((val&3) == FAR_PTR) {
		val = lookup_far(&s, &d, val);
	}

	datasz = U16(val >> 32);
	d += (I32(U32(val)) << 1) + 8;

	if (val != 0 && (val&3) != STRUCT_PTR && !datasz && s->data <= d && d < s->data + s->len) {
		return d;
	}

	return NULL;
}

struct capn_ptr capn_read_ptr(const struct capn_ptr *p, int off) {
	char *d, *e;
	struct capn_ptr ret;
	uint64_t val;

	switch (p->type) {
		case CAPN_LIST:
			/* Return an inner pointer */
			if (off >= p->size) {
				goto err;
			}
			ret.type = CAPN_STRUCT;
			ret.data = p->data + off * (p->datasz + p->ptrsz);
			ret.seg = p->seg;
			ret.datasz = p->datasz;
			ret.ptrsz = p->ptrsz;
			return ret;

		case CAPN_STRUCT:
			off *= 8;
			if (off >= p->ptrsz) {
				goto err;
			}

			d = p->data + p->datasz + off;
			break;

		case CAPN_PTR_LIST:
			if (off >= p->size) {
				goto err;
			}

			d = p->data + off * 8;
			break;

		default:
			goto err;
	}

	val = capn_flip64(*(uint64_t*) d);
	ret.seg = p->seg;

	if ((val&3) == FAR_PTR) {
		val = lookup_far(&ret.seg, &d, val);
	}

	d += (I32(U32(val)) << 1) + 8;
	ret.data = d;

	if ((val&3) > LIST_PTR || d < ret.seg->data) {
		goto err;
	}

	if ((val&3) == STRUCT_PTR) {
		ret.type = CAPN_STRUCT;
		ret.datasz = U32(U16(val >> 32)) * 8;
		ret.ptrsz = U32(U16(val >> 48)) * 8;
		e = d + ret.size * (ret.datasz + ret.ptrsz);
	} else {
		ret.type = CAPN_LIST;
		ret.size = val >> 35;
		ret.datasz = 0;
		ret.ptrsz = 0;

		switch ((val >> 32) & 7) {
			case VOID_LIST:
				e = d;
				break;
			case BIT_1_LIST:
				ret.type = CAPN_BIT_LIST;
				ret.datasz = (ret.size+7)/8;
				e = d + ret.datasz;
				break;
			case BYTE_1_LIST:
				ret.datasz = 1;
				e = d + ret.size;
				break;
			case BYTE_2_LIST:
				ret.datasz = 2;
				e = d + ret.size * 2;
				break;
			case BYTE_4_LIST:
				ret.datasz = 4;
				e = d + ret.size * 4;
				break;
			case BYTE_8_LIST:
				ret.datasz = 8;
				e = d + ret.size * 8;
				break;
			case PTR_LIST:
				ret.type = CAPN_PTR_LIST;
				ret.ptrsz = 8;
				e = d + ret.size * 8;
				break;
			case COMPOSITE_LIST:
				if (d+8-ret.seg->data > ret.seg->len) {
					goto err;
				}

				val = capn_flip64(*(uint64_t*) d);

				d += 8;
				e = d + ret.size * 8;

				ret.datasz = U32(U16(val >> 32)) * 8;
				ret.ptrsz = U32(U16(val >> 48)) * 8;
				ret.size = U32(val >> 2);

				if ((ret.datasz + ret.ptrsz) * ret.size != e - d) {
					goto err;
				}
				break;
		}
	}

	if (e - ret.seg->data <= ret.seg->len) {
		return ret;
	}

err:
	memset(&ret, 0, sizeof(ret));
	return ret;
}

static uint64_t ptr_value(const struct capn_ptr *p, int off) {
	uint64_t val = U64(U32(I32((off >> 3) << 2)));

	switch (p->type) {
	case CAPN_STRUCT:
		val |= STRUCT_PTR | (U64(p->datasz) << 32) | (U64(p->ptrsz) << 48);
		break;

	case CAPN_LIST:
		if (p->ptrsz || p->datasz > 8) {
			val |= LIST_PTR | (U64(COMPOSITE_LIST) << 32) | (U64(p->size * (p->datasz + p->ptrsz)/8) << 35);
		} else if (p->datasz == 8) {
			val |= LIST_PTR | (U64(BYTE_8_LIST) << 32) | (U64(p->size) << 35);
		} else if (p->datasz == 4) {
			val |= LIST_PTR | (U64(BYTE_4_LIST) << 32) | (U64(p->size) << 35);
		} else if (p->datasz == 2) {
			val |= LIST_PTR | (U64(BYTE_2_LIST) << 32) | (U64(p->size) << 35);
		} else if (p->datasz == 1) {
			val |= LIST_PTR | (U64(BYTE_1_LIST) << 32) | (U64(p->size) << 35);
		} else {
			val |= LIST_PTR | (U64(VOID_LIST) << 32) | (U64(p->size) << 35);
		}
		break;

	case CAPN_BIT_LIST:
		val |= LIST_PTR | (U64(BIT_1_LIST) << 32) | (U64(p->size) << 35);
		break;

	case CAPN_PTR_LIST:
		val |= LIST_PTR | (U64(PTR_LIST) << 32) | (U64(p->size) << 35);
		break;

	default:
		val = 0;
		break;
	}

	return capn_flip64(val);
}

static void write_far_ptr(char *d, struct capn_segment *s, char *tgt) {
	*(uint64_t*) d = capn_flip64(FAR_PTR | U64(tgt - s->data) | (U64(s->id) << 32));
}

static void write_double_far(char *d, struct capn_segment *s, char *tgt) {
	*(uint64_t*) d = capn_flip64(FAR_PTR | 4 | U64(tgt - s->data) | (U64(s->id) << 32));
}

static void write_ptr_tag(char *d, const struct capn_ptr *p, int off) {
	*(uint64_t*) d = ptr_value(p, off);
}

static int has_tag(const struct capn_ptr* p) {
	struct capn_segment *s;
	char *d = p->data - 8;
	return d >= s->data && ptr_value(p, 0) == *(uint64_t*) d;
}

static int write_ptr(struct capn_segment *s, char *d, const struct capn_ptr *p) {
	/* note p->seg can be NULL if its a ptr to static data */

	if (!p || p->type == CAPN_NULL) {
		*(uint64_t*) d = 0;
		return 0;

	} else if (p->seg && p->seg == s) {
		write_ptr_tag(d, p, p->data - d - 8);
		return 0;

	} else if (p->seg && p->seg->capn == s->capn && ((p->data - p->seg->data) & 7) == 0) {
		/* if its in the same context we can create a far pointer */

		if (has_tag(p)) {
			/* By lucky chance, the data has a tag in front
			 * of it. This happens when new_data had to move
			 * the data to a new segment. */
			write_far_ptr(d, p->seg, p->data);
			return 0;

		} else if (p->seg->len + 8 <= p->seg->cap) {
			/* The target segment has enough room for tag */
			char *t = p->seg->data + p->seg->len;
			write_ptr_tag(t, p, p->data - t - 8);
			write_far_ptr(d, p->seg, t);
			p->seg->len += 8;
			return 0;

		} else {
			/* have to allocate room for a double far
			 * pointer, but try to allocate it in our
			 * starting segment first */
			char *t;

			if (s->len + 16 > s->cap) {
				if (!s->capn->create)
					return -1;
				if ((s = s->capn->create(s->capn->user, 16)) == NULL)
					return -1;
			}

			t = s->data + s->len;
			write_far_ptr(t, p->seg, p->data);
			write_ptr_tag(t+8, p, 0);
			write_double_far(d, s, t);
			s->len += 16;
			return 0;
		}

	} else {
		/* different context or not aligned - have to copy */
		struct capn_ptr copy;

		switch (p->type) {
		case CAPN_STRUCT:
			copy = capn_new_struct(s, p->datasz, p->ptrsz);
			break;
		case CAPN_PTR_LIST:
			copy = capn_new_ptr_list(s, p->size);
			break;
		case CAPN_BIT_LIST:
			copy = capn_new_bit_list(s, p->size);
			break;
		case CAPN_LIST:
			copy = capn_new_list(s, p->size, p->datasz, p->ptrsz);
			break;
		default:
			return -1;
		}

		return capn_copy(&copy, p) || write_ptr(s, d, &copy);
	}
}

int capn_write_ptr(struct capn_ptr *p, int off, const struct capn_ptr *tgt) {
	struct capn_ptr inner;

	switch (p->type) {
	case CAPN_LIST:
		if (off >= p->size)
			return -1;
		inner = capn_read_ptr(p, off);
		return capn_copy(&inner, tgt);

	case CAPN_PTR_LIST:
		if (off >= p->size)
			return -1;
		return write_ptr(p->seg, p->data + off * 8, tgt);

	case CAPN_STRUCT:
		off *= 8;
		if (off >= p->ptrsz)
			return -1;
		return write_ptr(p->seg, p->data + p->datasz + off, tgt);

	default:
		return -1;
	}
}

static int copy_ptrs(struct capn_ptr *t, const struct capn_ptr *f, int reset_excess) {
	int tptrs = t->ptrsz / 8;
	int fptrs = f->ptrsz / 8;
	int i;

	for (i = 0; i < min(tptrs, fptrs); i++) {
		struct capn_ptr p = capn_read_ptr(f, i);
		if (capn_write_ptr(t, i, &p))
			return -1;
	}

	if (reset_excess) {
		for (i = min(tptrs, fptrs); i < tptrs; i++) {
			capn_write_ptr(t, i, NULL);
		}
	}

	return 0;
}

int capn_copy(struct capn_ptr *t, const struct capn_ptr *f) {
	int fsz = f->size * (f->datasz + f->ptrsz);
	int tsz = t->size * (t->datasz + t->ptrsz);
	int msz = min(fsz, tsz);

	switch (t->type) {
	case CAPN_STRUCT:
		if (f->type == CAPN_STRUCT) {
			/* For structs we reset the excess as reading
			 * from the end of a undersized struct just
			 * reads zeros */
			memcpy(t->data, f->data, msz);
			memset(t->data + msz, 0, tsz - msz);
			return copy_ptrs(t, f, 1);
		} else {
			return -1;
		}

	case CAPN_LIST:
		if (f->type == CAPN_LIST && !f->ptrsz && !t->ptrsz && f->datasz == t->datasz) {
			memcpy(t->data, f->data, msz);
			return 0;
		} else if (f->type == CAPN_LIST || f->type == CAPN_PTR_LIST) {
			return copy_ptrs(t, f, 0);
		} else {
			return -1;
		}

	case CAPN_PTR_LIST:
		if (f->type == CAPN_LIST || f->type == CAPN_PTR_LIST) {
			return copy_ptrs(t, f, 0);
		} else {
			return -1;
		}

	case CAPN_BIT_LIST:
		if (f->type != CAPN_BIT_LIST) {
			memcpy(t->data, f->data, min(t->datasz, f->datasz));
			return 0;
		} else {
			return -1;
		}

	default:
		return -1;
	}
}

int capn_read1(const struct capn_list1 *list, int off, uint8_t *data, int sz) {
	/* Note we only support aligned reads */
	int bsz;
	const struct capn_ptr *p = &list->p;
	if (p->type != CAPN_BIT_LIST || (off & 7) != 0)
		return -1;

	bsz = (sz + 7) / 8;
	off /= 8;

	if (off + sz > p->datasz) {
		memcpy(data, p->data + off, p->datasz - off);
		return p->size - off*8;
	} else {
		memcpy(data, p->data + off, bsz);
		return sz;
	}
}

int capn_write1(struct capn_list1 *list, int off, const uint8_t *data, int sz) {
	/* Note we only support aligned writes */
	int bsz;
	const struct capn_ptr *p = &list->p;
	if (p->type != CAPN_BIT_LIST || (off & 7) != 0)
		return -1;

	bsz = (sz + 7) / 8;
	off /= 8;

	if (off + sz > p->datasz) {
		memcpy(p->data + off, data, p->datasz - off);
		return p->size - off*8;
	} else {
		memcpy(p->data + off, data, bsz);
		return sz;
	}
}

#define SZ 8
#include "capn-list.c"
#undef SZ

#define SZ 16
#include "capn-list.c"
#undef SZ

#define SZ 32
#include "capn-list.c"
#undef SZ

#define SZ 64
#include "capn-list.c"
#undef SZ

static void new_data(struct capn_ptr *p, int bytes) {
	struct capn_segment *s = p->seg;

	/* all allocations are 8 byte aligned */
	bytes = (bytes + 7) & ~7;

	if (s->len + bytes <= s->cap) {
		p->data = p->data + s->len;
		s->len += bytes;
		return;
	}

	/* add a tag whenever we switch segments so that write_ptr can
	 * use it */
	if (!s->capn->create)
		goto err;

	s = s->capn->create(s->capn->user, bytes + 8);
	if (!s)
		goto err;

	write_ptr_tag(s->data + s->len, p, 0);
	return;

err:
	memset(p, 0, sizeof(*p));
}

struct capn_ptr capn_new_struct(struct capn_segment *seg, int datasz, int ptrs) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_STRUCT;
	p.datasz = (datasz + 7) & ~7;
	p.ptrsz = ptrs * 8;
	new_data(&p, p.datasz + p.ptrsz);
	return p;
}

struct capn_ptr capn_new_list(struct capn_segment *seg, int sz, int datasz, int ptrs) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_LIST;
	p.size = sz;

	if (ptrs || datasz > 4) {
		p.datasz = (datasz + 7) & ~7;
		p.ptrsz = ptrs*8;
	} else if (datasz == 3) {
		p.datasz = 4;
		p.ptrsz = 0;
	} else {
		p.datasz = datasz;
		p.ptrsz = 0;
	}

	new_data(&p, p.size * (p.datasz+p.ptrsz));
	return p;
}

struct capn_ptr capn_new_bit_list(struct capn_segment *seg, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_BIT_LIST;
	p.datasz = (sz+7)/8;
	p.size = sz;
	new_data(&p, p.datasz);
	return p;
}

struct capn_ptr capn_new_ptr_list(struct capn_segment *seg, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_PTR_LIST;
	p.size = sz;
	p.ptrsz = 8;
	p.datasz = 0;
	new_data(&p, sz*8);
	return p;
}

struct capn_ptr capn_new_string(struct capn_segment *seg, const char *str, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_LIST;
	p.size = ((sz >= 0) ? sz : strlen(str)) + 1;
	p.datasz = 1;
	p.ptrsz = 0;
	new_data(&p, p.size);
	if (p.data) {
		memcpy(p.data, str, p.size-1);
	}
	return p;
}

char *capn_to_string(const struct capn_ptr *p, int *psz) {
	if (p->type != CAPN_LIST || p->size < 1 || p->data[p->size - 1] != 0) {
		if (psz) *psz = 0;
		return NULL;
	}

	if (psz) *psz = p->size - 1;
	return p->data;
}

struct capn_text capn_read_text(const struct capn_ptr *p, int off) {
	struct capn_text ret;
	if (p->type == CAPN_LIST && p->datasz == 1 && p->size && p->data[p->size - 1] == 0) {
		ret.seg = p->seg;
		ret.str = p->data;
		ret.size = p->size - 1;
	} else {
		ret.seg = NULL;
		ret.str = NULL;
		ret.size = 0;
	}
	return ret;
}

struct capn_data capn_read_data(const struct capn_ptr *p, int off) {
	struct capn_data ret;
	if (p->type == CAPN_LIST && p->datasz == 1) {
		ret.seg = p->seg;
		ret.data = (uint8_t*) p->data;
		ret.size = p->size;
	} else {
		ret.seg = NULL;
		ret.data = NULL;
		ret.size = 0;
	}
	return ret;
}

int capn_write_text(struct capn_ptr *p, int off, struct capn_text tgt) {
	struct capn_ptr m;
	if (tgt.str) {
		m.type = CAPN_LIST;
		m.size = (tgt.size >= 0 ? tgt.size : strlen(tgt.str)) + 1;
		m.seg = tgt.seg;
		m.data = (char*)tgt.str;
		m.datasz = 1;
	} else {
		m.type = CAPN_NULL;
	}
	return capn_write_ptr(p, off, &m);
}

int capn_write_data(struct capn_ptr *p, int off, struct capn_data tgt) {
	struct capn_ptr m;
	if (tgt.data) {
		m.type = CAPN_LIST;
		m.data = (char*)tgt.data;
		m.size = tgt.size;
		m.datasz = 1;
		m.seg = tgt.seg;
	} else {
		m.type = CAPN_NULL;
	}
	return capn_write_ptr(p, off, &m);
}
