#!/bin/bash

config_dir=./settings
k=8  # 最大並列数

c=0  # 現在実行中のプロセス数

for file in "${config_dir}"/*.json; do
  while [[ $c -ge $k ]]; do
    wait -n  # いずれかのバックグラウンドプロセスが終了するまで待機
    c=$((c - 1)) # 終了したプロセス数を減算
  done
  if [ -f "$file" ]; then
    # goのプログラムをバックグラウンドで実行
    go run json_file_loading.go $file &
    c=$((c + 1)) # 実行中のプロセス数を加算
  fi
done

wait