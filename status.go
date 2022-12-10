package anamoni

// ServerStatus は処理済みのログとサーバーの状態を保持する構造体。
type ServerStatus struct {
	Logs         Logs
	IsBroken     bool
	IsOverloaded bool
}

// CheckBroken は現在故障中か否かを判定する。
func (ss ServerStatus) CheckBroken(judgment int) bool {
	if len(ss.Logs) < judgment {
		return false
	}

	isBroken := true
	for i := 0; i < judgment; i++ {
		if !ss.Logs[len(ss.Logs)-1-i].Timeouted {
			isBroken = false
			break
		}
	}
	return isBroken
}

// CheckOverloaded は現在過負荷状態か否かを判定する。
func (ss ServerStatus) CheckOverloaded(judgment, t int) bool {
	if len(ss.Logs) < judgment {
		return false
	}

	sum := 0
	for i := 0; i < judgment; i++ {
		sum += ss.Logs[len(ss.Logs)-1-i].ResponseTime
	}
	return sum/judgment >= t
}

// ServerStatusMap はサーバーごとの ServerStatus を記録する。
type ServerStatusMap map[string]*ServerStatus
