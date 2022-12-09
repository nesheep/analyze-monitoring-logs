package anamoni

import (
	"sort"
	"time"
)

// Trouble 障害情報を保持する構造体。
type Trouble struct {
	Addr  string
	Start *time.Time
	End   *time.Time
}

func NewTrouble(addr string, start time.Time) Trouble {
	return Trouble{Addr: addr, Start: &start}
}

func (t *Trouble) SetEnd(end time.Time) {
	t.End = &end
}

type Troubles []Trouble

// TroublesMap はIPアドレスごとの Troubles を記録する。
type TroublesMap map[string]Troubles

func (tm TroublesMap) ToSlice() Troubles {
	sl := Troubles{}
	for _, ts := range tm {
		sl = append(sl, ts...)
	}

	sort.Slice(sl, func(i, j int) bool {
		if sl[i].Start == nil {
			if sl[j].Start == nil {
				return sl[i].End.Before(*sl[j].End)
			}
			return true
		}
		if sl[j].Start == nil {
			return false
		}
		return sl[i].Start.Before(*sl[j].Start)
	})

	return sl
}
