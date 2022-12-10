package anamoni

// Subnet はサブネットの情報を保持する。
type Subnet struct {
	Addr    string
	Servers []string
}

func NewSubnet(addr string) *Subnet {
	return &Subnet{Addr: addr, Servers: []string{}}
}

func (s Subnet) Has(server string) bool {
	for _, srv := range s.Servers {
		if srv == server {
			return true
		}
	}
	return false
}

type Subnets []Subnet

func (ss Subnets) Brokens(brokenMap TroublesMap) Troubles {
	ts := Troubles{}
	for _, sn := range ss {
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
