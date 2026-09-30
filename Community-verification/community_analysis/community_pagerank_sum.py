from collections import defaultdict
import json
from pathlib import Path

import networkx as nx
import numpy as np
import pandas as pd
import scipy as sp


def load_data(data_dir: Path):
    """データファイル (.gzエッジリスト, JSON形式のコミュニティテキスト) を読み込む関数"""
    edgelist_path = data_dir / "LFR_edgelist.txt.gz"
    communities_path = data_dir / "LFR_communities.txt"

    # 1. 有向グラフの読み込み (.gz圧縮ファイルを自動解凍)
    G = nx.read_edgelist(
        edgelist_path,
        nodetype=int,
        create_using=nx.DiGraph
    )

    # 2. JSON形式の LFR_communities.txt を読み込み
    with open(communities_path, "r", encoding="utf-8") as f:
        node_to_community_raw = json.load(f)
        # キーを文字列から int 型に変換
        node_to_community = {int(k): v for k, v in node_to_community_raw.items()}

    # コミュニティIDごとにノードをリスト化
    comm_to_members = defaultdict(list)
    for node_id, comm_id in node_to_community.items():
        comm_to_members[comm_id].append(node_id)

    communities_list = list(comm_to_members.values())
    return G, communities_list


def main():
    # 実行ファイルの位置を基準に相対パスで data ディレクトリを指定
    base_dir = Path(__file__).resolve().parent.parent.parent
    data_dir = base_dir / "data"

    print(f"📁 データ参照元: {data_dir}")

    # データ読み込み
    G, communities_list = load_data(data_dir)

    # 1. グラフ全体で一度だけPageRankを計算する
    full_pr = nx.pagerank(G)

    # 2. 各コミュニティのPageRank集計
    results = []
    for i, comm_members in enumerate(communities_list):
        comm_scores = [full_pr[node] for node in comm_members if node in full_pr]

        results.append({
            "Community ID": i,
            "Size": len(comm_members),
            "PPR Sum": np.sum(comm_scores),
            "PPR Mean": np.mean(comm_scores)
        })

    results_df = pd.DataFrame(results)

    # 検算: 全コミュニティのPPR Sumの合計を表示 (1.0になるはず)
    print(f"Total PPR Sum (all communities): {results_df['PPR Sum'].sum():.4f}")

    # PPR Sumが高い順にソートして表示
    sorted_by_sum_df = results_df.sort_values(by="PPR Sum", ascending=False)
    total_ppr_sum = results_df["PPR Sum"].sum()

    # 表示用テキストの整形
    output_text = (
        f"Total PPR Sum (all communities): {total_ppr_sum:.4f}\n\n"
        "📊 【コミュニティごとの PageRank 集計結果】\n"
        f"{sorted_by_sum_df.to_string(index=False)}\n"
    )

    # テキストファイルへの保存処理
    output_file_path = data_dir / "community_pagerank_sum_result.txt"
    with open(output_file_path, "w", encoding="utf-8") as f:
        f.write(output_text)


if __name__ == "__main__":
    main()