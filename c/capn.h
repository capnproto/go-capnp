/* vim: set sw=8 ts=8 sts=8 noet: */

#ifndef CAPN_H
#define CAPN_H

#include <stdint.h>

typedef struct capn_segment *(*capn_create_t)(void* /*user*/, int /*sz*/);
typedef struct capn_segment *(*capn_lookup_t)(void* /*user*/, uint32_t /*id*/);

/* struct capn is a common structure shared between segments in the same
 * session/context so that far pointers between the segments will be created.
 * lookup is used to lookup segments by id when derefencing a far pointer
 * create is used to create or lookup an alternate segment that has at least
 * sz available (ie returned seg->len + sz <= seg->cap)
 * create and lookup can be NULL if you don't need multiple segments
 */
struct capn {
	struct capn_segment *(*lookup)(void* /*user*/, uint32_t /*id */);
	struct capn_segment *(*create)(void* /*user*/, int /*sz*/);
	void *user;
};

/* struct capn_segment contains the information about a single segment.
 * capn should point to a struct capn that is shared between segments in the
 * same session
 * id specifies the segment id, used for far pointers
 * data specifies the segment data. This should not move once.
 * len specifies the current segment length. This should be 0 for a blank
 * segment.
 * cap specifies the segment capacity.
 * When creating new structures len will be incremented until it reaces cap,
 * at which point a new segment will be requested via capn->create.
 *
 * data, len, and cap must all by 8 byte aligned
 */
struct capn_segment {
	struct capn *capn;
	uint32_t id;
	char *data;
	int len, cap;
};

enum CAPN_TYPE {
	CAPN_NULL = 0,
	CAPN_STRUCT = 1,
	CAPN_LIST = 2,
	CAPN_PTR_LIST = 3,
	CAPN_BIT_LIST = 4,
};

struct capn_ptr {
	enum CAPN_TYPE type;
	int size;
	char *data;
	struct capn_segment *seg;
	uint32_t datasz;
	uint32_t ptrsz;
};

struct capn_text {
	int size;
	const char *str;
	struct capn_segment *seg;
};

struct capn_data {
	int size;
	const uint8_t *data;
	struct capn_segment *seg;
};

struct capn_list1{struct capn_ptr p;};
struct capn_list8{struct capn_ptr p;};
struct capn_list16{struct capn_ptr p;};
struct capn_list32{struct capn_ptr p;};
struct capn_list64{struct capn_ptr p;};

/* capn_read|write_ptr functions read/write ptrs to list/structs
 * off is the list index or pointer index in a struct
 * capn_write_ptr will copy the data, create far pointers, etc if the target
 * is in a different segment/context.
 * Both of these will use/return inner pointers for composite lists.
 */
struct capn_ptr capn_read_ptr(const struct capn_ptr *p, int off);
int capn_write_ptr(struct capn_ptr *p, int off, const struct capn_ptr *tgt);

/* capn_to_string returns a pointer to a string
 * Use this instead of accessing the data directly as these checks that the
 * string is null terminated, the list type, etc.
 * psz is filled out with the string length if non NULL
 */
struct capn_text capn_read_text(const struct capn_ptr *p, int off);
struct capn_data capn_read_data(const struct capn_ptr *p, int off);
int capn_write_text(struct capn_ptr *p, int off, struct capn_text tgt);
int capn_write_data(struct capn_ptr *p, int off, struct capn_data tgt);

/* capn_copy copies data from 'from' to 'to'
 * returns 0 on success, non-zero on error (type mismatch, allocation error,
 * etc).
 */
int capn_copy(struct capn_ptr *to, const struct capn_ptr *from);

/* capn_read_* functions read data from a list
 * The length of the list is given by p->size
 * off specifies how far into the list to start
 * sz indicates the number of elements to read
 * The function returns the number of elements read or -1 on an error.
 * off must be byte aligned for capn_read_1
 */
int capn_read_1(const struct capn_list1 *p, int off, uint8_t *data, int sz);
int capn_read8(const struct capn_list8 *p, int off, uint8_t *data, int sz);
int capn_read16(const struct capn_list16 *p, int off, uint16_t *data, int sz);
int capn_read32(const struct capn_list32 *p, int off, uint32_t *data, int sz);
int capn_read64(const struct capn_list64 *p, int off, uint64_t *data, int sz);

/* capn_write_* function write data to a list
 * off specifies how far into the list to start
 * sz indicates the number of elements to write
 * The function returns the number of elemnts written or -1 on an error.
 * off must be byte aligned for capn_read_1
 */
int capn_write_1(struct capn_list1 *p, int off, const uint8_t *data, int sz);
int capn_write8(struct capn_list8 *p, int off, const uint8_t *data, int sz);
int capn_write16(struct capn_list16 *p, int off, const uint16_t *data, int sz);
int capn_write32(struct capn_list32 *p, int off, const uint32_t *data, int sz);
int capn_write64(struct capn_list64 *p, int off, const uint64_t *data, int sz);

/* capn_new_* functions create a new object
 * datasz is in bytes, ptrs is # of pointers, sz is # of elements in the list
 * If capn_new_string sz < 0, strlen is used to compute the string length
 * On an error a CAPN_NULL pointer is returned
 */
struct capn_ptr capn_new_struct(struct capn_segment *seg, int datasz, int ptrs);
struct capn_ptr capn_new_list(struct capn_segment *seg, int sz, int datasz, int ptrs);
struct capn_ptr capn_new_bit_list(struct capn_segment *seg, int sz);
struct capn_ptr capn_new_ptr_list(struct capn_segment *seg, int sz);
struct capn_ptr capn_new_string(struct capn_segment *seg, const char *str, int sz);

/* capn_get|set_* functions get/set struct values
 * off is the offset into the structure in bytes
 * Rarely should these be called directly, instead use the generated code.
 * Data must be xored with the default value
 * These are static in order to be inlined.
 */
static uint8_t capn_get8(const struct capn_ptr *p, int off);
static uint16_t capn_get16(const struct capn_ptr *p, int off);
static uint32_t capn_get32(const struct capn_ptr *p, int off);
static uint64_t capn_get64(const struct capn_ptr *p, int off);
static int capn_set8(struct capn_ptr *p, int off, uint8_t val);
static int capn_set16(struct capn_ptr *p, int off, uint16_t val);
static int capn_set32(struct capn_ptr *p, int off, uint32_t val);
static int capn_set64(struct capn_ptr *p, int off, uint64_t val);




/* Inline functions */


#define T(IDX) s.v[IDX] = (uint8_t) (v >> (8*IDX))
static uint8_t capn_flip8(uint8_t v) {
	return v;
}
static uint16_t capn_flip16(uint16_t v) {
	union { uint16_t u; uint8_t v[2]; } s;
	T(0); T(1);
	return s.u;
}
static uint32_t capn_flip32(uint32_t v) {
	union { uint32_t u; uint16_t v[2]; } s;
	T(0); T(1); T(2); T(3);
	return s.u;
}
static uint64_t capn_flip64(uint64_t v) {
	union { uint64_t u; uint32_t v[2]; } s;
	T(0); T(1); T(2); T(3); T(4); T(5); T(6); T(7);
	return s.u;
}
#undef T

static uint8_t capn_get8(const struct capn_ptr *p, int off) {
	return (p->type == CAPN_STRUCT && off < p->datasz) ? capn_flip8(*(uint8_t*) p->data) : 0;
}
static int capn_set8(struct capn_ptr *p, int off, uint8_t val) {
	if (p->type == CAPN_STRUCT && off < p->datasz) {
		*(uint8_t*) p->data = capn_flip8(val);
		return 0;
	} else {
		return -1;
	}
}

static uint16_t capn_get16(const struct capn_ptr *p, int off) {
	return (p->type == CAPN_STRUCT && off < p->datasz) ? capn_flip16(*(uint16_t*) p->data) : 0;
}
static int capn_set16(struct capn_ptr *p, int off, uint16_t val) {
	if (p->type == CAPN_STRUCT && off < p->datasz) {
		*(uint16_t*) p->data = capn_flip16(val);
		return 0;
	} else {
		return -1;
	}
}

static uint32_t capn_get32(const struct capn_ptr *p, int off) {
	return (p->type == CAPN_STRUCT && off < p->datasz) ? capn_flip32(*(uint32_t*) p->data) : 0;
}
static int capn_set32(struct capn_ptr *p, int off, uint32_t val) {
	if (p->type == CAPN_STRUCT && off < p->datasz) {
		*(uint32_t*) p->data = capn_flip32(val);
		return 0;
	} else {
		return -1;
	}
}

static uint64_t capn_get64(const struct capn_ptr *p, int off) {
	return (p->type == CAPN_STRUCT && off < p->datasz) ? capn_flip64(*(uint64_t*) p->data) : 0;
}
static int capn_set64(struct capn_ptr *p, int off, uint64_t val) {
	if (p->type == CAPN_STRUCT && off < p->datasz) {
		*(uint64_t*) p->data = capn_flip64(val);
		return 0;
	} else {
		return -1;
	}
}

static float capn_get_float(const struct capn_ptr *p, int off, float def) {
	union { float f; uint32_t u;} u;
	u.f = def;
	u.u ^= capn_get32(p, off);
	return u.f;
}
static int capn_set_float(struct capn_ptr *p, int off, float f, float def) {
	union { float f; uint32_t u;} u;
	union { float f; uint32_t u;} d;
	u.f = f;
	d.f = def;
	return capn_set32(p, off, u.u ^ d.u);
}

static double capn_get_double(const struct capn_ptr *p, int off, double def) {
	union { double f; uint64_t u;} u;
	u.f = def;
	u.u ^= capn_get64(p, off);
	return u.f;
}
static int capn_set_double(struct capn_ptr *p, int off, double f, double def) {
	union { double f; uint64_t u;} u;
	union { double f; uint64_t u;} d;
	d.f = f;
	u.f = f;
	return capn_set64(p, off, u.u ^ d.u);
}

#endif
