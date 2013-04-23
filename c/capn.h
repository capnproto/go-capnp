/* vim: set sw=8 ts=8 sts=8 noet: */

#ifndef CAPN_H
#define CAPN_H

#include <stdint.h>

enum CAPN_TYPE {
	CAPN_NULL = 0,
	CAPN_FAR = 2,
	CAPN_DOUBLE_FAR = 6,
	CAPN_STRUCT = 7,
	CAPN_VOID_LIST = 8,
	CAPN_1BIT_LIST = 9,
	CAPN_1BYTE_LIST = 10,
	CAPN_2BYTE_LIST = 11,
	CAPN_4BYTE_LIST = 12,
	CAPN_8BYTE_LIST = 13,
	CAPN_POINTER_LIST = 14,
	CAPN_COMPOSITE_LIST = 15,
};

struct capn_ptr {
	/* should not be used unless capn_deref_ptr has been called first */
	enum CAPN_TYPE type;
	int elements;

	/* rest is private */
	union {
		struct {
			uint16_t datasz;
			uint16_t ptrsz;
		} s;
		uint32_t id;
	} u;
	struct capn_segment *seg;
	char *data;
};

int capn_deref_ptr(struct capn_ptr*);

/* length is optional and can be set to NULL */
char *capn_to_string(struct capn_ptr*, int *length);

/* capn_*_struct reads/writes data to a from a struct. These should not be
 * called directly, but read|write_TYPE should be called instead (created by
 * the generator).
 *
 * WARNING: Fields in structs are always in little endian format. Use
 * capn_to_* when dealing with struct fields.
 */
void capn_read_struct(struct capn_ptr*, void *p, int data, int ptrs);
void capn_write_struct(struct capn_ptr*, const void *p, int data, int ptrs);

/* capn_read_* functions are read data from a list.
 *
 * The length of the list is given by capn_elements. off specifies how far
 * into the list to start, sz indicates the number of elements to read.
 *
 * With capn_read_1 off and sz are in bytes (_not_ bits).
 * Bits read are in little endian (lowest bit first) order.
 *
 * If off+sz > elements in list, then only elements up to the end of the list
 * will be read.
 *
 * Data read is in native byte order.
 */
void capn_read_1(struct capn_ptr*, int off, uint8_t*, int sz);
void capn_read_8(struct capn_ptr*, int off, uint8_t*, int sz);
void capn_read_16(struct capn_ptr*, int off, uint16_t*, int sz);
void capn_read_32(struct capn_ptr*, int off, uint32_t*, int sz);
void capn_read_64(struct capn_ptr*, int off, uint64_t*, int sz);
struct capn_ptr capn_read_ptr(struct capn_ptr*, int off);

/* capn_write_* functions write data to a list.
 *
 * Only up to the end of the list is written as given by capn_elemnts.
 *
 * With capn_write_1 off and sz are in bytes (_not_ bits). Bits are presented
 * in little endian (lowest bit and lowest byte first).
 *
 * Data provided should be in native byte order.
 */
void capn_write_1(struct capn_ptr*, int off, const void*, int sz);
void capn_write_8(struct capn_ptr*, int off, const uint8_t*, int sz);
void capn_write_16(struct capn_ptr*, int off, const uint16_t*, int sz);
void capn_write_32(struct capn_ptr*, int off, const uint32_t*, int sz);
void capn_write_64(struct capn_ptr*, int off, const uint64_t*, int sz);
int capn_write_ptr(struct capn_ptr*, int off, struct capn_ptr *to);

struct capn_ptr capn_new_struct(struct capn_segment*, int datasz, int ptrs, int want_tag);
struct capn_ptr capn_new_list(struct capn_segment*, enum CAPN_TYPE, int sz, int want_tag);
struct capn_ptr capn_new_composite(struct capn_segment*, int elems, int datasz, int ptrs, int want_tag);
/* use sz == -1 for null terminated string */
struct capn_ptr capn_new_string(struct capn_segment*, const char *s, int sz, int want_tag);

int capn_copy(struct capn_ptr *to, struct capn_ptr *from);

typedef struct capn_segment *(*capn_create_t)(void* /*user*/, int /*sz*/);
typedef struct capn_segment *(*capn_lookup_t)(void* /*user*/, uint32_t /*id*/);

struct capn *capn_new(void);
void capn_free(struct capn*);
void capn_add(struct capn*, struct capn_segment*);
void capn_set_lookup(struct capn*, capn_lookup_t lookup, void *user);
void capn_set_create(struct capn*, capn_create_t create, void *user);

struct capn_segment *capn_new_segment(void *data, int len, int cap, void (*free)(void*));

#define T(IDX) s.v[IDX] = (uint8_t) (v >> (8*IDX))
#define F(SZ, IDX) ((uint ## SZ ## _t) (s.v[IDX]) << (IDX*8))
static uint16_t capn_to_le_16(uint16_t v) {
	union { uint16_t u; uint8_t v[2]; } s;
	T(0); T(1);
	return s.u;
}
static uint32_t capn_to_le_32(uint32_t v) {
	union { uint32_t u; uint16_t v[2]; } s;
	T(0); T(1); T(2); T(3);
	return s.u;
}
static uint64_t capn_to_le_64(uint64_t v) {
	union { uint64_t u; uint32_t v[2]; } s;
	T(0); T(1); T(2); T(3); T(4); T(5); T(6); T(7);
	return s.u;
}
static uint16_t capn_from_le_16(uint16_t v) {
	union { uint16_t u; uint8_t v[2]; } s;
	s.u = v;
	return F(16, 1) | F(16, 0);
}
static uint32_t capn_from_le_32(uint32_t v) {
	union { uint32_t u; uint8_t v[4]; } s;
	s.u = v;
	return F(32, 3) | F(32, 2) | F(32, 1) | F(32, 0);
}
static uint64_t capn_from_le_64(uint64_t v) {
	union { uint64_t u; uint8_t v[8]; } s;
	s.u = v;
	return F(64,7) | F(64,6) | F(64,5) | F(64,4) | F(64,3) | F(64,2) | F(64,1) | F(64,0);
}
#undef T
#undef F

#endif
