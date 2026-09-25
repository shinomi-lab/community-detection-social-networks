import json
import os
from pathlib import Path
import random

import networkx as nx


def generate_lfr_digraph(seed: int = 6) -> nx.DiGraph:
    """LFRベンチマークを生成し、各エッジに方向を付与して有向グラフにする関数"""
    n = 5000            # ノード数
    tau1 = 3.0          # 次数分布の指数
    tau2 = 1.016        # コミュニティサイズ分布の指数
    mu = 0.8            # ミキシングパラメータ
    average_degree = 10 # 平均次数
    min_community = 100 # 最小コミュニティサイズ
    max_community = 1500

    # 1. 無向LFRグラフを生成 (directed 引数は削除)
    G_undirected = nx.LFR_benchmark_graph(
        n,
        tau1,
        tau2,
        mu,
        average_degree=average_degree,
        min_community=min_community,
        max_community=max_community,
        seed=seed,
    )

    # 2. 有向グラフを作成し、ノードとその属性（community情報）を引き継ぐ
    G = nx.DiGraph()
    G.add_nodes_from(G_undirected.nodes(data=True))

    # 3. エッジに向きをつけて追加（50%の確率で u -> v または v -> u）
    rng = random.Random(seed)
    for u, v in G_undirected.edges():
        if u == v:
            continue
        if rng.random() < 0.5:
            G.add_edge(u, v)
        else:
            G.add_edge(v, u)

    return G


def write_graph(G: nx.DiGraph, p: Path):
    """有向グラフのエッジリストとコミュニティ情報を保存する関数"""
    p.mkdir(parents=True, exist_ok=True)

    # 1. 有向エッジリストの保存（起点 -> 終点の順で保存されます）
    nx.write_edgelist(G, p.joinpath("LFR_edgelist.txt"), data=False)
    nx.write_edgelist(G, p.joinpath("LFR_edgelist.txt.gz"), data=False)

    # 2. コミュニティ所属情報の抽出とJSON保存
    communities = {frozenset(c) for c in nx.get_node_attributes(G, "community").values()}
    
    node_to_community = {}
    for comm_id, comm_nodes in enumerate(communities):
        for node in comm_nodes:
            node_to_community[int(node)] = comm_id

    with open(p.joinpath("LFR_communities_adj.json"), "w", encoding="utf-8") as f:
        json.dump(node_to_community, f, indent=2)

    print(f"保存完了: コミュニティ数 = {len(communities)}, ノード数 = {G.number_of_nodes()}, エッジ数 = {G.number_of_edges()}")


def read_graph(p: Path) -> tuple[nx.DiGraph, dict[int, int]]:
    """保存した有向エッジリストとコミュニティ情報を読み込み、検証する関数"""
    # create_using=nx.DiGraph を指定して有向グラフとして読み込み
    G = nx.read_edgelist(p.joinpath("LFR_edgelist.txt"), nodetype=int, create_using=nx.DiGraph)

    # コミュニティJSONの読み込み
    with open(p.joinpath("LFR_communities_adj.json"), "r", encoding="utf-8") as f:
        data = json.load(f)
        node_to_community = {int(k): v for k, v in data.items()}

    # 検証: 有向グラフであること & ノードの所属確認
    assert isinstance(G, nx.DiGraph), "有向グラフとして読み込まれていません"
    assert len(G.nodes()) == len(node_to_community)
    for node in G.nodes():
        assert node in node_to_community

    print("読み込み・有向グラフ検証完了！")
    return G, node_to_community


if __name__ == "__main__":
    dir_path = Path(os.getcwd()).parent.joinpath("/Users/kaoriogawa/研究/community-detection-social-networks-dev/difftools/network")
    print(f"保存先ディレクトリ: {dir_path}")

    # 有向グラフ生成
    digraph = generate_lfr_digraph()

    # 保存
    write_graph(digraph, dir_path)

    # 読み込みと検証
    loaded_digraph, loaded_communities = read_graph(dir_path)