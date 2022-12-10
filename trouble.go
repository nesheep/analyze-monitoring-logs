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
	start := laterStart(d.Start, o.Start)
	end := earlierEnd(d.End, o.End)
	if start != nil && end != nil && start.After(*end) {
		return nil
	}
	return &Duration{Start: start, End: end}
}

func laterStart(a, b *time.Time) *time.Time {
	if a == nil {
		if b == nil {
			return nil
		}
		return b
	}
	if b == nil {
		return a
	}
	if a.After(*b) {
		return a
	}
	return b
}

func earlierEnd(a, b *time.Time) *time.Time {
	if a == nil {
		if b == nil {
			return nil
		}
		return b
	}
	if b == nil {
		return a
	}
	if a.Before(*b) {
		return a
	}
	return b
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

// Sort は Troubles を古い順にソートする。
func (ts Troubles) Sort() {
	sort.Slice(ts, func(i, j int) bool {
		if ts[i].Start == nil {
			if ts[j].Start == nil {
				if ts[i].End == nil {
					return false
				}
				if ts[j].End == nil {
					return true
				}
				return ts[i].End.Before(*ts[j].End)
			}
			return true
		}
		if ts[j].Start == nil {
			return false
		}
		return ts[i].Start.Before(*ts[j].Start)
	})
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
	sl.Sort()
	return sl
}
