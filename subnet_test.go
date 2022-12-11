package anamoni

import "testing"

func TestSubnet_Has(t *testing.T) {
	type args struct {
		server string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{server: "10.20.30.1"},
			want: true,
		},
		{
			name: "false",
			args: args{server: "192.168.1.1"},
			want: false,
		},
	}

	sn := Subnet{
		Addr:    "10.20.0.0",
		Servers: []string{"10.20.30.1", "10.20.30.2", "10.20.30.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sn.Has(tt.args.server)
			if got != tt.want {
				t.Errorf("expect %v, but %v", tt.want, got)
			}
		})
	}
}

func TestSubnet_ExistsAll(t *testing.T) {
	type args struct {
		srvs []string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{srvs: []string{"10.20.30.1", "10.20.30.2", "10.20.30.3", "192.168.1.1"}},
			want: true,
		},
		{
			name: "false",
			args: args{srvs: []string{"10.20.30.1", "10.20.30.3", "192.168.1.1"}},
			want: false,
		},
	}

	sn := Subnet{
		Addr:    "10.20.0.0",
		Servers: []string{"10.20.30.1", "10.20.30.2", "10.20.30.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sn.ExistsAll(tt.args.srvs)
			if got != tt.want {
				t.Errorf("expect %v, but %v", tt.want, got)
			}
		})
	}
}
