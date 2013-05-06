/* vim: set sw=8 ts=8 sts=8 noet: */
#include "capn.h"
#include <stdlib.h>
#include <string.h>

static struct capn_segment *create(void *u, uint32_t id, int sz) {
	struct capn_segment *s;
	if (sz < 4096) {
		sz = 4096;
	} else {
		sz = (sz + 4095) & ~4095;
	}
	s = (struct capn_segment*) calloc(1, sizeof(*s));
	s->data = calloc(1, sz);
	s->len = 0;
	s->cap = sz;
	return s;
}

void capn_init_malloc(struct capn *c) {
	memset(c, 0, sizeof(*c));
	c->create = &create;
}

void capn_free_all(struct capn *c) {
	struct capn_segment *s = c->seglist;
	while (s != NULL) {
		struct capn_segment *n = s->next;
		free(s->data);
		free(s);
		s = n;
	}
	s = c->copylist;
	while (s != NULL) {
		struct capn_segment *n = s->next;
		free(s->data);
		free(s);
		s = n;
	}
}


