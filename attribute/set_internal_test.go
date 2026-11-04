// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestBinaryInsertionSortKVs(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	keyFuncs := map[string]func(n, i int) int{
		"Sorted":   func(_, i int) int { return i },
		"Reversed": func(n, i int) int { return n - i },
		"Equal":    func(int, int) int { return 0 },
		// Few distinct keys, so most inputs have duplicates to check stability.
		"Random": func(n, _ int) int { return r.IntN(n/2 + 1) },
	}
	for name, key := range keyFuncs {
		t.Run(name, func(t *testing.T) {
			for n := range 2 * maxInsertionSort {
				kvs := make([]KeyValue, n)
				for i := range kvs {
					kvs[i] = Int(fmt.Sprintf("k%03d", key(n, i)), i)
				}
				checkSortMatchesStdlib(t, kvs)
			}
		})
	}
}

func FuzzBinaryInsertionSortKVs(f *testing.F) {
	f.Add([]byte("hello, world"))
	f.Add([]byte{3, 2, 1, 0, 3, 2, 1, 0})
	f.Fuzz(func(t *testing.T, keys []byte) {
		// One byte per key, so the fuzzer can reorder and repeat keys freely.
		kvs := make([]KeyValue, len(keys))
		for i, k := range keys {
			kvs[i] = Int(string(k), i)
		}
		checkSortMatchesStdlib(t, kvs)
	})
}

// checkSortMatchesStdlib fails unless binaryInsertionSortKVs sorts kvs exactly
// like stdlibSortKVs. Values hold original positions, so stability is checked too.
func checkSortMatchesStdlib(t *testing.T, kvs []KeyValue) {
	t.Helper()
	want := slices.Clone(kvs)
	stdlibSortKVs(want)
	binaryInsertionSortKVs(kvs)
	if !slices.Equal(kvs, want) {
		t.Fatalf("n=%d: binaryInsertionSortKVs = %v, want %v", len(kvs), kvs, want)
	}
}
