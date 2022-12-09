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

// Analyze は logs 分析を行って TroublesMap を返す。
func Analyze(logs Logs) TroublesMap {
	tm := TroublesMap{}

	for _, l := range logs {
		addr := l.Address()
		ts := tm[addr]

		// 各サーバーの 1 件目のログ処理
		if ts == nil {
			tm[addr] = Troubles{}
			if l.Timeouted {
				tm[addr] = append(tm[addr], Trouble{Addr: addr})
			}
			continue
		}

		// 各サーバーの 2 件目以降のログ処理
		if l.Timeouted {
			if len(ts) == 0 || ts[len(ts)-1].End != nil {
				tm[addr] = append(tm[addr], NewTrouble(addr, l.Time))
			}
			continue
		}

		if len(ts) > 0 && ts[len(ts)-1].End == nil {
			ts[len(ts)-1].SetEnd(l.Time)
		}
	}

	return tm
}
