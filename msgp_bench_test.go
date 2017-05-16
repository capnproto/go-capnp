package capnp_test

import (
	"math/rand"
	"testing"
	"time"
)

//go:generate msgp -tests=false -o msgp_bench_gen_test.go
//msgp:Tuple Event

type Event struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func BenchmarkUnmarshalMsgp(b *testing.B) {
	r := rand.New(rand.NewSource(12345))
	data := make([][]byte, 1000)
	for i := range data {
		msg, _ := (*Event)(generateA(r)).MarshalMsg(nil)
		data[i] = msg
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var e Event
		msg := data[r.Intn(len(data))]
		_, err := e.UnmarshalMsg(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
