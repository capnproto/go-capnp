/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"

#include <stdlib.h>
#include <string.h>
#include <sys/param.h>

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

#ifdef BYTE_ORDER
#define CAPN_LITTLE (BYTE_ORDER == LITTLE_ENDIAN)
#elif defined(__BYTE_ORDER)
#define CAPN_LITTLE (__BYTE_ORDER == __LITTLE_ENDIAN)
#else
#define CAPN_LITTLE 0
#endif

/* from    to
 *     G         G
 *    / \       / \
 *   P   U     N   U
 *  / \       / \
 * 1   N     P   3
 *    / \   / \
 *   2   3 1   2
 */
static void rotate_left(struct capn_tree **pp, struct capn_tree **pn) {
	struct capn_tree *p = *pp, *n = *pn;
        struct capn_tree *g = p->parent;
	p->right = n->left;
	n->left = p;
	g->left = n;
	n->parent = g;
	p->parent = n;
	*pn = p;
	*pp = n;
}

/* from     to
 *     G      G
 *    / \    / \
 *   U   P  U   N
 *      / \    / \
 *     N   3  1   P
 *    / \        / \
 *   1   2      2   3
 */
static void rotate_right(struct capn_tree **pp, struct capn_tree **pn) {
	struct capn_tree *p = *pp, *n = *pn;
        struct capn_tree *g = p->parent;
	p->left = n->right;
	n->right = p;
	g->right = n;
	n->parent = g;
	p->parent = n;
	*pn = p;
	*pp = n;
}

static void insert_rebalance(struct capn_tree **root, struct capn_tree *n) {
	n->red = 1;
	n->left = n->right = NULL;

	for (;;) {
		struct capn_tree *p, *u, *g;

		p = n->parent;
		if (!p) {
			*root = n;
			n->red = 0;
			return;
		}

		g = p->parent;
		if (!g->red) {
			return;
		}

		if (p == g->left) {
			if ((u = g->right) != NULL && u->red) {
				p->red = 0;
				u->red = 0;
				g->red = 1;
				n = g;
				continue;
			}
			if (n == p->right) {
				rotate_left(&p, &n);
			}
			p->red = 0;
			g->red = 1;
			rotate_right(&g, &p);
		} else {
			if ((u = g->left) != NULL && u->red) {
				p->red = 0;
				u->red = 0;
				g->red = 1;
				n = g;
				continue;
			}
			if (n == p->left) {
				rotate_right(&p, &n);
			}
			p->red = 0;
			g->red = 1;
			rotate_left(&g, &p);
		}
	}

}

static void *new_data(struct capn_segment **ps, int sz) {
	struct capn *c = (*ps)->capn;
	struct capn_segment *s = *ps;

	/* find a segment with sufficient data */
	for (s = c->seglist; s != NULL; s = s->next) {
		if (s->len + sz <= s->cap) {
			goto end;
		}
	}

	s = c->create ? c->create(c->user, c->segnum, sz) : NULL;
	if (!s) {
		return NULL;
	}

	s->capn = c;
	s->id = c->segnum++;
	s->next = c->seglist;
	c->seglist = s;
	s->hdr.parent = c->lastseg;
	c->lastseg->right = &s->hdr;
	c->lastseg = &s->hdr;
	insert_rebalance(&c->segtree, &s->hdr);

end:
	*ps = s;
	s->len += sz;
	return s->data + s->len - sz;
}

void capn_append_segment(struct capn *c, struct capn_segment *s) {
	s->id = c->segnum++;
	s->capn = c;
	s->next = c->seglist;
	c->seglist = s;
	s->hdr.parent = c->lastseg;
	c->lastseg->right = &s->hdr;
	c->lastseg = &s->hdr;
	insert_rebalance(&c->segtree, &s->hdr);
}

static struct capn_segment *lookup_segment(struct capn* c, struct capn_segment *s, uint32_t id) {
	struct capn_tree **x;
	struct capn_segment *y;

	if (s && s->id == id)
		return s;
	if (!c)
		return NULL;

	x = &c->segtree;
	y = NULL;
	while (*x) {
		y = (struct capn_segment*) *x;
		if (id == y->id) {
			return y;
		} else if (id < y->id) {
			x = &y->hdr.left;
		} else {
			x = &y->hdr.right;
		}
	}

	s = c->lookup ? c->lookup(c->user, id) : NULL;
	if (!s)
		return NULL;

	if (id >= c->segnum) {
		c->lastseg = &s->hdr;
		c->segnum = id+1;
	}

	s->id = id;
	s->capn = c;
	s->next = c->seglist;
	c->seglist = s;
	s->hdr.parent = &y->hdr;
	*x = &s->hdr;
	insert_rebalance(&c->segtree, &s->hdr);

	return s;
}

static uint64_t lookup_far(struct capn_segment **s, char **d, uint64_t val) {
	uint32_t off = U32(val >> 3);

	if ((*s = lookup_segment((*s)->capn, *s, U32(val >> 32))) == NULL) {
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

		if ((*s = lookup_segment((*s)->capn, *s, U32(far >> 32))) == NULL) {
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

static struct capn_ptr read_ptr(struct capn_segment *s, char *d) {
	char *e;
	struct capn_ptr ret;
	uint64_t val;

	val = capn_flip64(*(uint64_t*) d);

	if ((val&3) == FAR_PTR) {
		val = lookup_far(&s, &d, val);
	}

	d += (I32(U32(val)) << 1) + 8;
	ret.data = d;

	if ((val&3) > LIST_PTR || d < s->data) {
		goto err;
	}

	if ((val&3) == STRUCT_PTR) {
		ret.type = CAPN_STRUCT;
		ret.size = 0;
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
			e = d + ret.size * 8;
			break;
		case COMPOSITE_LIST:
			if (d+8-s->data > s->len) {
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

	if (e - s->data > s->len)
		goto err;

	ret.seg = s;
	return ret;
err:
	memset(&ret, 0, sizeof(ret));
	return ret;
}

struct capn_ptr capn_read_ptr(const struct capn_ptr *p, int off) {
	struct capn_ptr ret;

	switch (p->type) {
	case CAPN_LIST:
		/* Return an inner pointer */
		if (off >= p->size) {
			goto err;
		}
		ret = *p;
		ret.type = CAPN_LIST_MEMBER;
		ret.data += off * (p->datasz + p->ptrsz);
		ret.size = 0;
		return ret;

	case CAPN_STRUCT:
	case CAPN_LIST_MEMBER:
		off *= 8;
		if (off >= p->ptrsz) {
			goto err;
		}

		return read_ptr(p->seg, p->data + p->datasz + off);

	case CAPN_PTR_LIST:
		if (off >= p->size) {
			goto err;
		}

		return read_ptr(p->seg, p->data + off * 8);

	default:
		goto err;
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
	char *d = p->data - 8;
	return d >= p->seg->data && ptr_value(p, 0) == *(uint64_t*) d;
}

#define NEED_TO_COPY 1

static int write_ptr(struct capn_segment *s, char *d, const struct capn_ptr *p) {
	/* note p->seg can be NULL if its a ptr to static data */

	if (!p || p->type == CAPN_NULL) {
		*(uint64_t*) d = 0;
		return 0;

	} else if (!p->seg || p->seg->capn != s->capn || p->type == CAPN_LIST_MEMBER) {
		return NEED_TO_COPY;

	} else if (p->seg == s) {
		write_ptr_tag(d, p, p->data - d - 8);
		return 0;

	} else {
		/* if its in the same context we can create a far pointer */

		if (has_tag(p)) {
			/* By lucky chance, the data has a tag in front
			 * of it. This happens when new_object had to move
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
			 * pointer */
			char *t;

			if (s->len + 16 <= s->cap) {
				/* Try and allocate in the src segment
				 * first. This should improve lookup on
				 * read. */
				t = s->data + s->len;
				s->len += 16;
			} else {
				t = new_data(&s, 16);
				if (!t) return -1;
			}

			write_far_ptr(t, p->seg, p->data);
			write_ptr_tag(t+8, p, 0);
			write_double_far(d, s, t);
			return 0;
		}
	}
}

struct copy {
	struct capn_tree hdr;
	struct capn_ptr to, from;
	int fsize;
};

static int data_size(const struct capn_ptr *p) {
	switch (p->type) {
	case CAPN_BIT_LIST:
		return p->datasz;
	case CAPN_PTR_LIST:
		return p->size*8;
	case CAPN_STRUCT:
	case CAPN_LIST_MEMBER:
		return p->datasz + p->ptrsz;
	case CAPN_LIST:
		return p->size * (p->datasz + p->ptrsz);
	default:
		return 0;
	}
}

static struct capn_ptr new_clone(struct capn_segment *s, const struct capn_ptr *p) {
	switch (p->type) {
	case CAPN_LIST_MEMBER:
	case CAPN_STRUCT:
		return capn_new_struct(s, p->datasz, p->ptrsz);
	case CAPN_PTR_LIST:
		return capn_new_ptr_list(s, p->size);
	case CAPN_BIT_LIST:
		return capn_new_bit_list(s, p->size);
	case CAPN_LIST:
		return capn_new_list(s, p->size, p->datasz, p->ptrsz);
	default:
		return *p;
	}

}

static int is_ptr_equal(const struct capn_ptr *a, const struct capn_ptr *b) {
	return a->data == b->data && a->type == b->type && a->size == b->size && a->datasz == b->datasz && a->ptrsz == b->ptrsz;
}

static int write_copy(struct capn_segment *seg, char *data, struct capn_ptr *t, struct capn_ptr *f, int *dep) {
	struct capn *c = seg->capn;
	struct copy *cp = (struct copy*) c->copy;
	int sz = data_size(f);

	/* We always copy list members as it would otherwise be an
	 * overlapped pointer (the data is owned by the inclosing list).
	 * We do not bother with the overlapped lookup for zero sized
	 * structures/lists as they never overlap. Nor do we add them to
	 * the copy list as there is no data to be shared by multiple
	 * pointers.
	 */

	while (c && sz) {
		if (f->data + sz <= cp->from.data) {
			cp = (struct copy*) cp->hdr.left;
		} else if (cp->from.data + cp->fsize <= f->data) {
			cp = (struct copy*) cp->hdr.right;
		} else if (is_ptr_equal(f, &cp->from)) {
			/* we already have a copy so just point to that */
			return write_ptr(seg, data, &cp->from);
		} else {
			/* pointer to overlapped data */
			return -1;
		}
	}

	/* no copy - have to copy */
	*t = new_clone(seg, f);

	/* add the copy to the copy tree so we can look for overlapping
	 * source pointers and handle recursive structures */
	if (sz && f->type != CAPN_LIST_MEMBER) {
		struct copy *n;
		struct capn_segment *cs = c->copylist;

		/* need to allocate a struct copy */
		if (!cs || cs->len + sizeof(*n) > cs->cap) {
			cs = c->create ? c->create(c->user, CAPN_SEGID_LOCAL, sizeof(*n)) : NULL;
			if (!cs) {
				/* can't allocate a copy structure */
				return -1;
			}
			cs->next = c->copylist;
			c->copylist = cs;
		}

		n = (struct copy*) (cs->data + cs->len);
		cs->len += sizeof(*n);

		n->hdr.parent = &cp->hdr;
		n->from = *f;
		n->to = *t;
		n->fsize = sz;

		if (f->data < cp->from.data) {
			cp->hdr.left = &n->hdr;
		} else {
			cp->hdr.right = &n->hdr;
		}

		insert_rebalance(&seg->capn->copy, &n->hdr);
	}

	/* minimize the number of types the main copy routine has to
	 * deal with to just CAPN_LIST and CAPN_PTR_LIST. ptr list only
	 * needs t->type, t->size, t->data, t->seg, f->data, f->seg to
	 * be valid */
	switch (t->type) {
	case CAPN_STRUCT:
	case CAPN_LIST_MEMBER:
		if (t->datasz) {
			memcpy(t->data, f->data, t->datasz);
			t->data += t->datasz;
			f->data += t->datasz;
		}
		if (t->ptrsz) {
			t->type = CAPN_PTR_LIST;
			t->size = t->ptrsz/8;
			(*dep)++;
		}
		return 0;

	case CAPN_BIT_LIST:
		memcpy(t->data, f->data, t->datasz);
		return 0;

	case CAPN_LIST:
		if (!t->size) {
			/* empty list - nothing to copy */
		} else if (t->ptrsz && t->datasz) {
			(*dep)++;
		} else if (t->datasz) {
			memcpy(t->data, f->data, t->size * t->datasz);
		} else if (t->ptrsz) {
			t->type = CAPN_PTR_LIST;
			t->size *= t->ptrsz/8;
			(*dep)++;
		}
		return 0;

	case CAPN_PTR_LIST:
		if (t->size) {
			(*dep)++;
		}
		return 0;

	default:
		return -1;
	}
}

#define MAX_COPY_DEPTH 32

int capn_write_ptr(struct capn_ptr *p, int off, const struct capn_ptr *tgt) {
	struct capn_ptr to[MAX_COPY_DEPTH], from[MAX_COPY_DEPTH];
	char *data;
	int err, dep;

	switch (p->type) {
	case CAPN_LIST:
		if (off < p->size && (tgt->type == CAPN_STRUCT || tgt->type == CAPN_LIST_MEMBER)) {
			struct capn_ptr *f, *t;
			char *d;
			int sz;

			/* copy struct data */
			d = p->data + off * (p->datasz + p->ptrsz);
			sz = min(p->datasz, tgt->datasz);
			memcpy(d, tgt->data, sz);
			memset(d + sz, 0, p->datasz - sz);

			/* reset excess pointers */
			d += p->datasz;
			sz = min(p->ptrsz, tgt->ptrsz);
			memset(d + sz, 0, p->ptrsz - sz);

			/* create a pointer list for the main loop to copy */
			dep = 1;

			/* main copy loop doesn't need the other fields
			 * for ptr lists */
			f = &from[0];
			f->data = tgt->data + tgt->datasz;
			f->seg = tgt->seg;

			t = &to[0];
			t->type = CAPN_PTR_LIST;
			t->data = d;
			t->size = sz/8;
			t->seg = p->seg;

			goto copy_loop;
		} else {
			return -1;
		}

	case CAPN_PTR_LIST:
		if (off >= p->size)
			return -1;
		data = p->data + off * 8;
		break;

	case CAPN_STRUCT:
	case CAPN_LIST_MEMBER:
		off *= 8;
		if (off >= p->ptrsz)
			return -1;
		data = p->data + p->datasz + off;
		break;

	default:
		return -1;
	}

	err = write_ptr(p->seg, data, tgt);
	if (err != NEED_TO_COPY)
		return err;

	/* Depth first copy the source whilst using a pointer stack to
	 * maintain the ptr to set and size left to copy at each level.
	 * We also maintain a rbtree (capn->copy) of the copies indexed
	 * by the source data. This way we can detect overlapped
	 * pointers in the source (and bail) and recursive structures
	 * (and point to the previous copy).
	 */

	dep = 0;
	from[0] = *tgt;
	if (write_copy(p->seg, data, to, from, &dep))
		return -1;

copy_loop:

	while (dep) {
		struct capn_ptr *tc = &to[dep-1], *tn = &to[dep];
		struct capn_ptr *fc = &from[dep-1], *fn = &from[dep];

		if (dep+1 == MAX_COPY_DEPTH) {
			return -1;
		}

		if (!tc->size) {
			dep--;
			continue;
		}

		switch (tc->type) {
		case CAPN_LIST:
			*fn = *fc;
			*tn = *tc;
			fn->type = tn->type = CAPN_LIST_MEMBER;
			fn->size = tn->size = 0;

			if (write_copy(tc->seg, tc->data, tn, fn, &dep))
				return -1;

			fc->data += tc->datasz + tc->ptrsz;
			tc->data += tc->datasz + tc->ptrsz;
			tc->size--;
			break;

		case CAPN_PTR_LIST:
		default:
			*fn = read_ptr(fc->seg, fc->data);

			if (write_copy(tc->seg, tc->data, tn, fn, &dep))
				return -1;

			fc->data += 8;
			tc->data += 8;
			tc->size--;
			break;
		}
	}

	return 0;
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
#include "capn-list.inc"
#undef SZ

#define SZ 16
#include "capn-list.inc"
#undef SZ

#define SZ 32
#include "capn-list.inc"
#undef SZ

#define SZ 64
#include "capn-list.inc"
#undef SZ

static void new_object(struct capn_ptr *p, int bytes) {
	struct capn_segment *s = p->seg;

	/* all allocations are 8 byte aligned */
	bytes = (bytes + 7) & ~7;

	if (s->len + bytes <= s->cap) {
		p->data = s->data + s->len;
		s->len += bytes;
		return;
	}

	/* add a tag whenever we switch segments so that write_ptr can
	 * use it */
	p->data = new_data(&s, bytes + 8);
	if (!p->data) {
		memset(p, 0, sizeof(*p));
		return;
	}

	write_ptr_tag(p->data, p, 0);
	p->data += 8;
}

struct capn_ptr capn_root(struct capn* c) {
	struct capn_ptr p;
	p.seg = lookup_segment(c, NULL, 0);
	if (!p.seg) goto err;
	p.data = p.seg->data;
	p.size = 1;
	p.type = CAPN_PTR_LIST;
	p.datasz = 0;
	p.ptrsz = 0;
	new_object(&p, 8);
	return p;
err:
	memset(&p, 0, sizeof(p));
	return p;
}

struct capn_ptr capn_new_struct(struct capn_segment *seg, int datasz, int ptrs) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_STRUCT;
	p.size = 0;
	p.datasz = (datasz + 7) & ~7;
	p.ptrsz = ptrs * 8;
	new_object(&p, p.datasz + p.ptrsz);
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

	new_object(&p, p.size * (p.datasz+p.ptrsz));
	return p;
}

struct capn_ptr capn_new_bit_list(struct capn_segment *seg, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_BIT_LIST;
	p.datasz = (sz+7)/8;
	p.ptrsz = 0;
	p.size = sz;
	new_object(&p, p.datasz);
	return p;
}

struct capn_ptr capn_new_ptr_list(struct capn_segment *seg, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_PTR_LIST;
	p.size = sz;
	p.ptrsz = 0;
	p.datasz = 0;
	new_object(&p, sz*8);
	return p;
}

struct capn_ptr capn_new_string(struct capn_segment *seg, const char *str, int sz) {
	struct capn_ptr p;
	p.seg = seg;
	p.type = CAPN_LIST;
	p.size = ((sz >= 0) ? sz : strlen(str)) + 1;
	p.datasz = 1;
	p.ptrsz = 0;
	new_object(&p, p.size);
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

int capn_marshal_iptr(const union capn_iptr *ip, struct capn_ptr *p, int off) {
	return capn_write_ptr(p, off, &ip->c);
}
