import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.linear_model import LinearRegression
from sklearn.preprocessing import PolynomialFeatures
from sklearn.metrics import r2_score

df = pd.read_csv('実行結果\calSameImpression_twitter8_1.csv',header=None)

max_infl_by_node_num = df.groupby(1)[2].max().reset_index()

df = max_infl_by_node_num



# 1. データの準備（例としてランダムなデータを使用）
# 実際にはあなたのデータを使用してください
x = df.iloc[:,0].values.reshape(-1, 1)
y = df.iloc[:,1].values#[0:60]

# サンプルデータの作成
X = x

# 特徴量の変換
poly = PolynomialFeatures(degree=2)
X_poly = poly.fit_transform(X)

# モデルの作成と学習
model = LinearRegression()
model.fit(X_poly, y)

# 新しいデータに対する予測
X_new = np.array([6, 7, 8]).reshape(-1, 1)
X_new_poly = poly.transform(X_new)
y_pred = model.predict(X_new_poly)

print("Predicted values:", y_pred)

# R2スコアの計算
y_train_pred = model.predict(X_poly)
r2 = r2_score(y, y_train_pred)
print("R2 score:", r2)

# プロット
plt.scatter(X, y, color='blue')
plt.plot(X, y_train_pred, color='red')
plt.show()
