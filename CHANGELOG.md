### 2026/9/16
#### ディレクトリ
- Goプロジェクトのルートを[./difftools](./difftools)に変更
- [./difftools/sample](./difftools/sample)のうち、プロジェクトと無関係のものを[./sample](./sample)に隔離

#### difftools
- 同一ロジックの共通化
- `main`関数ファイルのコマンド化（`difftools/cmd`配下）
- 乱数生成器のオブジェクトの明示・共通化（ver 1.18標準）
- 以下を満たすシード配列用の定数を定義
    - `SeedInfoF == InfoType_F + 1`
    - `SeedInfoT == InfoType_T + 1`
- グラフ情報を`Network`構造体に一元化
  - `Adj [][]int`: 隣接行列
  - `N`: ノード数 (`== len(Adj)`)
  - `FollowerNums`: 各ノードのフォロワー数 (`optimzation.FolowerSize`関数の廃止)
