package main

import (
	"anamoni"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("引数が足りません")
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}

	r := csv.NewReader(f)
	logs := anamoni.Logs{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("レコードの読み込みに失敗しました: %v", err)
		}

		l, err := anamoni.ParseLog(record)
		if err != nil {
			log.Fatal(err)
		}

		logs = append(logs, l)
	}

	logs.Sort()

	tm := anamoni.Analyze(logs)

	tmSlice := tm.ToSlice()
	for i, t := range tmSlice {
		fmt.Printf("%d\t%s\t%s\t%s\n", i+1, t.Addr, formatTime(t.Start), formatTime(t.End))
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return "                   "
	}
	return t.Format("2006/01/02 15:04:05")
}
