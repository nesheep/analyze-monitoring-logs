package anamoni

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Log はログファイルの各行の情報を格納する構造体。
type Log struct {
	Time         time.Time
	IP           net.IP
	Mask         net.IPMask
	ResponseTime int
	Timeouted    bool
}

// ParseLog はログファイルの各行の情報を []string として受け取って Log を返す。
func ParseLog(record []string) (Log, error) {
	if len(record) != 3 {
		return Log{}, errors.New("ログ形式が無効です")
	}

	t, err := parseTime(record[0])
	if err != nil {
		return Log{}, fmt.Errorf("時刻形式が無効です: %w", err)
	}

	ip, mask, err := parseAddress(record[1])
	if err != nil {
		return Log{}, fmt.Errorf("アドレス形式が無効です: %w", err)
	}

	rt, timeouted, err := parseResponseTime(record[2])
	if err != nil {
		return Log{}, fmt.Errorf("応答時間形式が無効です: %w", err)
	}

	return Log{
		Time:         t,
		IP:           ip,
		Mask:         mask,
		ResponseTime: rt,
		Timeouted:    timeouted,
	}, nil
}

func parseTime(s string) (time.Time, error) {
	layout := "20060102150405"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func parseAddress(s string) (net.IP, net.IPMask, error) {
	splited := strings.Split(s, "/")
	if len(splited) != 2 {
		return nil, nil, errors.New("'/'が必要です")
	}

	ip := net.ParseIP(splited[0])
	if ip == nil {
		return nil, nil, errors.New("IP形式が無効です")
	}

	ones, err := strconv.Atoi(splited[1])
	if err != nil {
		return nil, nil, err
	}
	mask := net.CIDRMask(ones, 32)

	return ip, mask, nil
}

func parseResponseTime(s string) (int, bool, error) {
	if s == "-" {
		return 0, true, nil
	}

	rt, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, err
	}

	return rt, false, nil
}

func (l Log) Address() string {
	return l.IP.String()
}

func (l Log) Subnet() string {
	return l.IP.Mask(l.Mask).String()
}

// Logs はログファイルの全ての行の情報を格納する。
type Logs []Log

// RaadLogs は CSV ファイルの filename を受け取って Logs を返す。
// CSV や ログの形式が無効の場合は error を返す。
func ReadLogs(filename string) (Logs, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	logs := Logs{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		l, err := ParseLog(record)
		if err != nil {
			return nil, err
		}

		logs = append(logs, l)
	}

	logs.Sort()
	return logs, nil
}

// Sort は Logs を Time の古い順にソートする。
func (ls Logs) Sort() {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Time.Before(ls[j].Time)
	})
}

func (ls Logs) Servers() []string {
	sl := []string{}
	m := map[string]bool{}
	for _, l := range ls {
		addr := l.Address()
		if !m[addr] {
			m[addr] = true
			sl = append(sl, addr)
		}
	}
	return sl
}

func (ls Logs) Subnets() Subnets {
	snMap := map[string]*Subnet{}
	srvMap := map[string]bool{}
	for _, l := range ls {
		subnet := l.Subnet()
		addr := l.Address()
		if snMap[subnet] == nil {
			snMap[subnet] = NewSubnet(subnet)
		}
		if !srvMap[addr] {
			srvMap[addr] = true
			snMap[subnet].Servers = append(snMap[subnet].Servers, addr)
		}
	}
	sl := Subnets{}
	for _, v := range snMap {
		sl = append(sl, *v)
	}
	return sl
}
