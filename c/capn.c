/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"

#include <stdlib.h>
#include <string.h>
#include <sys/param.h>

#define STRUCT_PTR 0
#define LIST_PTR 1
#define FAR_PTR 2
#define DOUBLE_PTR 6

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

struct capn_tree *capn_tree_insert(struct capn_tree *root, struct capn_tree *n) {
	n->red = 1;
	n->link[0] = n->link[1] = NULL;

	for (;;) {
		/* parent, uncle, grandparent, great grandparent link */
		struct capn_tree *p, *u, *g, **gglink;
		int dir;

		/* Case 1: N is root */
		p = n->parent;
		if (!p) {
			n->red = 0;
			root = n;
			break;
		}

		/* Case 2: p is black */
		if (!p->red) {
			break;
		}

		g = p->parent;
		dir = (p == g->link[1]);

		/* Case 3: P and U are red, switch g to red, but must
		 * loop as G could be root or have a red parent
		 *     g    to   G
		 *    / \       / \
		 *   P   U     p   u
		 *  /         /
		 * N         N
		 */
		u = g->link[!dir];
		if (u != NULL && u->red) {
			p->red = 0;
			u->red = 0;
			g->red = 1;
			n = g;
			continue;
		}

		if (!g->parent) {
			gglink = &root;
		} else if (g->parent->link[1] == g) {
			gglink = &g->parent->link[1];
		} else {
			gglink = &g->parent->link[0];
		}

		if (dir != (n == p->link[1])) {
			/* Case 4: rotate on P, then on g
			 * here dir is /
			 *     g    to   g   to   n
			 *    / \       / \      / \
			 *   P   u     N   u    P   G
			 *  / \       / \      /|  / \
			 * 1   N     P   3    1 2 3   u
			 *    / \   / \
			 *   2   3 1   2
			 */
			struct capn_tree *two = n->link[dir];
			struct capn_tree *three = n->link[!dir];
			p->link[!dir] = two;
			g->link[dir] = three;
			n->link[dir] = p;
			n->link[!dir] = g;
			*gglink = n;
			n->parent = g->parent;
			p->parent = n;
			g->parent = n;
			if (two)
				two->parent = p;
			if (three)
				three->parent = g;
			n->red = 0;
			g->red = 1;
		} else {
			/* Case 5: rotate on g
			 * here dir is /
			 *       g   to   p
			 *      / \      / \
			 *     P   u    N   G
			 *    / \      /|  / \
			 *   N   3    1 2 3   u
			 *  / \
			 * 1   2
			 */
			struct capn_tree *three = p->link[!dir];
			g->link[dir] = three;
			p->link[!dir] = g;
			*gglink = p;
			p->parent = g->parent;
			g->parent = p;
			if (three)
				three->parent = g;
			g->red = 1;
			p->red = 0;
		}

		break;
	}

	return root;
}

void capn_append_segment(struct capn *c, struct capn_segment *s) {
	s->id = c->segnum++;
	s->capn = c;
	s->next = NULL;

	if (c->lastseg) {
		c->lastseg->next = s;
		c->lastseg->hdr.link[1] = &s->hdr;
		s->hdr.parent = &c->lastseg->hdr;
	} else {
		c->seglist = s;
		s->hdr.parent = NULL;
	}

	c->lastseg = s;
	c->segtree = capn_tree_insert(c->segtree, &s->hdr);
}

static char *new_data(struct capn *c, int sz, struct capn_segment **ps) {
	struct capn_segment *s;

	/* find a segment with sufficient data */
	for (s = c->seglist; s != NULL; s = s->next) {
		if (s->len + sz <= s->cap) {
			goto end;
		}
	}

	s = c->create ? c->create(c->user, c->segnum, sz) : NULL;
	if (!s) {
		*ps = NULL;
		return NULL;
	}

	capn_append_segment(c, s);
end:
	*ps = s;
	s->len += sz;
	return s->data + s->len - sz;
}

static struct capn_segment *lookup_segment(struct capn* c, struct capn_segment *s, uint32_t id) {
	struct capn_tree **x;
	struct capn_segment *y;

	if (s && s->id == id)
		return s;
	if (!c)
		return NULL;

	if (id < c->segnum) {
		x = &c->segtree;
		y = NULL;
		while (*x) {
			y = (struct capn_segment*) *x;
			if (id == y->id) {
				return y;
			} else if (id < y->id) {
				x = &y->hdr.link[0];
			} else {
				x = &y->hdr.link[1];
			}
		}
	}

	s = c->lookup ? c->lookup(c->user, id) : NULL;
	if (!s)
		return NULL;

	if (id < c->segnum) {
		s->id = id;
		s->capn = c;
		s->next = c->seglist;
		c->seglist = s;
		s->hdr.parent = &y->hdr;
		*x = &s->hdr;
		c->segtree = capn_tree_insert(c->segtree, &s->hdr);
	} else {
		c->segnum = id;
		capn_append_segment(c, s);
	}

	return s;
}

static uint64_t lookup_double(struct capn_segment **s, char **d, uint64_t val) {
	uint64_t far, tag;
	uint32_t off = (U32(val) >> 3) * 8;
	char *p;

	if ((*s = lookup_segment((*s)->capn, *s, U32(val >> 32))) == NULL) {
		return 0;
	}

	p = (*s)->data + off;
	if (off + 16 > (*s)->len) {
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

	/* -8 because far pointers reference from the start of
	 * the segment, but offsets reference the end of the
	 * pointer data. Here *d points to where an equivalent
	 * ptr would be.
	 */
	*d = (*s)->data - 8;
	return U64(U32(far) >> 3 << 2) | tag;
}

static uint64_t lookup_far(struct capn_segment **s, char **d, uint64_t val) {
	uint32_t off = (U32(val) >> 3) * 8;

	if ((*s = lookup_segment((*s)->capn, *s, U32(val >> 32))) == NULL) {
		return 0;
	}

	if (off + 8 > (*s)->len) {
		return 0;
	}

	*d = (*s)->data + off;
	return capn_flip64(*(uint64_t*)*d);
}

static char *struct_ptr(struct capn_segment *s, char *d, int minsz) {
	uint64_t val = capn_flip64(*(uint64_t*)d);
	uint16_t datasz;

	switch (val&7) {
	case FAR_PTR:
		val = lookup_far(&s, &d, val);
		break;
	case DOUBLE_PTR:
		val = lookup_double(&s, &d, val);
		break;
	}

	datasz = U16(val >> 32);
	d += (I32(U32(val)) << 1) + 8;

	if (val != 0 && (val&3) != STRUCT_PTR && datasz >= minsz && s->data <= d && d < s->data + s->len) {
		return d;
	}

	return NULL;
}

static capn_ptr read_ptr(struct capn_segment *s, char *d) {
	capn_ptr ret = {CAPN_NULL};
	uint64_t val;
	char *e;

	val = capn_flip64(*(uint64_t*) d);

	switch (val&7) {
	case FAR_PTR:
		val = lookup_far(&s, &d, val);
		ret.has_ptr_tag = (U32(val) >> 2) == 0;
		break;
	case DOUBLE_PTR:
		val = lookup_double(&s, &d, val);
		break;
	}

	d += (I32(U32(val)) >> 2) * 8 + 8;

	if ((val&3) > LIST_PTR || d < s->data) {
		goto err;
	}

	if ((val&3) == STRUCT_PTR) {
		ret.type = CAPN_STRUCT;
		ret.datasz = U32(U16(val >> 32)) * 8;
		ret.ptrsz = U32(U16(val >> 48)) * 8;
		e = d + ret.len * (ret.datasz + ret.ptrsz);
	} else {
		ret.type = CAPN_LIST;
		ret.len = val >> 35;

		switch ((val >> 32) & 7) {
		case VOID_LIST:
			e = d;
			break;
		case BIT_1_LIST:
			ret.type = CAPN_BIT_LIST;
			ret.datasz = (ret.len+7)/8;
			e = d + ret.datasz;
			break;
		case BYTE_1_LIST:
			ret.datasz = 1;
			e = d + ret.len;
			break;
		case BYTE_2_LIST:
			ret.datasz = 2;
			e = d + ret.len * 2;
			break;
		case BYTE_4_LIST:
			ret.datasz = 4;
			e = d + ret.len * 4;
			break;
		case BYTE_8_LIST:
			ret.datasz = 8;
			e = d + ret.len * 8;
			break;
		case PTR_LIST:
			ret.type = CAPN_PTR_LIST;
			e = d + ret.len * 8;
			break;
		case COMPOSITE_LIST:
			if (d+8-s->data > s->len) {
				goto err;
			}

			val = capn_flip64(*(uint64_t*) d);

			d += 8;
			e = d + ret.len * 8;

			ret.datasz = U32(U16(val >> 32)) * 8;
			ret.ptrsz = U32(U16(val >> 48)) * 8;
			ret.len = U32(val) >> 2;
			ret.has_composite_tag = 1;

			if ((ret.datasz + ret.ptrsz) * ret.len != e - d) {
				goto err;
			}
			break;
		}
	}

	if (e - s->data > s->len)
		goto err;

	ret.data = d;
	ret.seg = s;
	return ret;
err:
	memset(&ret, 0, sizeof(ret));
	return ret;
}

capn_ptr capn_getp(capn_ptr p, int off) {
	switch (p.type) {
	case CAPN_LIST:
		/* Return an inner pointer */
		if (off < p.len) {
			capn_ptr ret = {CAPN_NULL};
			ret.type = CAPN_STRUCT;
			ret.is_list_member = 1;
			ret.data = p.data + off * (p.datasz + p.ptrsz);
			ret.seg = p.seg;
			ret.datasz = p.datasz;
			ret.ptrsz = p.ptrsz;
			return ret;
		} else {
			goto err;
		}

	case CAPN_STRUCT:
		off *= 8;
		if (off >= p.ptrsz) {
			goto err;
		}

		return read_ptr(p.seg, p.data + p.datasz + off);

	case CAPN_PTR_LIST:
		if (off >= p.len) {
			goto err;
		}

		return read_ptr(p.seg, p.data + off * 8);

	default:
		goto err;
	}

err:
	memset(&p, 0, sizeof(p));
	return p;
}

static uint64_t ptr_value(capn_ptr p, int off) {
	uint64_t val = U64(U32(I32(off/8) << 2));

	switch (p.type) {
	case CAPN_STRUCT:
		val |= STRUCT_PTR | (U64(p.datasz/8) << 32) | (U64(p.ptrsz/8) << 48);
		break;

	case CAPN_LIST:
		if (p.has_composite_tag) {
			val |= LIST_PTR | (U64(COMPOSITE_LIST) << 32) | (U64(p.len * (p.datasz + p.ptrsz)/8) << 35);
		} else if (p.datasz == 8) {
			val |= LIST_PTR | (U64(BYTE_8_LIST) << 32) | (U64(p.len) << 35);
		} else if (p.datasz == 4) {
			val |= LIST_PTR | (U64(BYTE_4_LIST) << 32) | (U64(p.len) << 35);
		} else if (p.datasz == 2) {
			val |= LIST_PTR | (U64(BYTE_2_LIST) << 32) | (U64(p.len) << 35);
		} else if (p.datasz == 1) {
			val |= LIST_PTR | (U64(BYTE_1_LIST) << 32) | (U64(p.len) << 35);
		} else {
			val |= LIST_PTR | (U64(VOID_LIST) << 32) | (U64(p.len) << 35);
		}
		break;

	case CAPN_BIT_LIST:
		val |= LIST_PTR | (U64(BIT_1_LIST) << 32) | (U64(p.len) << 35);
		break;

	case CAPN_PTR_LIST:
		val |= LIST_PTR | (U64(PTR_LIST) << 32) | (U64(p.len) << 35);
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
	*(uint64_t*) d = capn_flip64(DOUBLE_PTR | U64(tgt - s->data) | (U64(s->id) << 32));
}

static void write_ptr_tag(char *d, capn_ptr p, int off) {
	*(uint64_t*) d = ptr_value(p, off);
}

#define NEED_TO_COPY 1

static int write_ptr(struct capn_segment *s, char *d, capn_ptr p) {
	/* note p.seg can be NULL if its a ptr to static data */
	char *pdata = p.data;

	if (p.has_composite_tag) {
		pdata -= 8;
	}

	if (p.type == CAPN_NULL) {
		*(uint64_t*) d = 0;
		return 0;

	} else if (!p.seg || p.seg->capn != s->capn || p.is_list_member) {
		return NEED_TO_COPY;

	} else if (p.seg == s) {
		write_ptr_tag(d, p, pdata - d - 8);
		return 0;

	} else {
		/* if its in the same context we can create a far pointer */

		if (p.has_ptr_tag) {
			/* By lucky chance, the data has a tag in front
			 * of it. This happens when new_object had to move
			 * the data to a new segment. */
			write_far_ptr(d, p.seg, pdata-8);
			return 0;

		} else if (p.seg->len + 8 <= p.seg->cap) {
			/* The target segment has enough room for tag */
			char *t = p.seg->data + p.seg->len;
			write_ptr_tag(t, p, pdata - t - 8);
			write_far_ptr(d, p.seg, t);
			p.seg->len += 8;
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
				t = new_data(s->capn, 16, &s);
				if (!t) return -1;
			}

			write_far_ptr(t, p.seg, pdata);
			write_ptr_tag(t+8, p, 0);
			write_double_far(d, s, t);
			return 0;
		}
	}
}

struct copy {
	struct capn_tree hdr;
	struct capn_ptr to, from;
	char *fbegin, *fend;
};

static int data_size(const struct capn_ptr *p) {
	switch (p->type) {
	case CAPN_BIT_LIST:
		return p->datasz;
	case CAPN_PTR_LIST:
		return p->len*8;
	case CAPN_STRUCT:
		return p->datasz + p->ptrsz;
	case CAPN_LIST:
		return p->len * (p->datasz + p->ptrsz);
	default:
		return 0;
	}
}

static capn_ptr new_clone(struct capn_segment *s, capn_ptr p) {
	switch (p.type) {
	case CAPN_STRUCT:
		return capn_new_struct(s, p.datasz, p.ptrsz/8);
	case CAPN_PTR_LIST:
		return capn_new_ptr_list(s, p.len);
	case CAPN_BIT_LIST:
		return capn_new_list1(s, p.len).p;
	case CAPN_LIST:
		return capn_new_list(s, p.len, p.datasz, p.ptrsz/8);
	default:
		return p;
	}
}

static int is_ptr_equal(const struct capn_ptr *a, const struct capn_ptr *b) {
	return a->data == b->data
		&& a->type == b->type
		&& a->len == b->len
		&& a->datasz == b->datasz
		&& a->ptrsz == b->ptrsz
		&& a->has_composite_tag == b->has_composite_tag;
}

static int copy_ptr(struct capn_segment *seg, char *data, struct capn_ptr *t, struct capn_ptr *f, int *dep) {
	struct capn *c = seg->capn;
	struct copy *cp = NULL;
	struct capn_tree **xcp;
	char *fbegin = f->data;
	char *fend = fbegin + data_size(f);

	if (f->has_composite_tag) {
		fbegin -= 8;
	} else if (f->is_list_member) {
		fend = fbegin;
	}

	/* We always copy list members as it would otherwise be an
	 * overlapped pointer (the data is owned by the enclosing list).
	 * We do not bother with the overlapped lookup for zero sized
	 * structures/lists as they never overlap. Nor do we add them to
	 * the copy list as there is no data to be shared by multiple
	 * pointers.
	 */

	xcp = &c->copy;
	while (*xcp && fend > fbegin) {
		cp = (struct copy*) *xcp;
		if (fend <= cp->fbegin) {
			xcp = &cp->hdr.link[0];
		} else if (cp->fend <= fbegin) {
			xcp = &cp->hdr.link[1];
		} else if (is_ptr_equal(f, &cp->from)) {
			/* we already have a copy so just point to that */
			return write_ptr(seg, data, cp->to);
		} else {
			/* pointer to overlapped data */
			return -1;
		}
	}

	/* no copy found - have to create a new copy */
	*t = new_clone(seg, *f);

	if (write_ptr(seg, data, *t))
		return -1;

	/* add the copy to the copy tree so we can look for overlapping
	 * source pointers and handle recursive structures */
	if (fend > fbegin) {
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

		n->from = *f;
		n->to = *t;
		n->fbegin = fbegin;
		n->fend = fend;

		*xcp = &n->hdr;
		n->hdr.parent = &cp->hdr;

		c->copy = capn_tree_insert(c->copy, &n->hdr);
	}

	/* minimize the number of types the main copy routine has to
	 * deal with to just CAPN_LIST and CAPN_PTR_LIST. ptr list only
	 * needs t->type, t->len, t->data, t->seg, f->data, f->seg to
	 * be valid */
	switch (t->type) {
	case CAPN_STRUCT:
		if (t->datasz) {
			memcpy(t->data, f->data, t->datasz);
			t->data += t->datasz;
			f->data += t->datasz;
		}
		if (t->ptrsz) {
			t->type = CAPN_PTR_LIST;
			t->len = t->ptrsz/8;
			(*dep)++;
		}
		return 0;

	case CAPN_BIT_LIST:
		memcpy(t->data, f->data, t->datasz);
		return 0;

	case CAPN_LIST:
		if (!t->len) {
			/* empty list - nothing to copy */
		} else if (t->ptrsz && t->datasz) {
			(*dep)++;
		} else if (t->datasz) {
			memcpy(t->data, f->data, t->len * t->datasz);
		} else if (t->ptrsz) {
			t->type = CAPN_PTR_LIST;
			t->len *= t->ptrsz/8;
			(*dep)++;
		}
		return 0;

	case CAPN_PTR_LIST:
		if (t->len) {
			(*dep)++;
		}
		return 0;

	default:
		return -1;
	}
}

void copy_list_member(capn_ptr* t, capn_ptr *f, int *dep) {
	/* copy struct data */
	int sz = min(t->datasz, f->datasz);
	memcpy(t->data, f->data, sz);
	memset(t->data + sz, 0, t->datasz - sz);
	t->data += t->datasz;
	f->data += f->datasz;

	/* reset excess pointers */
	sz = min(t->ptrsz, f->ptrsz);
	memset(t->data + sz, 0, t->ptrsz - sz);

	/* create a pointer list for the main loop to copy */
	if (t->ptrsz) {
		t->type = CAPN_PTR_LIST;
		t->len = t->ptrsz/8;
		(*dep)++;
	}
}

#define MAX_COPY_DEPTH 32

int capn_setp(capn_ptr p, int off, capn_ptr tgt) {
	struct capn_ptr to[MAX_COPY_DEPTH], from[MAX_COPY_DEPTH];
	char *data;
	int err, dep = 0;

	switch (p.type) {
	case CAPN_LIST:
		if (off >= p.len || tgt.type != CAPN_STRUCT)
			return -1;

		to[0] = p;
		to[0].data += off * (p.datasz + p.ptrsz);
		from[0] = tgt;
		copy_list_member(to, from, &dep);
		break;

	case CAPN_PTR_LIST:
		if (off >= p.len)
			return -1;
		data = p.data + off * 8;
		goto copy_ptr;

	case CAPN_STRUCT:
		off *= 8;
		if (off >= p.ptrsz)
			return -1;
		data = p.data + p.datasz + off;
		goto copy_ptr;

	copy_ptr:
		err = write_ptr(p.seg, data, tgt);
		if (err != NEED_TO_COPY)
			return err;

		/* Depth first copy the source whilst using a pointer stack to
		 * maintain the ptr to set and size left to copy at each level.
		 * We also maintain a rbtree (capn->copy) of the copies indexed
		 * by the source data. This way we can detect overlapped
		 * pointers in the source (and bail) and recursive structures
		 * (and point to the previous copy).
		 */

		from[0] = tgt;
		if (copy_ptr(p.seg, data, to, from, &dep))
			return -1;
		break;

	default:
		return -1;
	}

	while (dep) {
		struct capn_ptr *tc = &to[dep-1], *tn = &to[dep];
		struct capn_ptr *fc = &from[dep-1], *fn = &from[dep];

		if (dep+1 == MAX_COPY_DEPTH) {
			return -1;
		}

		if (!tc->len) {
			dep--;
			continue;
		}

		switch (tc->type) {
		case CAPN_LIST:
			*fn = capn_getp(*fc, 0);
			*tn = capn_getp(*tc, 0);

			copy_list_member(tn, fn, &dep);

			fc->data += fc->datasz + fc->ptrsz;
			tc->data += tc->datasz + tc->ptrsz;
			tc->len--;
			break;

		case CAPN_PTR_LIST:
		default:
			*fn = read_ptr(fc->seg, fc->data);

			if (copy_ptr(tc->seg, tc->data, tn, fn, &dep))
				return -1;

			fc->data += 8;
			tc->data += 8;
			tc->len--;
			break;
		}
	}

	return 0;
}

int capn_get1(capn_list1 l, int off) {
	return l.p.type == CAPN_BIT_LIST
		&& off < l.p.len
		&& (l.p.data[off/8] & (1 << (off%8))) != 0;
}

int capn_set1(capn_list1 l, int off, int val) {
	if (l.p.type != CAPN_BIT_LIST || off >= l.p.len)
		return -1;
	if (val) {
		l.p.data[off/8] |= 1 << (off%8);
	} else {
		l.p.data[off/8] &= ~(1 << (off%8));
	}
	return 0;
}

int capn_getv1(capn_list1 l, int off, uint8_t *data, int sz) {
	/* Note we only support aligned reads */
	int bsz;
	capn_ptr p = l.p;
	if (p.type != CAPN_BIT_LIST || (off & 7) != 0)
		return -1;

	bsz = (sz + 7) / 8;
	off /= 8;

	if (off + sz > p.datasz) {
		memcpy(data, p.data + off, p.datasz - off);
		return p.len - off*8;
	} else {
		memcpy(data, p.data + off, bsz);
		return sz;
	}
}

int capn_setv1(capn_list1 l, int off, const uint8_t *data, int sz) {
	/* Note we only support aligned writes */
	int bsz;
	capn_ptr p = l.p;
	if (p.type != CAPN_BIT_LIST || (off & 7) != 0)
		return -1;

	bsz = (sz + 7) / 8;
	off /= 8;

	if (off + sz > p.datasz) {
		memcpy(p.data + off, data, p.datasz - off);
		return p.len - off*8;
	} else {
		memcpy(p.data + off, data, bsz);
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

/* pull out whether we add a tag or not as a define so the unit test can
 * test double far pointers by not creating tags */
#ifndef ADD_TAG
#define ADD_TAG 1
#endif

static void new_object(capn_ptr *p, int bytes) {
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
	p->data = new_data(s->capn, bytes + ADD_TAG*8, &p->seg);
	if (!p->data) {
		memset(p, 0, sizeof(*p));
		return;
	}

	if (ADD_TAG) {
		write_ptr_tag(p->data, *p, 0);
		p->data += 8;
		p->has_ptr_tag = 1;
	}
}

capn_ptr capn_get_root(struct capn* c) {
	struct capn_segment* s = lookup_segment(c, NULL, 0);
	if (s->len < 8) {
		capn_ptr ret = {CAPN_NULL};
		return ret;
	} else {
		return read_ptr(s, s->data);
	}
}

capn_ptr capn_new_root(struct capn *c) {
	capn_ptr p = {CAPN_NULL};
	struct capn_segment *s = lookup_segment(c, NULL, 0);

	/* don't use new_object as we don't want the tag */
	if ((s || new_data(c, 8, &s) != NULL) && s->len >= 8) {
		p.seg = s;
		p.data = p.seg->data;
		p.len = 1;
		p.type = CAPN_PTR_LIST;
	}

	return p;
}

capn_ptr capn_new_struct(struct capn_segment *seg, int datasz, int ptrs) {
	capn_ptr p = {CAPN_NULL};
	p.seg = seg;
	p.type = CAPN_STRUCT;
	p.datasz = (datasz + 7) & ~7;
	p.ptrsz = ptrs * 8;
	new_object(&p, p.datasz + p.ptrsz);
	return p;
}

capn_ptr capn_new_list(struct capn_segment *seg, int sz, int datasz, int ptrs) {
	capn_ptr p = {CAPN_NULL};
	p.seg = seg;
	p.type = CAPN_LIST;
	p.len = sz;

	if (ptrs || datasz > 8) {
		p.datasz = (datasz + 7) & ~7;
		p.ptrsz = ptrs*8;
		p.has_composite_tag = 1;
		new_object(&p, p.len * (p.datasz + p.ptrsz) + 8);
		if (p.data) {
			uint64_t hdr = STRUCT_PTR | (U64(p.len) << 2) | (U64(p.datasz/8) << 32) | (U64(ptrs) << 48);
			*(uint64_t*) p.data = capn_flip64(hdr);
			p.data += 8;
		}
	} else if (datasz > 4) {
		p.datasz = 8;
		new_object(&p, p.len * 8);
	} else if (datasz > 2) {
		p.datasz = 4;
		new_object(&p, p.len * 4);
	} else {
		p.datasz = datasz;
		new_object(&p, p.len * datasz);
	}

	return p;
}

capn_list1 capn_new_list1(struct capn_segment *seg, int sz) {
	capn_list1 l = {{CAPN_NULL}};
	l.p.seg = seg;
	l.p.type = CAPN_BIT_LIST;
	l.p.datasz = (sz+7)/8;
	l.p.len = sz;
	new_object(&l.p, l.p.datasz);
	return l;
}

capn_ptr capn_new_ptr_list(struct capn_segment *seg, int sz) {
	capn_ptr p = {CAPN_NULL};
	p.seg = seg;
	p.type = CAPN_PTR_LIST;
	p.len = sz;
	p.ptrsz = 0;
	p.datasz = 0;
	new_object(&p, sz*8);
	return p;
}

capn_ptr capn_new_string(struct capn_segment *seg, const char *str, int sz) {
	capn_ptr p = {CAPN_NULL};
	p.seg = seg;
	p.type = CAPN_LIST;
	p.len = ((sz >= 0) ? sz : strlen(str)) + 1;
	p.datasz = 1;
	new_object(&p, p.len);
	if (p.data) {
		memcpy(p.data, str, p.len-1);
	}
	return p;
}

capn_text capn_get_text(capn_ptr p, int off, capn_text def) {
	capn_ptr m = capn_getp(p, off);
	capn_text ret = def;
	if (m.type == CAPN_LIST && m.datasz == 1 && m.len && m.data[m.len - 1] == 0) {
		ret.seg = m.seg;
		ret.str = m.data;
		ret.len = m.len - 1;
	}
	return ret;
}

int capn_set_text(capn_ptr p, int off, capn_text tgt) {
	capn_ptr m = {CAPN_NULL};
	if (tgt.seg) {
		m.type = CAPN_LIST;
		m.seg = tgt.seg;
		m.data = (char*)tgt.str;
		m.len = tgt.len + 1;
		m.datasz = 1;
	} else if (tgt.str) {
		m = capn_new_string(p.seg, tgt.str, tgt.len);
	}
	return capn_setp(p, off, m);
}

capn_data capn_get_data(capn_ptr p, int off) {
	capn_data ret = {capn_getp(p, off)};
	if (ret.p.type != CAPN_LIST || ret.p.datasz != 1) {
		memset(&ret, 0, sizeof(ret));
	}
	return ret;
}
