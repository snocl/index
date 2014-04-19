package index_test

import (
    "github.com/snorredc/index"
    "sort"
    "strings"
    "testing"
)

type SortByValue struct {
    *System
}

func (s SortByValue) Len() int {
    return len(s.Values)
}

func (s SortByValue) Less(i, j int) bool {
    return s.Values[i] < s.Values[j]
}

func (s SortByValue) Swap(i, j int) {
    s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
    s.Names[i], s.Names[j] = s.Names[j], s.Names[i]
}

func isSorted(s *System) bool {
    p := -1
    for _, v := range s.Values {
        if p >= v {
            return false
        }
        p = v
    }
    return true
}

func TestSort(t *testing.T) {
    s := NewSystem(7)

    text := "This sentence consists of many ordered parts."
    words := strings.Fields(text)
    randomish := [7]int{2, 3, 5, 1, 0, 4, 6}
    var is [7]index.I

    for i, j := range randomish {
        is[i] = s.Add(words[j], j)
    }
    if isSorted(s) {
        t.Errorf("not unsorted")
    }

    sort.Sort(index.Sorter{&s.index, SortByValue{s}})

    if !isSorted(s) {
        t.Errorf("sorting failed")
    }
    if newtext := strings.Join(s.Names[:], " "); newtext != text {
        t.Errorf("new text doesn't match old text")
    }
}
