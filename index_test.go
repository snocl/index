package index_test

import (
    "github.com/snorredc/index"
    "math/rand"
    "strings"
    "testing"
)

type System struct {
    index  index.Index
    Names  []string
    Values []int
}

func NewSystem(size int) *System {
    s := &System{
        Names:  make([]string, 0, size),
        Values: make([]int, 0, size),
    }
    s.index.InitSize(size)
    return s
}

func (s *System) Add(n string, v int) index.I {
    s.Names = append(s.Names, n)
    s.Values = append(s.Values, v)
    return s.index.Add()
}

func (s *System) Remove(i index.I) {
    l := len(s.Names) - 1
    s.Names[s.index.I[i]] = s.Names[l]
    s.Names = s.Names[:l]
    s.Values[s.index.I[i]] = s.Values[l]
    s.Values = s.Values[:l]
    s.index.Remove(i)
}

func (s *System) Clear() {
    s.Names = s.Names[:0]
    s.Values = s.Values[:0]
    s.index.Clear()
}

func (s *System) Name(i index.I) string {
    return s.Names[s.index.I[i]]
}

func (s *System) Value(i index.I) int {
    return s.Values[s.index.I[i]]
}

func (s *System) JoinNames(sep string) string {
    return strings.Join(s.Names, sep)
}

func (s *System) Each(f func(string) bool) bool {
    for _, n := range s.Names {
        if f(n) {
            return true
        }
    }
    return false
}

func CheckValid(t *testing.T, s *System) {
    if len(s.Names) != len(s.Values) || len(s.Names) != len(s.index.Reverse) {
        t.Error("slices not the same length")
    }
    for j, i := range s.index.Reverse {
        if s.index.I[i] != j {
            t.Error("index broken")
        }
    }
}

func (s *System) CheckAdd(t *testing.T, n string, v int, os ...index.I) index.I {
    i := s.Add(n, v)
    for _, o := range os {
        if i == o {
            t.Error("duplicate index")
        }
    }
    return i
}

func TestInsert(t *testing.T) {
    s := NewSystem(1024)

    i := s.Add("Hello", 42)
    for x := 0; x < 1111; x++ {
        s.CheckAdd(t, "not important", 0, i)
    }

    j := s.CheckAdd(t, "world!", 43, i)
    for x := 0; x < 1111; x++ {
        s.CheckAdd(t, "not important 2", 1, i, j)
    }
    if s.Name(i) != "Hello" || s.Value(i) != 42 {
        t.Errorf("unexpected value (1) {%q %d}", s.Name(i), s.Value(i))
    }
    if s.Name(j) != "world!" || s.Value(j) != 43 {
        t.Errorf("unexpected value (2) {%q %d}", s.Name(j), s.Value(j))
    }

    CheckValid(t, s)
}

func TestRandom(t *testing.T) {
    s := NewSystem(1024)

    rand.Seed(0x7254856391272548)
    names := []string{
        "Antepenult", "Disproves", "Patchboards", "Tanna",
        "Ichthyological", "Proles", "Renounce", "Encomiastically",
        "Ravaging", "Duane", "Overproud", "Frequences",
        "Vapour", "Earbob", "Frostiness", "Dracunculuses",
        "Wheedled", "Vomitings", "Sands", "Chain-smoking",
        "One-way", "T-junction", "Ruing", "Cautionary",
        "Digressive", "Sutler", "Berated", "Duff",
        "Beggarman", "Denationalises", "Asceticism", "Medusan",
        "Mailcoaches", "Extensiveness", "Engraft", "Inculcator",
        "Bestiaries", "Unjointed", "Weapons", "Misdating",
        "Revanches", "Suburbans", "Wattlings", "Square-rigger",
        "Fledgling", "Thoughtless", "Dzos", "Agonizedly",
        "Releasing", "Citification", "Millwrights", "Bales",
        "Bitty", "Arafat", "Parquetry", "Intestinal",
        "Roseola", "Resoundingly", "Mogul", "Meekness",
        "Cricketers", "Unpoised", "Arbitrager", "Corks",
        "Suffusions", "Aruba", "Iatrochemist", "Thrivingly",
        "Ganoin", "Canalisation", "Oncogenes", "Aliveness",
        "Amritsar", "Ridders", "Unpardoning", "Champlain",
        "Childishly", "Skins", "Sun", "Dissuades",
        "Lixiviations", "Miscreative", "Aniconisms", "Decay",
        "Granite", "Cou-cou", "Girlie", "Occupied",
        "Vibraphonists", "Ascetics", "Multidisciplinary", "Waysides",
        "Gleets", "Patrilineal", "Suitors", "Trichophytons",
        "Critically", "Auriculated", "Gall", "Unisons",
    }
    indices := make([]index.I, 0, len(names))
    mapping := make(map[index.I]string, len(names))

    for x := 0; x < 2048; x++ {
        switch rand.Intn(2) {
        case 0:
            if len(names) == 0 {
                x--
                continue
            }
            j := rand.Intn(len(names))
            n := names[j]
            names[j] = names[len(names)-1]
            names = names[:len(names)-1]

            i := s.CheckAdd(t, n, 111, indices...)
            indices = append(indices, i)

            if _, ok := mapping[i]; ok {
                t.Error("index", i, "already in use")
            }
            mapping[i] = n
        case 1:
            if len(indices) == 0 {
                x--
                continue
            }
            j := rand.Intn(len(indices))
            i := indices[j]
            indices[j] = indices[len(indices)-1]
            indices = indices[:len(indices)-1]
            names = append(names, s.Name(i))

            s.Remove(i)
            delete(mapping, i)
        default:
            panic("unreachable")
        }
    }

    for i, n := range mapping {
        if s.Name(i) != n {
            t.Errorf("broken index; %q != %q", s.Name(i), n)
        }
    }

    CheckValid(t, s)
}

func BenchmarkAdd4096(b *testing.B) {
    s := NewSystem(4096)
    for i := 0; i < b.N; i++ {
        s.Clear()

        for i := 0; i < 4096; i++ {
            s.Add("Benchmarking", i)
        }
    }
}

func BenchmarkAdd4096GrowFrom0(b *testing.B) {
    s := NewSystem(0)

    for i := 0; i < b.N; i++ {
        s.Clear()

        for i := 0; i < 4096; i++ {
            s.Add("Benchmarking", i)
        }
    }
}

func BenchmarkAdd4096Cap1024(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        s := NewSystem(1024)
        b.StartTimer()

        for i := 0; i < 4096; i++ {
            s.Add("Benchmarking", i)
        }
    }
}

func BenchmarkAdd4096Cap0(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        s := NewSystem(0)
        b.StartTimer()

        for i := 0; i < 4096; i++ {
            s.Add("Benchmarking", i)
        }
    }
}

func BenchmarkRemove4096(b *testing.B) {
    s := NewSystem(4096)
    is := [4096]index.I{}

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        s.Clear()
        for i := range is {
            is[i] = s.Add("Benchmarking", i)
        }
        b.StartTimer()

        for _, i := range is {
            s.Remove(i)
        }
    }
}

func BenchmarkRandomRemove4096(b *testing.B) {
    s := NewSystem(4096)
    is := [4096]index.I{}

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        s.Clear()
        for i := range is {
            is[i] = s.Add("Benchmarking", i)
        }
        for i := range is {
            j := rand.Intn(len(is))
            is[i], is[j] = is[j], is[i]
        }
        b.StartTimer()

        for _, i := range is {
            s.Remove(i)
        }
    }
}

func BenchmarkIterateSum4096(b *testing.B) {
    s := NewSystem(4096)

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        s.Clear()
        for i := 0; i < 4096; i++ {
            s.Add("Benchmarking", i)
        }
        b.StartTimer()

        sum := 0
        for _, v := range s.Values {
            sum += v
        }
    }
}
