package index

import (
    "sort"
)

// A Sorter ensures that sorting on an index also updates the index maps.
type Sorter struct {
    Index *Index
    sort.Interface
}

func (s Sorter) Swap(i, j int) {
    s.Interface.Swap(i, j)
    a, b := s.Index.Reverse[i], s.Index.Reverse[j]
    s.Index.Reverse[i], s.Index.Reverse[j] = b, a
    s.Index.I[a], s.Index.I[b] = j, i
}
