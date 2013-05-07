
#define CAT2(A,B) A ## B
#define CAT(A,B) CAT2(A, B)
#define UINT_T CAT(CAT(uint, SZ), _t)
#define FLIP CAT(capn_flip, SZ)

UINT_T CAT(capn_get,SZ) (const struct capn_ptr *p, int off) {
	char *d;
	if (off >= p->size) {
		return 0;
	}

	switch (p->type) {
	case CAPN_LIST:
		if (p->datasz < SZ/8)
			return 0;
		d = p->data + off * (p->datasz + p->ptrsz);
		return FLIP(*(UINT_T*)d);

	case CAPN_PTR_LIST:
		d = struct_ptr(p->seg, p->data + 8*off, SZ/8);
		if (d) {
			return FLIP(*(UINT_T*)d);
		} else {
			return 0;
		}

	default:
		return 0;
	}
}

int CAT(capn_getv,SZ) (const struct capn_ptr *p, int off, UINT_T *to, int sz) {
	int i;

	if (off + sz > p->size) {
		sz = p->size - off;
	}

	switch (p->type) {
	case CAPN_LIST:
		if (p->datasz == SZ/8 && !p->ptrsz && (SZ == 8 || CAPN_LITTLE)) {
			memcpy(to, p->data + off, sz * (SZ/8));
			return sz;
		} else if (p->datasz < SZ/8) {
			return -1;
		}

		for (i = 0; i < sz; i++) {
			char *d = p->data + (i + off) * (p->datasz + p->ptrsz);
			to[i] = FLIP(*(UINT_T*)d);
		}
		return sz;

	case CAPN_PTR_LIST:
		for (i = 0; i < sz; i++) {
			char *d = struct_ptr(p->seg, p->data + 8*(i+off), SZ/8);
			if (d) {
				to[i] = FLIP(*(UINT_T*)d);
			} else {
				return -1;
			}
		}
		return sz;

	default:
		return -1;
	}
}

int CAT(capn_set,SZ) (struct capn_ptr *p, int off, UINT_T v) {
	char *d;
	if (off >= p->size) {
		return -1;
	}

	switch (p->type) {
	case CAPN_LIST:
		if (p->datasz < SZ/8)
			return -1;
		d = p->data + off * (p->datasz + p->ptrsz);
		*(UINT_T*) d = FLIP(v);
		return 0;

	case CAPN_PTR_LIST:
		d = struct_ptr(p->seg, p->data + 8*off, SZ/8);
		if (!d) {
			return -1;
		}
		*(UINT_T*) d = FLIP(v);
		return 0;

	default:
		return -1;
	}
}

int CAT(capn_setv,SZ) (struct capn_ptr *p, int off, const UINT_T *from, int sz) {
	int i;

	if (off + sz > p->size) {
		sz = p->size - off;
	}

	switch (p->type) {
	case CAPN_LIST:
		if (p->datasz == SZ/8 && !p->ptrsz && (SZ == 8 || CAPN_LITTLE)) {
			memcpy(p->data + off, from, sz * (SZ/8));
			return sz;
		} else if (p->datasz < SZ/8) {
			return -1;
		}

		for (i = 0; i < sz; i++) {
			char *d = p->data + (i + off) * (p->datasz + p->ptrsz);
			*(UINT_T*) d = FLIP(from[i]);
		}
		return sz;

	case CAPN_PTR_LIST:
		for (i = 0; i < sz; i++) {
			char *d = struct_ptr(p->seg, p->data + 8*(i+off), SZ/8);
			if (d) {
				*(UINT_T*) d = FLIP(from[i]);
			} else {
				return -1;
			}
		}
		return sz;

	default:
		return -1;
	}
}

#undef FLIP
#undef UINT_T
#undef CAT

