package anamoni

import (
	"sort"
	"time"
)

// Duration は期間を表す構造体。
type Duration struct {
	Start *time.Time
	End   *time.Time
}

func NewDuration(start, end *time.Time) Duration {
	return Duration{Start: start, End: end}
}

// Intersection は他の Duration を受け取って共通期間を返す。
// 共通期間が存在しないときは nil を返す。
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

// Trouble は障害情報を保持する構造体。
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

// Troubles は複数の障害情報を保持するスライス。
type Troubles []Trouble

// Intersections は他の Troubles を受け取って共通期間の Troubles を返す。
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

// TroublesMap は IP アドレスごとの Troubles を記録する。
type TroublesMap map[string]Troubles

func (tm TroublesMap) Servers() []string {
	sl := []string{}
	for k := range tm {
		sl = append(sl, k)
	}
	return sl
}

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
