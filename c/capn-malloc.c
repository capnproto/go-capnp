/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"
#include <stdlib.h>
#include <string.h>
#include <limits.h>

struct check_segment_alignment {
	unsigned int foo : (sizeof(struct capn_segment)&7) ? -1 : 1;
};

static struct capn_segment *create(void *u, uint32_t id, int sz) {
	struct capn_segment *s;
	sz += sizeof(*s);
	if (sz < 4096) {
		sz = 4096;
	} else {
		sz = (sz + 4095) & ~4095;
	}
	s = (struct capn_segment*) calloc(1, sz);
	s->data = (char*) (s+1);
	s->cap = sz - sizeof(*s);
	s->user = s;
	return s;
}

void capn_init_malloc(struct capn *c) {
	memset(c, 0, sizeof(*c));
	c->create = &create;
}

void capn_free(struct capn *c) {
	struct capn_segment *s = c->seglist;
	while (s != NULL) {
		struct capn_segment *n = s->next;
		free(s->user);
		s = n;
	}
	capn_reset_copy(c);
}

void capn_reset_copy(struct capn *c) {
	struct capn_segment *s = c->copylist;
	while (s != NULL) {
		struct capn_segment *n = s->next;
		free(s->user);
		s = n;
	}
	c->copy = NULL;
	c->copylist = NULL;
}

#define ZBUF_SZ 4096

static int read_fp(void *p, size_t sz, FILE *f, struct capn_stream *z, uint8_t* zbuf, int packed) {
	if (f && packed) {
		z->next_out = (uint8_t*) p;
		z->avail_out = sz;

		while (z->avail_out && capn_inflate(z) == CAPN_NEED_MORE) {
			int r;
			memmove(zbuf, z->next_in, z->avail_in);
			r = fread(zbuf+z->avail_in, 1, ZBUF_SZ - z->avail_in, f);
			if (r <= 0)
				return -1;
			z->avail_in += r;
		}
		return 0;

	} else if (f && !packed) {
		return fread(p, sz, 1, f) != 1;

	} else if (packed) {
		z->next_out = (uint8_t*) p;
		z->avail_out = sz;
		return capn_inflate(z) != 0;

	} else {
		if (z->avail_in < sz)
			return -1;
		memcpy(p, z->next_in, sz);
		z->next_in += sz;
		z->avail_in -= sz;
		return 0;
	}
}

static int init_fp(struct capn *c, FILE *f, struct capn_stream *z, int packed) {
	struct capn_segment *s = NULL;
	uint32_t i, segnum, total = 0;
	uint32_t hdr[1024];
	uint8_t zbuf[ZBUF_SZ];
	char *data = NULL;

	capn_init_malloc(c);

	if (read_fp(&segnum, 4, f, z, zbuf, packed))
		goto err;

	segnum = capn_flip32(segnum);
	if (segnum > 1023)
		goto err;
	segnum++;

	if (read_fp(hdr, 8 * (segnum/2) + 4, f, z, zbuf, packed))
		goto err;

	for (i = 0; i < segnum; i++) {
		uint32_t n = capn_flip32(hdr[i]);
		if (n > INT_MAX/8 || n > UINT32_MAX/8 || UINT32_MAX - total < n*8)
			goto err;
		hdr[i] = n*8;
		total += hdr[i];
	}

	s = (struct capn_segment*) calloc(1, total + (sizeof(*s) * segnum));
	if (!s)
		goto err;

	data = (char*) (s+segnum);
	if (read_fp(data, total, f, z, zbuf, packed))
		goto err;

	for (i = 0; i < segnum; i++) {
		s[i].len = s[i].cap = hdr[i];
		s[i].data = data;
		data += s[i].len;
		capn_append_segment(c, &s[i]);
	}

	s[segnum-1].user = s;

	return 0;

err:
	memset(c, 0, sizeof(*c));
	free(s);
	return -1;
}

int capn_init_fp(struct capn *c, FILE *f, int packed) {
	struct capn_stream z;
	memset(&z, 0, sizeof(z));
	return init_fp(c, f, &z, packed);
}

int capn_init_mem(struct capn *c, const uint8_t *p, size_t sz, int packed) {
	struct capn_stream z;
	memset(&z, 0, sizeof(z));
	z.next_in = p;
	z.avail_in = sz;
	return init_fp(c, NULL, &z, packed);
}
