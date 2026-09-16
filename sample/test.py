import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.linear_model import LinearRegression

# 1. pandasライブラリのインポート
# 2. CSVファイルの読み込み
# 例: 'data.csv' という名前のCSVファイルを読み込む
df = pd.read_csv('実行結果\calSameImpression_twitter8_1_kirinuki.csv',header=None)

# 3. データの表示
print(df)


# 2. 品種ごとにグループ化して最大値を計算
max_infl_by_node_num = df.groupby(1)[2].max().reset_index()

# 結果の表示
print(max_infl_by_node_num)

df = max_infl_by_node_num



# 1. データの準備（例としてランダムなデータを使用）
# 実際にはあなたのデータを使用してください
x = df.iloc[:,0].values#[0:50]
y = df.iloc[:,1].values#[0:50] # 例: 対数曲線に近いデータ

# 2. 対数変換を行う
log_x = np.log(x)
log_x_2d = log_x.reshape(-1, 1)

# 3. 線形回帰モデルの作成
model = LinearRegression()

# 4. モデルのトレーニング
model.fit(log_x_2d, y)

# 5. 回帰直線の予測
y_pred = model.predict(log_x_2d)

# 6. 結果の表示
plt.scatter(x, y, color='blue', label='実データ')
plt.plot(x, y_pred, color='red', label='回帰曲線（対数変換）')
plt.xlabel('X')
plt.ylabel('Y')
plt.legend()
plt.show()

# 7. 回帰曲線の数式の表示
slope = model.coef_[0]
intercept = model.intercept_
print(f"回帰曲線の数式: y = {intercept} + {slope} * log(x)")
print(model.score(log_x_2d,y))
