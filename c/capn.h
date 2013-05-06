/* vim: set sw=8 ts=8 sts=8 noet: */

#ifndef CAPN_H
#define CAPN_H

#include <stdint.h>

#define CAPN_SEGID_LOCAL 0xFFFFFFFF

/* struct capn is a common structure shared between segments in the same
 * session/context so that far pointers between the segments will be created.
 *
 * lookup is used to lookup segments by id when derefencing a far pointer
 *
 * create is used to create or lookup an alternate segment that has at least
 * sz available (ie returned seg->len + sz <= seg->cap)
 *
 * Allocated segments must be zero initialized.
 *
 * create and lookup can be NULL if you don't need multiple segments and don't
 * want to support copying
 *
 * create is also used to allocate room for the copy tree with id ==
 * CAPN_SEGID_LOCAL. This data should be allocated in the local memory space
 *
 * seglist and copylist are linked lists which can be used to free up segments
 * on cleanup
 *
 * lookup, create, and user can be set by the user. Other values should be
 * zero initialized.
 */
struct capn {
	struct capn_segment *(*lookup)(void* /*user*/, uint32_t /*id */);
	struct capn_segment *(*create)(void* /*user*/, uint32_t /*id */, int /*sz*/);
	void *user;
	uint32_t segnum;
	struct capn_tree *copy;
	struct capn_tree *segtree, *lastseg;
	struct capn_segment *seglist;
	struct capn_segment *copylist;
};

struct capn_tree {
	struct capn_tree *left, *right, *parent;
	unsigned int red : 1;
};

/* struct capn_segment contains the information about a single segment.
 * capn should point to a struct capn that is shared between segments in the
 * same session
 * id specifies the segment id, used for far pointers
 * data specifies the segment data. This should not move after creation.
 * len specifies the current segment length. This should be 0 for a blank
 * segment.
 * cap specifies the segment capacity.
 * When creating new structures len will be incremented until it reaces cap,
 * at which point a new segment will be requested via capn->create.
 *
 * data, len, and cap must all by 8 byte aligned
 *
 * data, len, cap should all set by the user. Other values should be zero
 * initialized.
 */
struct capn_segment {
	struct capn_tree hdr;
	struct capn_segment *next;
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
	CAPN_LIST_MEMBER = 5,
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

union capn_iptr {
	struct capn_ptr c;
	uintptr_t u;
	void *p;
};

struct capn_ret_vt {
	void (*free)(void*);
};

struct capn_list1{struct capn_ptr p;};
struct capn_list8{struct capn_ptr p;};
struct capn_list16{struct capn_ptr p;};
struct capn_list32{struct capn_ptr p;};
struct capn_list64{struct capn_ptr p;};

/* capn_append_segment appends a segment to a session */
void capn_append_segment(struct capn*, struct capn_segment*);

/* capn_root returns a fake pointer that can be used to read/write the session
 * root object using capn_(read|write)_ptr at index 0. The root is the object
 * pointed to by a ptr at offset 0 in segment 0. This will allocate room for
 * the root if not already.
 */
struct capn_ptr capn_root(struct capn*);

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

/* capn_read_* functions read data from a list
 * The length of the list is given by p->size
 * off specifies how far into the list to start
 * sz indicates the number of elements to read
 * The function returns the number of elements read or -1 on an error.
 * off must be byte aligned for capn_read_1
 */
int capn_read1(const struct capn_list1 *p, int off, uint8_t *data, int sz);
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
int capn_write1(struct capn_list1 *p, int off, const uint8_t *data, int sz);
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

#if defined(__cplusplus) || (defined(__STDC_VERSION__) && __STDC_VERSION__ >= 199901L)
#define CAPN_INLINE inline
#else
#define CAPN_INLINE static
#endif

/* capn_get|set_* functions get/set struct values
 * off is the offset into the structure in bytes
 * Rarely should these be called directly, instead use the generated code.
 * Data must be xored with the default value
 * These are inlined
 */
CAPN_INLINE uint8_t capn_get8(const struct capn_ptr *p, int off);
CAPN_INLINE uint16_t capn_get16(const struct capn_ptr *p, int off);
CAPN_INLINE uint32_t capn_get32(const struct capn_ptr *p, int off);
CAPN_INLINE uint64_t capn_get64(const struct capn_ptr *p, int off);
CAPN_INLINE int capn_set8(struct capn_ptr *p, int off, uint8_t val);
CAPN_INLINE int capn_set16(struct capn_ptr *p, int off, uint16_t val);
CAPN_INLINE int capn_set32(struct capn_ptr *p, int off, uint32_t val);
CAPN_INLINE int capn_set64(struct capn_ptr *p, int off, uint64_t val);


/* capn_init_malloc inits the capn struct with a create function which
 * allocates segments on the heap using malloc
 *
 * capn_free_all frees all the segment headers and data created by the create
 * function setup by capn_init_malloc
 */
void capn_init_malloc(struct capn *c);
void capn_free_all(struct capn *c);

/* capn_stream encapsulates the needed fields for capn_(deflate|inflate) in a
 * similar manner to z_stream from zlib
 *
 * The user should set next_in, avail_in, next_out, avail_out to the
 * available in/out buffers before calling capn_(deflate|inflate).
 *
 * Other fields should be zero initialized.
 */
struct capn_stream {
	uint8_t *next_in;
	int avail_in;
	uint8_t *next_out;
	int avail_out;
	int zeros, raw;
};

#define CAPN_NEED_MORE_IN -1
#define CAPN_NEED_MORE_OUT -2

/* capn_deflate deflates a stream to the packed format
 * capn_inflate inflates a stream from the packed format
 *
 * They will return CAPN_NEED_MORE_(IN|OUT) as appropriate or 0 if the entire
 * input has been processed.
 */
int capn_deflate(struct capn_stream*);
int capn_inflate(struct capn_stream*);

int capn_marshal_iptr(const union capn_iptr*, struct capn_ptr*, int off);

/* Inline functions */


#define T(IDX) s.v[IDX] = (uint8_t) (v >> (8*IDX))
CAPN_INLINE uint8_t capn_flip8(uint8_t v) {
	return v;
}
CAPN_INLINE uint16_t capn_flip16(uint16_t v) {
	union { uint16_t u; uint8_t v[2]; } s;
	T(0); T(1);
	return s.u;
}
CAPN_INLINE uint32_t capn_flip32(uint32_t v) {
	union { uint32_t u; uint8_t v[4]; } s;
	T(0); T(1); T(2); T(3);
	return s.u;
}
CAPN_INLINE uint64_t capn_flip64(uint64_t v) {
	union { uint64_t u; uint8_t v[8]; } s;
	T(0); T(1); T(2); T(3); T(4); T(5); T(6); T(7);
	return s.u;
}
#undef T

CAPN_INLINE uint8_t capn_get8(const struct capn_ptr *p, int off) {
	return off < p->datasz ? capn_flip8(*(uint8_t*) p->data) : 0;
}
CAPN_INLINE int capn_set8(struct capn_ptr *p, int off, uint8_t val) {
	if (off < p->datasz) {
		*(uint8_t*) p->data = capn_flip8(val);
		return 0;
	} else {
		return -1;
	}
}

CAPN_INLINE uint16_t capn_get16(const struct capn_ptr *p, int off) {
	return off < p->datasz ? capn_flip16(*(uint16_t*) p->data) : 0;
}
CAPN_INLINE int capn_set16(struct capn_ptr *p, int off, uint16_t val) {
	if (off < p->datasz) {
		*(uint16_t*) p->data = capn_flip16(val);
		return 0;
	} else {
		return -1;
	}
}

CAPN_INLINE uint32_t capn_get32(const struct capn_ptr *p, int off) {
	return off < p->datasz ? capn_flip32(*(uint32_t*) p->data) : 0;
}
CAPN_INLINE int capn_set32(struct capn_ptr *p, int off, uint32_t val) {
	if (off < p->datasz) {
		*(uint32_t*) p->data = capn_flip32(val);
		return 0;
	} else {
		return -1;
	}
}

CAPN_INLINE uint64_t capn_get64(const struct capn_ptr *p, int off) {
	return off < p->datasz ? capn_flip64(*(uint64_t*) p->data) : 0;
}
CAPN_INLINE int capn_set64(struct capn_ptr *p, int off, uint64_t val) {
	if (off < p->datasz) {
		*(uint64_t*) p->data = capn_flip64(val);
		return 0;
	} else {
		return -1;
	}
}

CAPN_INLINE float capn_get_float(const struct capn_ptr *p, int off, float def) {
	union { float f; uint32_t u;} u;
	u.f = def;
	u.u ^= capn_get32(p, off);
	return u.f;
}
CAPN_INLINE int capn_set_float(struct capn_ptr *p, int off, float f, float def) {
	union { float f; uint32_t u;} u;
	union { float f; uint32_t u;} d;
	u.f = f;
	d.f = def;
	return capn_set32(p, off, u.u ^ d.u);
}

CAPN_INLINE double capn_get_double(const struct capn_ptr *p, int off, double def) {
	union { double f; uint64_t u;} u;
	u.f = def;
	u.u ^= capn_get64(p, off);
	return u.f;
}
CAPN_INLINE int capn_set_double(struct capn_ptr *p, int off, double f, double def) {
	union { double f; uint64_t u;} u;
	union { double f; uint64_t u;} d;
	d.f = f;
	u.f = f;
	return capn_set64(p, off, u.u ^ d.u);
}

#endif
