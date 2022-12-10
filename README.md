# Analyze monitoring logs

監視ログファイル ([サンプル](./testdata/0.csv)) を分析して障害情報を出力する CLI ツール。

## Install

```bash
go install github.com/nesheep/analyze-monitoring-logs/cmd/anamoni@latest
```

## Usage

```bash
anamoni <OPTIONS> <FILENAME>
```

### OPTIONS

- n: n 回以上連続してタイムアウトした場合にサーバー故障とみなす (default 1)
- m: 直近 m 回の平均応答時間が t ミリ秒を超えた場合にサーバー過負荷状態とみなす (default 1)
- t: 直近 m 回の平均応答時間が t ミリ秒を超えた場合にサーバー過負荷状態とみなす (default 100)

### Example

```bash
$ anamoni -n 3 -m 2 -t 100 testdata/3.csv
サーバー故障期間
1       10.20.30.1      2022/12/01 00:04:01     2022/12/01 00:08:01
2       10.20.30.2      2022/12/01 00:05:02     2022/12/01 00:08:02
サーバー過負荷期間
1       10.20.30.1      2022/12/01 00:08:01     2022/12/01 00:11:01
2       10.20.30.2      2022/12/01 00:09:02     2022/12/01 00:14:02
サブネット故障期間
1       10.20.0.0       2022/12/01 00:05:02     2022/12/01 00:08:01
```
