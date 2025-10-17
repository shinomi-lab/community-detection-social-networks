# 必要なライブラリをインポートします
import community.community_louvain as community_louvain
import networkx as nx
import requests
import tarfile
import os
import pandas as pd # データ整理と表示のためにpandasをインポートします
import numpy as np   # NumPyを隣接行列の型変換で使用

# --- 1. データのダウンロードと展開 (前回と同様) ---
url = "https://snap.stanford.edu/data/twitter.tar.gz"
filename = "ego-twitter/twitter.tar.gz"
data_file = "ego-twitter/twitter_combined.txt.gz"
output_dir = "twitter_communities2"

if not os.path.exists(data_file):
    print(f"'{filename}'をダウンロードしています...")
    r = requests.get(url, stream=True)
    if r.status_code == 200:
        with open(filename, "wb") as f:
            for chunk in r.iter_content(chunk_size=1024):
                f.write(chunk)
        print("ダウンロードが完了しました。")

        print(f"'{filename}'を展開しています...")
        with tarfile.open(filename, "r:gz") as tar:
            tar.extractall()
        print("展開が完了しました。")
    else:
        print("データのダウンロードに失敗しました。")
        exit()

# --- 2. グラフの構築 (前回と同様) ---
print("グラフを構築しています...")
G = nx.read_edgelist(data_file, create_using=nx.Graph(), nodetype=int)
print("グラフの構築が完了しました。")
print(f"ノード数: {G.number_of_nodes()}")
print(f"エッジ数: {G.number_of_edges()}")

# --- 3. Louvain法によるコミュニティ検出 (前回と同様) ---
print("Louvain法によるコミュニティ検出を実行しています...")
partition = community_louvain.best_partition(G)
print("コミュニティ検出が完了しました。")

# 検出されたコミュニティの数を計算します (partition辞書の値のユニークな数を数えます)
num_communities = len(set(partition.values()))
print(f"\n検出されたコミュニティの数: {num_communities}")

# コミュニティ分割の良さを示すモジュラリティを計算します
modularity = community_louvain.modularity(partition, G)
print(f"モジュラリティ: {modularity:.4f}")

'''
print("コミュニティごとにエッジリストをファイルに保存しています...")
# 出力用のディレクトリが存在しない場合は作成します
if not os.path.exists(output_dir):
    os.makedirs(output_dir)

# コミュニティIDをキー、そのコミュニティに属するノードのリストを値とする辞書を作成します
communities = {}
for node, community_id in partition.items():
    if community_id not in communities:
        communities[community_id] = []
    communities[community_id].append(node)

# 各コミュニティについてループ処理を行います
for community_id, nodes in communities.items():
    # 現在のコミュニティに属するノードからサブグラフを作成します
    subgraph = G.subgraph(nodes)
    # 保存するファイル名を定義します
    output_filename = os.path.join(output_dir, f"community_{community_id}.txt")
    # ファイルを書き込みモードで開きます
    with open(output_filename, "w") as f:
        # サブグラフのすべてのエッジについてループします
        for u, v in subgraph.edges():
            # エッジを構成する2つのノードをファイルに書き込みます
            f.write(f"{u} {v}\n")

print(f"'{output_dir}'ディレクトリに{num_communities}個のコミュニティファイルを保存しました。")
'''

# --- 4. コミュニティごとのグループ化 (前回と同様) ---
if not os.path.exists(output_dir):
    os.makedirs(output_dir)
communities = {}
for node, community_id in partition.items():
    if community_id not in communities:
        communities[community_id] = []
    communities[community_id].append(node)

# 各コミュニティについてループ処理を行います
for community_id, nodes in communities.items():
    # 現在のコミュニティに属するノードからサブグラフを作成します
    subgraph = G.subgraph(nodes)
    # 保存するファイル名を定義します
    output_filename = os.path.join(output_dir, f"community_{community_id}.txt")
    # ファイルを書き込みモードで開きます
    with open(output_filename, "w") as f:
        # サブグラフのすべてのエッジについてループします
        for u, v in subgraph.edges():
            # エッジを構成する2つのノードをファイルに書き込みます
            f.write(f"{u} {v}\n")

# --- 5. 【追加】コミュニティの密度の計算と分析 ---
print("\n各コミュニティの密度を計算し、分析しています...")
community_stats = [] # 各コミュニティの統計情報を格納するリスト

# 各コミュニティについてループ処理を行います
for community_id, nodes in communities.items():
    # コミュニティのノード数が2以上の場合のみ密度を計算します (ノード1つでは密度が定義できないため)
    if len(nodes) > 1:
        # 現在のコミュニティに属するノードからサブグラフを作成します
        subgraph = G.subgraph(nodes)
        # サブグラフのノード数を取得します
        num_nodes = subgraph.number_of_nodes()
        # サブグラフのエッジ数を取得します
        num_edges = subgraph.number_of_edges()
        # サブグラフの密度を計算します
        density = nx.density(subgraph)

        # 計算結果をリストに追加します
        community_stats.append({
            "community_id": community_id,
            "num_nodes": num_nodes,
            "num_edges": num_edges,
            "density": density
        })

# 統計情報をpandasのDataFrameに変換して見やすくします
df_stats = pd.DataFrame(community_stats)

# --- 6. 結果の表示 ---

# 密度が1.0の完全グラフ（クリック）のコミュニティを除外して、最も密なコミュニティを探します
# （完全グラフは自明な密なコミュニティのため、それ以外で興味深いものを探します）
most_dense_community = df_stats[df_stats['density'] < 1.0].sort_values(by='density', ascending=False).iloc[0]

# エッジが1つ以上あるコミュニティの中で、最も疎なコミュニティを探します
most_sparse_community = df_stats[df_stats['num_edges'] > 0].sort_values(by='density', ascending=True).iloc[0]

# 密度の中央値を計算
median_density = df_stats['density'].median()

# 中央値との差の絶対値を計算し、最も差が小さいコミュニティを特定
df_stats['median_diff'] = (df_stats['density'] - median_density).abs()
median_community = df_stats.sort_values(by='median_diff', ascending=True).iloc[0]

print("\n--- 分析結果 ---")
print(f"\n密度の中央値: {median_density:.6f}")

print("\n✅ 最も密なコミュニティ（完全グラフを除く）:")
print(most_dense_community.to_string())

print("\n\n✅ 中央値に近いコミュニティ:")
print(median_community[['community_id', 'num_nodes', 'num_edges', 'density']].to_string())

print("\n\n✅ 最も疎なコミュニティ（エッジが1つ以上あるもの）:")
print(most_sparse_community.to_string())
print("\n---------------\n")

print(f"\nコミュニティごとの隣接行列を'{output_dir}'に保存しています...")

# 出力用のディレクトリが存在しない場合は作成します
if not os.path.exists(output_dir):
    os.makedirs(output_dir)

# 各コミュニティについてループ処理を行います
for community_id, nodes in communities.items():
    # サブグラフを作成します
    subgraph = G.subgraph(nodes)
    
    # サブグラフの隣接行列をNumPy配列として取得します
    # ノードの順序を固定するためにnodelistを指定します
    adj_matrix = nx.to_numpy_array(subgraph, nodelist=sorted(subgraph.nodes()))
    
    # ご提示のコードに合わせてデータ型を整数に変換します
    adj_matrix = adj_matrix.astype(np.int32)
    
    # NumPy配列をpandas DataFrameに変換します
    # indexとcolumnsにノードIDリストを指定することで、JSONがより分かりやすくなります
    node_list = sorted(subgraph.nodes())
    df_adj = pd.DataFrame(adj_matrix, index=node_list, columns=node_list)
    
    # 保存するファイル名を定義します
    output_filename = os.path.join(output_dir, f"community_{community_id}_adj.txt")
    
    # DataFrameをJSONファイルとして保存します
    df_adj.to_json(output_filename, orient="index")

print("JSONファイルの保存が完了しました。")