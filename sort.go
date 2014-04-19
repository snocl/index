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
    s.Index.SwapBacking(i, j)
}
