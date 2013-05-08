/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

#ifndef min
static int min(int a, int b) { return (a < b) ? a : b; }
#endif

int capn_deflate(struct capn_stream* s) {
	if (s->avail_in % 8) {
		return CAPN_MISALIGNED;
	}

	while (s->avail_in) {
		int i, sz = 0;
		uint8_t hdr = 0;
		uint8_t *p;

		if (!s->avail_out)
			return CAPN_NEED_MORE;

		if (s->raw > 0) {
			sz = min(s->raw, min(s->avail_in, s->avail_out));
			memcpy(s->next_out, s->next_in, sz);
			s->next_out += sz;
			s->next_in += sz;
			s->avail_out -= sz;
			s->avail_in -= sz;
			s->raw -= sz;
			continue;
		}

		if (s->avail_in < 8)
			return CAPN_NEED_MORE;

		for (i = 0; i < 8; i++) {
			if (s->next_in[i]) {
				sz ++;
				hdr |= 1 << i;
			}
		}

		switch (sz) {
		case 0:
			if (s->avail_out < 2)
				return CAPN_NEED_MORE;

			s->next_out[0] = hdr;
			for (sz = 1; sz < min(s->avail_in/8, 256); sz++) {
				if (((uint64_t*) s->next_in)[sz] != 0) {
					break;
				}
			}

			s->next_out[1] = (uint8_t) (sz-1);
			s->next_in += sz*8;
			s->avail_in -= sz*8;
			s->next_out += 2;
			s->avail_out -= 2;
			continue;

		case 8:
			if (s->avail_out < 10)
				return CAPN_NEED_MORE;

			s->next_out[0] = hdr;
			memcpy(s->next_out+1, s->next_in, 8);
			s->next_in += 8;
			s->avail_in -= 8;

			s->raw = min(s->avail_in, 256*8);
			if ((p = (uint8_t*) memchr(s->next_in, 0, s->raw)) != NULL) {
				s->raw = (p - s->next_in) & ~7;
			}

			s->next_out[9] = (uint8_t) (s->raw/8);
			s->next_out += 10;
			s->avail_out -= 10;
			continue;

		default:
			if (s->avail_out < 1 + sz)
				return CAPN_NEED_MORE;

			*(s->next_out++) = hdr;
			for (i = 0; i < 8; i++) {
				if (s->next_in[i]) {
					*(s->next_out++) = s->next_in[i];
				}
			}
			s->avail_out -= sz + 1;
			s->next_in += 8;
			s->avail_in -= 8;
			continue;
		}
	}

	return 0;
}

int capn_inflate(struct capn_stream* s) {
	if (s->avail_out % 8) {
		return CAPN_MISALIGNED;
	}

	while (s->avail_out) {
		int i, sz;
		uint8_t hdr;

		if (s->zeros > 0) {
			sz = min(s->avail_out, s->zeros);
			memset(s->next_out, 0, sz);
			s->next_out += sz;
			s->avail_out -= sz;
			s->zeros -= sz;
			continue;
		}

		if (s->raw > 0) {
			if (s->avail_in == 0)
				return CAPN_NEED_MORE;

			sz = min(min(s->avail_out, s->raw), s->avail_in);
			memcpy(s->next_out, s->next_in, sz);
			s->next_in += sz;
			s->next_out += sz;
			s->avail_in -= sz;
			s->avail_out -= sz;
			s->raw -= sz;
			continue;
		}

		if (s->avail_in == 0)
			return 0;
		else if (s->avail_in < 2)
			return CAPN_NEED_MORE;

		switch (s->next_in[0]) {
		case 0xFF:
			/* 0xFF is followed by 8 bytes raw, followed by
			 * a byte with length in words to read raw */
			if (s->avail_in < 10)
				return CAPN_NEED_MORE;

			memcpy(s->next_out, s->next_in+1, 8);
			s->raw = (int) s->next_in[9] * 8;
			s->next_in += 10;
			s->avail_in -= 10;
			s->next_out += 8;
			s->avail_out -= 8;
			continue;

		case 0x00:
			/* 0x00 is followed by a single byte indicating
			 * the count of consecutive zero value words
			 * minus 1 */
			s->zeros = (int) (s->next_in[1] + 1) * 8;
			s->next_in += 2;
			s->avail_in -= 2;
			continue;

		default:
			sz = 0;
			hdr = s->next_in[1];
			for (i = 0; i < 8; i++) {
				if (hdr & (1 << i))
					sz++;
			}
			if (s->avail_in < 2 + sz)
				return CAPN_NEED_MORE;

			s->next_in += 2;

			for (i = 0; i < 8; i++) {
				if (hdr & (1 << i)) {
					*(s->next_out++) = *(s->next_in++);
				} else {
					*(s->next_out++) = 0;
				}
			}

			s->avail_out -= 8;
			s->avail_in -= 2 + sz;
			continue;
		}
	}

	return 0;
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
	memset(c, 0, sizeof(*c));

	if (read_fp(&segnum, 4, f, z, zbuf, packed))
		goto err;

	segnum = capn_flip32(segnum);
	if (segnum > 1023)
		goto err;
	segnum++;

	s = (struct capn_segment*) calloc(segnum, sizeof(*s));
	if (!s)
		goto err;

	if (read_fp(hdr, 8 * (segnum/2) + 4, f, z, zbuf, packed))
		goto err;

	for (i = 0; i < segnum; i++) {
		uint32_t n = capn_flip32(hdr[i]);
		if (n > INT_MAX/8 || n > UINT32_MAX/8 || UINT32_MAX - total < n*8)
			goto err;
		s[i].cap = s[i].len = n * 8;
		total += s[i].len;
	}

	data = (char*) calloc(1, total);
	if (!data)
		goto err;

	if (read_fp(data, total, f, z, zbuf, packed))
		goto err;

	for (i = 0; i < segnum; i++) {
		s[i].data = data;
		data += s[i].len;
		capn_append_segment(c, &s[i]);
	}

	return 0;

err:
	memset(c, 0, sizeof(*c));
	free(data);
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

void capn_free_mem(struct capn *c) {
	capn_free_fp(c);
}

void capn_free_fp(struct capn *c) {
	if (c->seglist) {
		free(c->seglist->data);
		free(c->seglist);
	}
}
