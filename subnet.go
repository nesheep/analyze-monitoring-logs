package anamoni

// Subnet はサブネットの情報を保持する構造体。
type Subnet struct {
	Addr    string
	Servers []string
}

func NewSubnet(addr string) *Subnet {
	return &Subnet{Addr: addr, Servers: []string{}}
}

// Has は server として受け取った IP アドレスが Subnet に含まれるか判定する。
func (s Subnet) Has(server string) bool {
	for _, srv := range s.Servers {
		if srv == server {
			return true
		}
	}
	return false
}

// ExistsAll は srvs として受け取った IP アドレスのスライスの中に
// Subnet の Servers が全て含まれているかを判定する。
func (s Subnet) ExistsAll(srvs []string) bool {
	for _, a := range s.Servers {
		exists := false
		for _, b := range srvs {
			if b == a {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}

// Subnets は複数のサブネットの情報を保持するスライス。
type Subnets []Subnet

// Brokens は brokenMap としてサーバー故障期間を受け取ってサブネットの故障期間を Troubles として返す。
func (ss Subnets) Brokens(brokenMap TroublesMap) Troubles {
	ts := Troubles{}
	for _, sn := range ss {
		if !sn.ExistsAll(brokenMap.Servers()) {
			continue
		}
		tmp := brokenMap[sn.Servers[0]]
		for k, troubles := range brokenMap {
			if !sn.Has(k) {
				continue
			}
			intersections := troubles.Intersections(sn.Addr, tmp)
			tmp = intersections
		}
		ts = append(ts, tmp...)
	}
	return ts
}
