#include "capn-stream.c"
#include <gtest/gtest.h>

template <int wordCount>
union AlignedData {
  uint8_t bytes[wordCount * 8];
  uint64_t words[wordCount];
};

TEST(Stream, ReadEmptyStream_Even) {
  AlignedData<2> data = {{
    1, 0, 0, 0, // num of segs - 1
    0, 0, 0, 0,
    0, 0, 0, 0,
    2, 3, 4, 0, // garbage that should be ignored
  }};

  struct capn ctx;
  ASSERT_NE(0, capn_init_mem(&ctx, data.bytes, 12, 0));
  ASSERT_EQ(0, capn_init_mem(&ctx, data.bytes, 16, 0));
  EXPECT_EQ(2, ctx.segnum);
  EXPECT_EQ(0, ctx.seglist->len);
  EXPECT_EQ(0, ctx.seglist->next->len);
  capn_free(&ctx);
}

TEST(Stream, ReadEmptyStream_Odd) {
  AlignedData<3> data = {{
    2, 0, 0, 0, // num of segs - 1
    0, 0, 0, 0,
    0, 0, 0, 0,
    0, 0, 0, 0,
    2, 3, 4, 0, // garbage that should be ignored
  }};

  struct capn ctx;
  ASSERT_NE(0, capn_init_mem(&ctx, data.bytes, 12, 0));

  ASSERT_EQ(0, capn_init_mem(&ctx, data.bytes, 16, 0));
  EXPECT_EQ(3, ctx.segnum);
  EXPECT_EQ(0, ctx.seglist->len);
  EXPECT_EQ(0, ctx.seglist->next->len);
  capn_free(&ctx);

  ASSERT_EQ(0, capn_init_mem(&ctx, data.bytes, 20, 0));
  EXPECT_EQ(3, ctx.segnum);
  EXPECT_EQ(0, ctx.seglist->len);
  EXPECT_EQ(0, ctx.seglist->next->len);
  capn_free(&ctx);
}

TEST(Stream, ReadStream_Even) {
  AlignedData<5> data = {{
     1, 0, 0, 0, // num of segs - 1
     1, 0, 0, 0,
     2, 0, 0, 0,
     2, 3, 4, 0, // garbage that should be ignored
     1, 2, 3, 4, 5, 6, 7, 8,
     9,10,11,12,13,14,15,16,
    17,18,19,20,21,22,23,24,
  }};

  struct capn ctx;
  ASSERT_NE(0, capn_init_mem(&ctx, data.bytes, 36, 0));
  ASSERT_EQ(0, capn_init_mem(&ctx, data.bytes, 40, 0));
  EXPECT_EQ(2, ctx.segnum);
  EXPECT_EQ(8, ctx.seglist->len);
  EXPECT_EQ(1, ctx.seglist->data[0]);
  EXPECT_EQ(16, ctx.seglist->next->len);
  EXPECT_EQ(9, ctx.seglist->next->data[0]);
  capn_free(&ctx);
}
