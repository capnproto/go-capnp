package bufferpool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBucketIndex(t *testing.T) {
	tests := []struct {
		size int
		get  int
		put  int
	}{
		// Only sizes that are powers of two are obtained and returned
		// to the same bucket.
		//
		// Sizes that are not a power of two must be fetched by the next
		// higher power of two, but are returned to the lower one.
		{size: 0, get: 0, put: -1},
		{size: 1, get: 0, put: 0},
		{size: 26, get: 5, put: 4}, // 26 == 0b00011010
		{size: 32, get: 5, put: 5},
		{size: 1024, get: 10, put: 10},
		{size: 1025, get: 11, put: 10},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(fmt.Sprintf("%d", tc.size), func(t *testing.T) {
			get := bucketToGet(tc.size)
			require.Equal(t, tc.get, get)
			put := bucketToPut(tc.size)
			require.Equal(t, tc.put, put)
		})
	}
}

func TestBucketSlice(t *testing.T) {
	const minAlloc = 8
	const bucketCount = 10
	const sizeLastBucket = 1 << (bucketCount - 1)

	tests := []struct {
		size    int
		wantLen int
		wantCap int
	}{{
		size:    -1, // Negative values are skipped.
		wantLen: 0,
		wantCap: 0,
	}, {
		size:    0,
		wantLen: minAlloc,
		wantCap: minAlloc,
	}, {
		size:    1,
		wantLen: minAlloc,
		wantCap: minAlloc,
	}, {
		size:    minAlloc,
		wantLen: minAlloc,
		wantCap: minAlloc,
	}, {
		size:    minAlloc + 1, // Goes to next bucket.
		wantLen: minAlloc * 2,
		wantCap: minAlloc * 2,
	}, {
		size:    minAlloc*2 + 1,
		wantLen: minAlloc * 4,
		wantCap: minAlloc * 4,
	}, {
		size:    sizeLastBucket - 1,
		wantLen: sizeLastBucket,
		wantCap: sizeLastBucket,
	}, {
		size:    sizeLastBucket,
		wantLen: sizeLastBucket,
		wantCap: sizeLastBucket,
	}, {
		size:    sizeLastBucket + 1, // Anything > last bucket size is not allocated.
		wantLen: 0,
		wantCap: 0,
	}}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d", tc.size), func(t *testing.T) {
			bs := makeBucketSlice(minAlloc, bucketCount)
			require.Len(t, bs, bucketCount)
			buf := bs.Get(tc.size)
			require.Len(t, buf, tc.wantLen)
			require.Equal(t, tc.wantCap, cap(buf))
		})
	}
}
