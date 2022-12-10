package anamoni

import (
	"sort"
	"time"
)

type Duration struct {
	Start *time.Time
	End   *time.Time
}

func NewDuration(start, end *time.Time) Duration {
	return Duration{Start: start, End: end}
}

func (d *Duration) Intersection(o Duration) *Duration {
	start := *d.Start
	end := *d.End
	if start.Before(*o.Start) {
		start = *o.Start
	}
	if end.After(*o.End) {
		end = *o.End
	}
	if start.After(end) {
		return nil
	}
	return &Duration{
		Start: &start,
		End:   &end,
	}
}

// Trouble 障害情報を保持する構造体。
type Trouble struct {
	Addr string
	Duration
}

func NewTrouble(addr string, d Duration) Trouble {
	return Trouble{Addr: addr, Duration: d}
}

func (t *Trouble) SetEnd(end time.Time) {
	t.End = &end
}

type Troubles []Trouble

func (ts Troubles) Intersections(addr string, o Troubles) Troubles {
	is := Troubles{}
	for _, a := range ts {
		for _, b := range o {
			intersection := a.Duration.Intersection(b.Duration)
			if intersection != nil {
				is = append(is, NewTrouble(addr, *intersection))
			}
		}
	}
	return is
}

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
