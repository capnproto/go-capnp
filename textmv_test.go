package capn_test

import (
	"math/rand"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
)

// highlight how much faster text movement between segments
// is when special casing Text and Data
//
// run this test with capn.go:1334-1341 commented in/out to compare.
//
func BenchmarkTextMovementBetweenSegments(b *testing.B) {

	buf := make([]byte, 1<<21)
	buf2 := make([]byte, 1<<21)

	text := make([]byte, 1<<20)
	for i := range text {
		text[i] = byte(65 + rand.Int()%26)
	}
	//stext := string(text)
	//fmt.Printf("text = %#v\n", stext)

	astr := make([]string, 1000)
	for i := range astr {
		astr[i] = string(text[i*1000 : (i+1)*1000])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		seg := capn.NewBuffer(buf[:0])
		scratch := capn.NewBuffer(buf2[:0])

		ht := air.NewRootHoldsText(seg)
		tl := scratch.NewTextList(1000)

		for j := 0; j < 1000; j++ {
			tl.Set(j, astr[j])
		}

		ht.SetLst(tl)

	}
}
