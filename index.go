// Package index provides functionality for accessing compact arrays with
// spaced indices with reduced garbage generation and faster insertions and
// deletions.
package index

// I is a handle to a value in an index.
type I int

// DefaultSize is the size of indices initialised with Init.
const DefaultSize = 1024

// An Index is an index with a holed array pointing virtual index.I's into the
// underlying data array (which is specified by the using system).
type Index struct {
    I       []int // holed
    Reverse []I   // compact
    Empty   []I   // stack
}

// Init initialises an Index in-place allocating DefaultSize slots.
func (idx *Index) Init() {
    idx.InitSize(DefaultSize)
}

// InitSize initialises an Index in-place allocating size slots.
func (idx *Index) InitSize(size int) {
    idx.I = make([]int, 0, size)
    idx.Reverse = make([]I, 0, size)
    idx.Empty = make([]I, 0, size/8)
}

// New returns the next I corresponding to the value just appended to the data.
// NOT THREAD-SAFE!
func (idx *Index) Add() (i I) {
    if len(idx.Empty) != 0 {
        i = I(idx.Empty[len(idx.Empty)-1])
        idx.Empty = idx.Empty[:len(idx.Empty)-1]
        idx.I[i] = len(idx.Reverse)
    } else {
        i = I(len(idx.I))
        idx.I = append(idx.I, len(idx.Reverse))
    }
    idx.Reverse = append(idx.Reverse, i)
    return
}

// Remove removes the data pointed to with i from the index. The caller should
// have replaced the index (idx.I[i]) with the end of the data before calling
// Remove.
// NOT THREAD-SAFE!
func (idx *Index) Remove(i I) {
    j := idx.I[i]
    // idx.I[i] = -1 // should not be necessary...
    idx.Empty = append(idx.Empty, i)
    lastj := len(idx.Reverse) - 1
    if j != lastj {
        i = idx.Reverse[lastj]
        idx.Reverse[j], idx.I[i] = i, j
    }
    idx.Reverse = idx.Reverse[:lastj]
}

// Swap swaps the slots i and j ensuring the reverse table is up-to-date as
// well.
func (idx *Index) Swap(i, j I) {
    a, b := idx.I[i], idx.I[j]
    idx.I[i], idx.I[j] = b, a
    idx.Reverse[a], idx.Reverse[b] = j, i
}

// SwapBacking swaps the slots with the backing indices i and j ensuring the
// direct table is up-to-date as well.
func (idx *Index) SwapBacking(i, j int) {
    a, b := idx.Reverse[i], idx.Reverse[j]
    idx.Reverse[i], idx.Reverse[j] = b, a
    idx.I[a], idx.I[b] = j, i
}

// Clears the entire index, removing everything. This does not deallocate data.
func (idx *Index) Clear() {
    idx.I = idx.I[:0]
    idx.Reverse = idx.Reverse[:0]
    idx.Empty = idx.Empty[:0]
}
