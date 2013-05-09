#include <gtest/gtest.h>

static int g_AddTag = 1;
#define ADD_TAG g_AddTag

#include "capn.c"
#include "capn-malloc.c"

template <int wordCount>
union AlignedData {
  uint8_t bytes[wordCount * 8];
  uint64_t words[wordCount];
};

class Session {
public:
  Session() {capn_init_malloc(&capn);}
  ~Session() {capn_free(&capn);}
  struct capn capn;
};





TEST(WireFormat, SimpleRawDataStruct) {
  AlignedData<2> data = {{
    // Struct ref, offset = 1, dataSize = 1, referenceCount = 0
    0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
    // Content for the data segment.
    0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef
  }};

  struct capn_segment seg;
  memset(&seg, 0, sizeof(seg));
  seg.data = (char*) data.bytes;
  seg.len = seg.cap = sizeof(data.bytes);

  struct capn ctx;
  memset(&ctx, 0, sizeof(ctx));
  capn_append_segment(&ctx, &seg);

  EXPECT_EQ(&seg, ctx.seglist);
  EXPECT_EQ(&seg, ctx.lastseg);
  EXPECT_EQ(&seg.hdr, ctx.segtree);
  EXPECT_EQ(1, ctx.segnum);
  EXPECT_EQ(0, seg.id);

  struct capn_ptr ptr = capn_get_root(&ctx);
  EXPECT_EQ(CAPN_STRUCT, ptr.type);
  EXPECT_EQ(8, ptr.datasz);
  EXPECT_EQ(0, ptr.ptrsz);

  EXPECT_EQ(UINT64_C(0xefcdab8967452301), capn_read64(ptr, 0));
  EXPECT_EQ(UINT64_C(0), capn_read64(ptr, 8));
  EXPECT_EQ(UINT32_C(0x67452301), capn_read32(ptr, 0));
  EXPECT_EQ(UINT32_C(0xefcdab89), capn_read32(ptr, 4));
  EXPECT_EQ(UINT32_C(0), capn_read32(ptr, 8));
  EXPECT_EQ(UINT16_C(0x2301), capn_read16(ptr, 0));
  EXPECT_EQ(UINT16_C(0x6745), capn_read16(ptr, 2));
  EXPECT_EQ(UINT16_C(0xab89), capn_read16(ptr, 4));
  EXPECT_EQ(UINT16_C(0xefcd), capn_read16(ptr, 6));
  EXPECT_EQ(UINT16_C(0), capn_read16(ptr, 8));
}

static const AlignedData<2> SUBSTRUCT_DEFAULT = {{0,0,0,0,1,0,0,0,  0,0,0,0,0,0,0,0}};
static const AlignedData<2> STRUCTLIST_ELEMENT_SUBSTRUCT_DEFAULT =
    {{0,0,0,0,1,0,0,0,  0,0,0,0,0,0,0,0}};

static void setupStruct(struct capn *ctx) {
  struct capn_ptr root = capn_new_root(ctx);
  ASSERT_EQ(CAPN_PTR_LIST, root.type);
  ASSERT_EQ(1, root.size);

  struct capn_ptr ptr = capn_new_struct(root.seg, 16, 4);
  ASSERT_EQ(CAPN_STRUCT, ptr.type);
  EXPECT_EQ(16, ptr.datasz);
  EXPECT_EQ(32, ptr.ptrsz);
  EXPECT_EQ(0, capn_setp(root, 0, ptr));

  EXPECT_EQ(0, capn_write64(ptr, 0, UINT64_C(0x1011121314151617)));
  EXPECT_EQ(0, capn_write32(ptr, 8, UINT32_C(0x20212223)));
  EXPECT_EQ(0, capn_write16(ptr, 12, UINT16_C(0x3031)));
  EXPECT_EQ(0, capn_write8(ptr, 14, 0x40));
  EXPECT_EQ(0, capn_write8(ptr, 15, (1 << 6) | (1 << 5) | (1 << 4) | (1 << 2)));

  struct capn_ptr subStruct = capn_new_struct(ptr.seg, 8, 0);
  ASSERT_EQ(CAPN_STRUCT, subStruct.type);
  EXPECT_EQ(8, subStruct.datasz);
  EXPECT_EQ(0, subStruct.ptrsz);
  EXPECT_EQ(0, capn_write32(subStruct, 0, 123));
  EXPECT_NE(0, capn_write32(subStruct, 8, 124));
  EXPECT_EQ(0, capn_setp(ptr, 0, subStruct));

  struct capn_ptr list = capn_new_list(ptr.seg, 3, 4, 0);
  ASSERT_EQ(CAPN_LIST, list.type);
  EXPECT_EQ(3, list.size);
  EXPECT_EQ(4, list.datasz);
  EXPECT_EQ(0, capn_set32(list, 0, 200));
  EXPECT_EQ(0, capn_set32(list, 1, 201));
  EXPECT_EQ(0, capn_set32(list, 2, 202));
  EXPECT_NE(0, capn_set32(list, 3, 203));
  EXPECT_NE(0, capn_set64(list, 0, 405));
  EXPECT_EQ(0, capn_setp(ptr, 1, list));

  list = capn_new_list(ptr.seg, 4, 4, 1);
  ASSERT_EQ(CAPN_LIST, list.type);
  EXPECT_EQ(4, list.size);
  EXPECT_EQ(8, list.datasz);
  EXPECT_EQ(8, list.ptrsz);
  EXPECT_EQ(0, capn_setp(ptr, 2, list));
  for (int i = 0; i < 4; i++) {
    struct capn_ptr element = capn_getp(list, i);
    ASSERT_EQ(CAPN_STRUCT, element.type);
    EXPECT_EQ(1, element.is_list_member);
    EXPECT_EQ(8, element.datasz);
    EXPECT_EQ(8, element.ptrsz);
    EXPECT_EQ(0, capn_write32(element, 0, 300+i));

    struct capn_ptr subelement = capn_new_struct(element.seg, 8, 0);
    ASSERT_EQ(CAPN_STRUCT, subelement.type);
    EXPECT_EQ(8, subelement.datasz);
    EXPECT_EQ(0, subelement.ptrsz);
    EXPECT_EQ(0, capn_write32(subelement, 0, 400+i));
    EXPECT_EQ(0, capn_setp(element, 0, subelement));
  }

  list = capn_new_ptr_list(ptr.seg, 5);
  ASSERT_EQ(CAPN_PTR_LIST, list.type);
  EXPECT_EQ(5, list.size);
  EXPECT_EQ(0, capn_setp(ptr, 3, list));
  for (int i = 0; i < 5; i++) {
    struct capn_ptr element = capn_new_list(list.seg, i+1, 2, 0);
    ASSERT_EQ(CAPN_LIST, element.type);
    EXPECT_EQ(i+1, element.size);
    EXPECT_EQ(2, element.datasz);
    EXPECT_EQ(0, element.ptrsz);
    EXPECT_EQ(0, capn_setp(list, i, element));
    for (int j = 0; j <= i; j++) {
      EXPECT_EQ(0, capn_set16(element, j, 500+j));
    }
  }
}

static void checkStruct(struct capn *ctx) {
  struct capn_ptr ptr = capn_get_root(ctx);
  EXPECT_EQ(CAPN_STRUCT, ptr.type);
  EXPECT_EQ(16, ptr.datasz);
  EXPECT_EQ(32, ptr.ptrsz);
  EXPECT_EQ(UINT64_C(0x1011121314151617), capn_read64(ptr, 0));
  EXPECT_EQ(UINT32_C(0x20212223), capn_read32(ptr, 8));
  EXPECT_EQ(0x3031, capn_read16(ptr, 12));
  EXPECT_EQ(0x40, capn_read8(ptr, 14));
  EXPECT_EQ((1 << 6) | (1 << 5) | (1 << 4) | (1 << 2), capn_read8(ptr, 15));

  struct capn_ptr subStruct = capn_getp(ptr, 0);
  EXPECT_EQ(CAPN_STRUCT, subStruct.type);
  EXPECT_EQ(8, subStruct.datasz);
  EXPECT_EQ(0, subStruct.ptrsz);
  EXPECT_EQ(123, capn_read32(subStruct, 0));

  struct capn_ptr list = capn_getp(ptr, 1);
  EXPECT_EQ(CAPN_LIST, list.type);
  EXPECT_EQ(3, list.size);
  EXPECT_EQ(4, list.datasz);
  EXPECT_EQ(0, list.ptrsz);
  EXPECT_EQ(200, capn_get32(list, 0));
  EXPECT_EQ(201, capn_get32(list, 1));
  EXPECT_EQ(202, capn_get32(list, 2));
  EXPECT_EQ(0, capn_get32(list, 3));
  EXPECT_EQ(0, capn_get64(list, 0));
  EXPECT_EQ(201, capn_get8(list, 1));
  EXPECT_EQ(202, capn_get16(list, 2));

  list = capn_getp(ptr, 2);
  EXPECT_EQ(CAPN_LIST, list.type);
  EXPECT_EQ(4, list.size);
  EXPECT_EQ(8, list.datasz);
  EXPECT_EQ(8, list.ptrsz);

  for (int i = 0; i < 4; i++) {
    struct capn_ptr element = capn_getp(list, i);
    EXPECT_EQ(CAPN_STRUCT, element.type);
    EXPECT_EQ(1, element.is_list_member);
    EXPECT_EQ(8, element.datasz);
    EXPECT_EQ(8, element.ptrsz);
    EXPECT_EQ(300+i, capn_read32(element,0));

    struct capn_ptr subelement = capn_getp(element, 0);
    EXPECT_EQ(CAPN_STRUCT, subelement.type);
    EXPECT_EQ(8, subelement.datasz);
    EXPECT_EQ(0, subelement.ptrsz);
    EXPECT_EQ(400+i, capn_read32(subelement, 0));
  }

  list = capn_getp(ptr, 3);
  EXPECT_EQ(CAPN_PTR_LIST, list.type);
  EXPECT_EQ(5, list.size);
  for (int i = 0; i < 5; i++) {
    struct capn_ptr element = capn_getp(list, i);
    EXPECT_EQ(CAPN_LIST, element.type);
    EXPECT_EQ(i+1, element.size);
    EXPECT_EQ(2, element.datasz);
    EXPECT_EQ(0, element.ptrsz);
    for (int j = 0; j <= i; j++) {
      EXPECT_EQ(500+j, capn_get16(element, j));
    }
  }
}

TEST(WireFormat, StructRoundTrip_OneSegment) {
  Session ctx;
  setupStruct(&ctx.capn);

  // word count:
  //    1  root reference
  //    6  root struct
  //    1  sub message
  //    2  3-element int32 list
  //   13  struct list
  //         1 tag
  //        12 4x struct
  //           1 data segment
  //           1 reference segment
  //           1 sub-struct
  //   11  list list
  //         5 references to sub-lists
  //         6 sub-lists (4x 1 word, 1x 2 words)
  // -----
  //   34
  ASSERT_EQ(1, ctx.capn.segnum);
  EXPECT_EQ(34*8, ctx.capn.seglist->len);

  checkStruct(&ctx.capn);

  struct capn ctx2;
  memset(&ctx2, 0, sizeof(ctx2));
  capn_append_segment(&ctx2, ctx.capn.seglist);
  checkStruct(&ctx2);
}

static struct capn_segment *CreateSmallSegment(void *u, uint32_t id, int sz) {
  struct capn_segment *s = (struct capn_segment*) calloc(1, sizeof(*s));
  s->data = (char*) calloc(1, sz);
  s->cap = sz;
  return s;
}

static void getSegments(struct capn *c, struct capn_segment **s, size_t num) {
  ASSERT_EQ(num, c->segnum);

  s[0] = c->seglist;
  for (int i = 1; i < num; i++) {
    s[i] = s[i-1]->next;
  }

  for (int i = 0; i < num; i++) {
    EXPECT_EQ(s[i]->id, i);
  }
}

TEST(WireFormat, StructRoundTrip_OneSegmentPerAllocation) {
  Session ctx;
  ctx.capn.create = &CreateSmallSegment;

  setupStruct(&ctx.capn);

  struct capn_segment *segments[15];
  getSegments(&ctx.capn, segments, 15);

  // Check that each segment has the expected size.  Recall that the first word of each segment will
  // actually be a reference to the first thing allocated within that segment.
  EXPECT_EQ( 8, segments[ 0]->len);  // root ref
  EXPECT_EQ(56, segments[ 1]->len);  // root struct
  EXPECT_EQ(16, segments[ 2]->len);  // sub-struct
  EXPECT_EQ(24, segments[ 3]->len);  // 3-element int32 list
  EXPECT_EQ(80, segments[ 4]->len);  // struct list
  EXPECT_EQ(16, segments[ 5]->len);  // struct list substruct 1
  EXPECT_EQ(16, segments[ 6]->len);  // struct list substruct 2
  EXPECT_EQ(16, segments[ 7]->len);  // struct list substruct 3
  EXPECT_EQ(16, segments[ 8]->len);  // struct list substruct 4
  EXPECT_EQ(48, segments[ 9]->len);  // list list
  EXPECT_EQ(16, segments[10]->len);  // list list sublist 1
  EXPECT_EQ(16, segments[11]->len);  // list list sublist 2
  EXPECT_EQ(16, segments[12]->len);  // list list sublist 3
  EXPECT_EQ(16, segments[13]->len);  // list list sublist 4
  EXPECT_EQ(24, segments[14]->len);  // list list sublist 5

  checkStruct(&ctx.capn);

  struct capn ctx2;
  memset(&ctx2, 0, sizeof(ctx2));
  for (int i = 0; i < 15; i++) {
    capn_append_segment(&ctx2, segments[i]);
  }

  checkStruct(&ctx2);
}

TEST(WireFormat, StructRoundTrip_OneSegmentPerAllocation_NoTag) {
  Session ctx;
  ctx.capn.create = &CreateSmallSegment;

  g_AddTag = 0;
  setupStruct(&ctx.capn);
  g_AddTag = 1;

  struct capn_segment *segments[29];
  getSegments(&ctx.capn, segments, 29);

  // Check that each segment has the expected size. Note that we have plenty
  // of 16 byte double far ptrs.
  EXPECT_EQ( 8, segments[ 0]->len);  // root ref
  EXPECT_EQ(48, segments[ 1]->len);  // root struct
  EXPECT_EQ(16, segments[ 2]->len);  // root struct ptr
  EXPECT_EQ( 8, segments[ 3]->len);  // sub-struct
  EXPECT_EQ(16, segments[ 4]->len);  // sub-struct ptr
  EXPECT_EQ(16, segments[ 5]->len);  // 3-element int32 list
  EXPECT_EQ(16, segments[ 6]->len);  // 3-element int32 list ptr
  EXPECT_EQ(72, segments[ 7]->len);  // struct list
  EXPECT_EQ(16, segments[ 8]->len);  // struct list ptr
  EXPECT_EQ( 8, segments[ 9]->len);  // struct list substruct 1
  EXPECT_EQ(16, segments[10]->len);  // struct list substruct 1 ptr
  EXPECT_EQ( 8, segments[11]->len);  // struct list substruct 2
  EXPECT_EQ(16, segments[12]->len);  // struct list substruct 2 ptr
  EXPECT_EQ( 8, segments[13]->len);  // struct list substruct 3
  EXPECT_EQ(16, segments[14]->len);  // struct list substruct 3 ptr
  EXPECT_EQ( 8, segments[15]->len);  // struct list substruct 4
  EXPECT_EQ(16, segments[16]->len);  // struct list substruct 4 ptr
  EXPECT_EQ(40, segments[17]->len);  // list list
  EXPECT_EQ(16, segments[18]->len);  // list list ptr
  EXPECT_EQ( 8, segments[19]->len);  // list list sublist 1
  EXPECT_EQ(16, segments[20]->len);  // list list sublist 1 ptr
  EXPECT_EQ( 8, segments[21]->len);  // list list sublist 2
  EXPECT_EQ(16, segments[22]->len);  // list list sublist 2 ptr
  EXPECT_EQ( 8, segments[23]->len);  // list list sublist 3
  EXPECT_EQ(16, segments[24]->len);  // list list sublist 3 ptr
  EXPECT_EQ( 8, segments[25]->len);  // list list sublist 4
  EXPECT_EQ(16, segments[26]->len);  // list list sublist 4 ptr
  EXPECT_EQ(16, segments[27]->len);  // list list sublist 5
  EXPECT_EQ(16, segments[28]->len);  // list list sublist 5 ptr

  checkStruct(&ctx.capn);

  struct capn ctx2;
  memset(&ctx2, 0, sizeof(ctx2));
  for (int i = 0; i < 29; i++) {
    capn_append_segment(&ctx2, segments[i]);
  }

  checkStruct(&ctx2);
}


static struct capn_segment *CreateSegment64(void *u, uint32_t id, int sz) {
  if (sz < 64) {
    sz = 64;
  }
  struct capn_segment *s = (struct capn_segment*) calloc(1, sizeof(*s));
  s->data = (char*) calloc(1, sz);
  s->cap = sz;
  return s;
}

TEST(WireFormat, StructRoundTrip_MultipleSegmentsWithMultipleAllocations) {
  Session ctx;
  ctx.capn.create = &CreateSegment64;

  setupStruct(&ctx.capn);

  // Verify that we made 6 segments.
  ASSERT_EQ(6, ctx.capn.segnum);

  struct capn_segment *segments[6];
  segments[0] = ctx.capn.seglist;
  for (int i = 1; i < 6; i++) {
    segments[i] = segments[i-1]->next;
  }

  for (int i = 0; i < 6; i++) {
    EXPECT_EQ(segments[i]->id, i);
  }

  // Check that each segment has the expected size.  Recall that each object will be prefixed by an
  // extra word if its parent is in a different segment.
  EXPECT_EQ(64, segments[0]->len);  // root ref + struct + sub
  EXPECT_EQ(56, segments[1]->len);  // 3-element int32 list (+tag) + struct list substructs 1+2 (+tag)
  EXPECT_EQ(80, segments[2]->len);  // struct list (+tag)
  EXPECT_EQ(64, segments[3]->len);  // struct list substructs 3+4 (+tag) + list list sublist 1,2 (+tag)
  EXPECT_EQ(64, segments[4]->len);  // list list (+tag) + sublist 3,4
  EXPECT_EQ(24, segments[5]->len);  // list list sublist 5

  checkStruct(&ctx.capn);

  struct capn ctx2;
  memset(&ctx2, 0, sizeof(ctx2));
  for (int i = 0; i < 6; i++) {
    capn_append_segment(&ctx2, segments[i]);
  }

  checkStruct(&ctx2);
}

int main(int argc, char *argv[]) {
  testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

