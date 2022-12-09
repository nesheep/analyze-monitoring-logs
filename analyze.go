package anamoni

// Analyze は logs 分析を行って TroublesMap を返す。
func Analyze(logs Logs, breakJudgment int) TroublesMap {
	tm := TroublesMap{}
	srvm := ServerStatusMap{}

	servers := logs.Servers()
	for _, addr := range servers {
		tm[addr] = Troubles{}
		srvm[addr] = &ServerStatus{}
	}

	for _, l := range logs {
		addr := l.Address()
		srvs := srvm[addr]
		srvm[addr].Logs = append(srvm[addr].Logs, l)

		if len(srvs.Logs) < breakJudgment {
			continue
		}

		isBroken := true
		for i := 0; i < breakJudgment; i++ {
			if !srvs.Logs[len(srvs.Logs)-1-i].Timeouted {
				isBroken = false
				break
			}
		}

		if isBroken {
			if len(srvs.Logs) == breakJudgment {
				// 1 件目のログから故障していた場合
				tm[addr] = append(tm[addr], Trouble{Addr: addr})
			} else if !srvs.IsBroken {
				// 非故障中から故障中になる場合
				tm[addr] = append(tm[addr], NewTrouble(addr, srvs.Logs[len(srvs.Logs)-breakJudgment].Time))
			}
		} else if srvs.IsBroken {
			// 故障中から非故障中になる場合
			tm[addr][len(tm[addr])-1].SetEnd(l.Time)
		}

		srvm[addr].IsBroken = isBroken
	}

	return tm
}
