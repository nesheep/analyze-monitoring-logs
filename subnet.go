package anamoni

// Subnet はサブネットの情報を保持する。
type Subnet struct {
	Addr    string
	Servers []string
}

func NewSubnet(addr string) *Subnet {
	return &Subnet{Addr: addr, Servers: []string{}}
}

func (s Subnet) has(server string) bool {
	for _, srv := range s.Servers {
		if srv == server {
			return true
		}
	}
	return false
}

type Subnets []Subnet

func (ss Subnets) Brokens(bm TroublesMap) Troubles {
	ts := Troubles{}
	for _, sn := range ss {
		tmp := bm[sn.Servers[0]]
		for k, troubles := range bm {
			if !sn.has(k) {
				continue
			}
			intersections := troubles.intersections(sn.Addr, tmp)
			tmp = intersections
		}
		ts = append(ts, tmp...)
	}
	return ts
}
