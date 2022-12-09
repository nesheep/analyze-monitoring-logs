package anamoni

import (
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseLog(t *testing.T) {
	type args struct {
		record []string
	}

	tests := []struct {
		name   string
		args   args
		want   Log
		hasErr bool
	}{
		{
			name: "valid",
			args: args{record: []string{"20201019133124", "10.20.30.1/16", "2"}},
			want: Log{
				Time:         time.Date(2020, time.October, 19, 13, 31, 24, 0, time.UTC),
				IP:           net.IPv4(10, 20, 30, 1),
				Mask:         net.CIDRMask(16, 32),
				ResponseTime: 2,
				Timeouted:    false,
			},
			hasErr: false,
		},
		{
			name: "timeouted",
			args: args{record: []string{"20201019133124", "10.20.30.1/16", "-"}},
			want: Log{
				Time:         time.Date(2020, time.October, 19, 13, 31, 24, 0, time.UTC),
				IP:           net.IPv4(10, 20, 30, 1),
				Mask:         net.CIDRMask(16, 32),
				ResponseTime: 0,
				Timeouted:    true,
			},
			hasErr: false,
		},
		{
			name:   "invalid_time",
			args:   args{record: []string{"202010191331240", "10.20.30.1/16", "2"}},
			want:   Log{},
			hasErr: true,
		},
		{
			name:   "invalid_address",
			args:   args{record: []string{"20201019133124", "10.20.30.1:16", "2"}},
			want:   Log{},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLog(tt.args.record)
			if (err != nil) != tt.hasErr {
				t.Errorf("ParseLog() error = %v, hasErr %v", err, tt.hasErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ParseLog() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
