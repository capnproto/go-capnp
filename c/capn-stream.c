/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"
#include <string.h>

#ifndef min
static int min(int a, int b) { return (a < b) ? a : b; }
#endif

int capn_deflate(struct capn_stream* s) {
	while (s->avail_in) {
		int i, sz = 0;
		uint8_t hdr = 0;
		uint8_t *p;

		if (!s->avail_out)
			return CAPN_NEED_MORE_OUT;

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
			return CAPN_NEED_MORE_IN;

		for (i = 0; i < 8; i++) {
			if (s->next_in[i]) {
				sz ++;
				hdr |= 1 << i;
			}
		}

		switch (sz) {
		case 0:
			if (s->avail_out < 2)
				return CAPN_NEED_MORE_OUT;

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
				return CAPN_NEED_MORE_OUT;

			s->next_out[0] = hdr;
			memcpy(s->next_out+1, s->next_in, 8);
			s->next_in += 8;
			s->avail_in -= 8;

			s->raw = min(s->avail_in, 256*8);
			if ((p = memchr(s->next_in, 0, s->raw)) != NULL) {
				s->raw = (p - s->next_in) & ~7;
			}

			s->next_out[9] = (uint8_t) (s->raw/8);
			s->next_out += 10;
			s->avail_out -= 10;
			continue;

		default:
			if (s->avail_out < 1 + sz)
				return CAPN_NEED_MORE_OUT;

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
	for (;;) {
		int i, sz;
		uint8_t hdr;

		if (s->zeros > 0) {
			if (s->avail_out == 0)
				return CAPN_NEED_MORE_OUT;

			sz = min(s->avail_out, s->zeros);
			memset(s->next_out, 0, sz);
			s->next_out += sz;
			s->avail_out -= sz;
			s->zeros -= sz;
			continue;
		}

		if (s->raw > 0) {
			if (s->avail_in == 0)
				return CAPN_NEED_MORE_IN;
			else if (s->avail_out == 0)
				return CAPN_NEED_MORE_OUT;

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
		else if (s->avail_out < 8)
			return CAPN_NEED_MORE_OUT;
		else if (s->avail_in < 2)
			return CAPN_NEED_MORE_IN;

		switch (s->next_in[0]) {
		case 0xFF:
			/* 0xFF is followed by 8 bytes raw, followed by
			 * a byte with length in words to read raw */
			if (s->avail_in < 10)
				return CAPN_NEED_MORE_IN;

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
				return CAPN_NEED_MORE_IN;

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
}

