package anamoni

// ServerStatus は処理済みのログとサーバーの状態を保持する構造体。
type ServerStatus struct {
	Logs     Logs
	IsBroken bool
}

// ServerStatusMap はサーバーごとの ServerStatus を記録する。
type ServerStatusMap map[string]*ServerStatus
