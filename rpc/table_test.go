package rpc

import (
	"reflect"
	"testing"
)

type testTableID uint32

func TestLocalTable(t *testing.T) {
	var table localTable[testTableID, string]
	if id := table.Add("zero"); id != 0 {
		t.Fatalf("first Add ID = %d; want 0", id)
	}
	id1 := table.Add("one")
	id2 := table.Add("two")
	if got, ok := table.Remove(id1); !ok || got != "one" {
		t.Fatalf("Remove(%d) = %q, %t; want one, true", id1, got, ok)
	}
	if _, ok := table.Remove(id1); ok {
		t.Fatalf("second Remove(%d) succeeded", id1)
	}
	if id := table.Add("replacement"); id != id1 {
		t.Fatalf("replacement ID = %d; want lowest free ID %d", id, id1)
	}
	if got, ok := table.Find(id2); !ok || got != "two" {
		t.Fatalf("Find(%d) = %q, %t; want two, true", id2, got, ok)
	}

	var got []string
	table.Range(func(_ testTableID, value string) bool {
		got = append(got, value)
		return true
	})
	if want := []string{"zero", "replacement", "two"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Range values = %q; want %q", got, want)
	}

	cleared := table.Clear()
	if !reflect.DeepEqual(cleared, got) {
		t.Fatalf("Clear() = %q; want %q", cleared, got)
	}
	if id := table.Add("new zero"); id != 0 {
		t.Fatalf("ID after Clear = %d; want 0", id)
	}
}

func TestLocalTableReservedID(t *testing.T) {
	var table localTable[testTableID, string]
	id := table.Add("reserved")
	if _, ok := table.take(id); !ok {
		t.Fatal("take failed")
	}
	if got := table.Add("next"); got != 1 {
		t.Fatalf("Add while ID reserved = %d; want 1", got)
	}
	table.release(id)
	if got := table.Add("reused"); got != id {
		t.Fatalf("Add after release = %d; want %d", got, id)
	}
}

func TestLocalTableExhaustion(t *testing.T) {
	table := localTable[testTableID, string]{next: ^uint32(0)}
	defer func() {
		if recover() == nil {
			t.Fatal("Add did not panic at ID exhaustion")
		}
	}()
	table.Add("overflow")
}

func TestRemoteTable(t *testing.T) {
	var table remoteTable[testTableID, string]
	if !table.Create(100, "first") {
		t.Fatal("initial Create failed")
	}
	if table.Create(100, "duplicate") {
		t.Fatal("duplicate Create succeeded")
	}
	if got, ok := table.Find(100); !ok || got != "first" {
		t.Fatalf("Find(100) = %q, %t; want first, true", got, ok)
	}
	got, created := table.FindOrCreate(100, func() string { return "unused" })
	if created || got != "first" {
		t.Fatalf("FindOrCreate(existing) = %q, %t; want first, false", got, created)
	}
	got, created = table.FindOrCreate(1, func() string { return "second" })
	if !created || got != "second" {
		t.Fatalf("FindOrCreate(new) = %q, %t; want second, true", got, created)
	}
	if got, ok := table.Remove(100); !ok || got != "first" {
		t.Fatalf("Remove(100) = %q, %t; want first, true", got, ok)
	}
	if _, ok := table.Remove(100); ok {
		t.Fatal("second Remove(100) succeeded")
	}
	if !table.Create(1_000_000, "sparse") {
		t.Fatal("sparse Create failed")
	}
	values := make(map[testTableID]string)
	table.Range(func(id testTableID, value string) bool {
		values[id] = value
		return true
	})
	if want := map[testTableID]string{1: "second", 1_000_000: "sparse"}; !reflect.DeepEqual(values, want) {
		t.Fatalf("Range values = %v; want %v", values, want)
	}
	cleared := table.Clear()
	if len(cleared) != 2 {
		t.Fatalf("Clear() returned %d entries; want 2", len(cleared))
	}
}

func TestUintSet(t *testing.T) {
	var set uintSet
	for _, id := range []uint{0, 1, 63, 64, 65, 127, 128, 256} {
		set.add(id)
	}
	for _, id := range []uint{0, 1, 63, 64, 65, 127, 128, 256} {
		if !set.has(id) {
			t.Errorf("has(%d) = false; want true", id)
		}
	}
	set.remove(0)
	set.remove(64)
	set.remove(256)
	if got, ok := set.min(); !ok || got != 1 {
		t.Fatalf("min() = %d, %t; want 1, true", got, ok)
	}
	for _, id := range []uint{0, 64, 256, ^uint(0)} {
		if set.has(id) {
			t.Errorf("has(%d) = true; want false", id)
		}
	}
}
