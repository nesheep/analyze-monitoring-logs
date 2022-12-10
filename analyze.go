package anamoni

// Analyze は logs の分析を行ってサーバー故障期間、サーバー過負荷期間、サブネット故障期間を Troubles として返す。
func Analyze(logs Logs, breakJudgment, overloadJudgment, overloadTime int) (Troubles, Troubles, Troubles) {
	brokenMap := TroublesMap{}
	overlaodMap := TroublesMap{}
	srvMap := ServerStatusMap{}

	servers := logs.Servers()
	for _, addr := range servers {
		srvMap[addr] = &ServerStatus{}
	}

	// サーバーに関する検査
	for _, l := range logs {
		addr := l.Address()
		srvStat := srvMap[addr]
		srvMap[addr].Logs = append(srvMap[addr].Logs, l)

		// 故障検査
		isBroken := srvStat.CheckBroken(breakJudgment)
		if isBroken {
			if len(srvStat.Logs) == breakJudgment {
				// 1 件目のログから故障していた場合: Start は nil とする
				brokenMap[addr] = append(brokenMap[addr], NewTrouble(addr, NewDuration(nil, nil)))
			} else if !srvStat.IsBroken {
				// 非故障中から故障中になる場合
				d := NewDuration(&srvStat.Logs[len(srvStat.Logs)-breakJudgment].Time, nil)
				brokenMap[addr] = append(brokenMap[addr], NewTrouble(addr, d))
			}
		} else if srvStat.IsBroken {
			// 故障中から非故障中になる場合
			brokenMap[addr][len(brokenMap[addr])-1].SetEnd(l.Time)
		}
		srvMap[addr].IsBroken = isBroken

		// 過負荷検査
		isOverloaded := srvStat.CheckOverloaded(overloadJudgment, overloadTime)
		if isOverloaded {
			if len(srvStat.Logs) == overloadJudgment {
				// 1 件目のログから過負荷だった場合: Start は nil とする
				overlaodMap[addr] = append(overlaodMap[addr], NewTrouble(addr, NewDuration(nil, nil)))
			} else if !srvStat.IsOverloaded {
				// 非過負荷状態から過負荷状態になる場合
				d := NewDuration(&srvStat.Logs[len(srvStat.Logs)-overloadJudgment].Time, nil)
				overlaodMap[addr] = append(overlaodMap[addr], NewTrouble(addr, d))
			}
		} else if srvStat.IsOverloaded {
			// 過負荷状態から非過負荷状態になる場合
			overlaodMap[addr][len(overlaodMap[addr])-1].SetEnd(l.Time)
		}
		srvMap[addr].IsOverloaded = isOverloaded
	}

	// サブネットに関する検査
	subnets := logs.Subnets()
	snBrokens := subnets.Brokens(brokenMap)

	return brokenMap.ToSlice(), overlaodMap.ToSlice(), snBrokens
}
