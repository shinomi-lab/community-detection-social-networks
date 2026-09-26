# Community-verification

## Python環境の構築（初回のみ）
```
cd Community-verification
uv sync
```

## Pythonスクリプトの実行
例：foo/bar.py
```
cd Community-verification
uv run python foo/bar.py
```

# difftools

## Go言語での実行用プログラムの作成
1. [./difftools/cmd](./difftools/cmd) 直下に適当なコマンド名（空白文字は避ける）のフォルダを追加する。
2. 追加したフォルダ直下に `main.go` ファイルを作成する。
3. `main.go` に以下の内容を記述する。
   ```go
   package main

   func main() {
   }
   ``` 

## プログラムの実行
```shell
cd difftools
go run ./cmd/[コマンド名]
```
