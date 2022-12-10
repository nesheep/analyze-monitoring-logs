package anamoni

// Analyze は logs 分析を行って TroublesMap を返す。
func Analyze(logs Logs, breakJudgment, overloadJudgment, overloadTime int) (Troubles, Troubles, Troubles) {
	bm := TroublesMap{}
	om := TroublesMap{}
	srvm := ServerStatusMap{}

	servers := logs.Servers()
	for _, addr := range servers {
		srvm[addr] = &ServerStatus{}
	}

	// サーバーに関する検査
	for _, l := range logs {
		addr := l.Address()
		srvs := srvm[addr]
		srvm[addr].Logs = append(srvm[addr].Logs, l)

		// 故障検査
		isBroken := srvs.CheckBroken(breakJudgment)
		if isBroken {
			if len(srvs.Logs) == breakJudgment {
				// 1 件目のログから故障していた場合
				bm[addr] = append(bm[addr], Trouble{Addr: addr})
			} else if !srvs.IsBroken {
				// 非故障中から故障中になる場合
				bm[addr] = append(bm[addr], NewTrouble(addr, srvs.Logs[len(srvs.Logs)-breakJudgment].Time))
			}
		} else if srvs.IsBroken {
			// 故障中から非故障中になる場合
			bm[addr][len(bm[addr])-1].SetEnd(l.Time)
		}
		srvm[addr].IsBroken = isBroken

		// 過負荷検査
		isOverloaded := srvs.CheckOverloaded(overloadJudgment, overloadTime)
		if isOverloaded {
			if len(srvs.Logs) == overloadJudgment {
				// 1 件目のログから過負荷だった場合
				om[addr] = append(om[addr], Trouble{Addr: addr})
			} else if !srvs.IsOverloaded {
				// 非過負荷状態から過負荷状態になる場合
				om[addr] = append(om[addr], NewTrouble(addr, srvs.Logs[len(srvs.Logs)-overloadJudgment].Time))
			}
		} else if srvs.IsOverloaded {
			// 過負荷状態から非過負荷状態になる場合
			om[addr][len(om[addr])-1].SetEnd(l.Time)
		}
		srvm[addr].IsOverloaded = isOverloaded
	}

	// サブネットに関する検査
	subnets := logs.Subnets()
	snBrokens := subnets.Brokens(bm)

	return bm.ToSlice(), om.ToSlice(), snBrokens
}
