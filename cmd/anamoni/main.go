package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	anamoni "github.com/nesheep/analyze-monitoring-logs"
)

var n int
var m int
var t int

func init() {
	flag.IntVar(&n, "n", 1, "n 回以上連続してタイムアウトした場合にサーバーの故障とみなす。")
	flag.IntVar(&m, "m", 1, "直近 m 回の平均応答時間が t ミリ秒を超えた場合はサーバが過負荷状態になっているとみなす。")
	flag.IntVar(&t, "t", 100, "直近 m 回の平均応答時間が t ミリ秒を超えた場合はサーバが過負荷状態になっているとみなす。")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("引数が足りません")
	}

	if n < 1 {
		log.Fatal("n は 1 以上の整数を指定してください")
	}

	if m < 1 {
		log.Fatal("m は 1 以上の整数を指定してください")
	}

	if t < 1 {
		log.Fatal("t は 1 以上の整数を指定してください")
	}

	filename := args[0]
	logs, err := anamoni.ReadLogs(filename)
	if err != nil {
		log.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}

	brokens, overloads, snBrokens := anamoni.Analyze(logs, n, m, t)

	if len(brokens) > 0 {
		fmt.Println("サーバー故障期間")
	}
	for i, t := range brokens {
		fmt.Printf("%d\t%s\t%s\t%s\n", i+1, t.Addr, formatTime(t.Start), formatTime(t.End))
	}

	if len(overloads) > 0 {
		fmt.Println("サーバー過負荷期間")
	}
	for i, t := range overloads {
		fmt.Printf("%d\t%s\t%s\t%s\n", i+1, t.Addr, formatTime(t.Start), formatTime(t.End))
	}

	if len(snBrokens) > 0 {
		fmt.Println("サブネット故障期間")
	}
	for i, t := range snBrokens {
		fmt.Printf("%d\t%s\t%s\t%s\n", i+1, t.Addr, formatTime(t.Start), formatTime(t.End))
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return "                   "
	}
	return t.Format("2006/01/02 15:04:05")
}
