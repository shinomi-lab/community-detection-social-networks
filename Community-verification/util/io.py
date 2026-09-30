import gzip
import json
from collections import defaultdict
from pathlib import Path

import networkx as nx


def write_communitydict(G: nx.DiGraph, p: Path):
    communities = {frozenset(c) for c in nx.get_node_attributes(G, "community").values()}

    node_to_community = {}
    for comm_id, comm_nodes in enumerate(communities):
        for node in comm_nodes:
            node_to_community[int(node)] = comm_id

    with gzip.open(p, "wt", encoding="utf-8") as f:
        json.dump(node_to_community, f, indent=2)


def read_communitydict(p: Path) -> dict[int, int]:
    with gzip.open(p, "rt", encoding="utf-8") as f:
        node_to_community_raw = json.load(f)
        node_to_community = {int(k): v for k, v in node_to_community_raw.items()}
    return node_to_community

def communitydict_to_communitylist(comm_dict: dict[int, int]) -> list[list[int]]:
    # コミュニティIDごとにノードをリスト化
    comm_to_members = defaultdict(list)
    for node_id, comm_id in comm_dict.items():
        comm_to_members[comm_id].append(node_id)

    return list(comm_to_members.values())

# def load_graph_data(data_dir: Path) -> tuple[nx.DiGraph, list[list[int]]]:
#     """データファイル (.gzエッジリスト, JSON形式のコミュニティテキスト) を読み込む関数"""
#     edgelist_path = data_dir / "edgelist.txt.gz"
#     communities_path = data_dir / "communities.txt.gz"

#     # 1. 有向グラフの読み込み
#     G: nx.DiGraph = nx.read_edgelist(
#         edgelist_path,
#         nodetype=int,
#         create_using=nx.DiGraph
#     )

# def load_data(data_dir: Path):
#     """データファイル (.gzエッジリスト, JSON形式のコミュニティテキスト) を読み込む関数"""
#     edgelist_path = data_dir / "LFR_edgelist.txt.gz"
#     communities_path = data_dir / "LFR_communities.txt"

#     # 1. 有向グラフの読み込み (.gz圧縮ファイルを自動解凍)
#     G = nx.read_edgelist(
#         edgelist_path,
#         nodetype=int,
#         create_using=nx.DiGraph
#     )

#     # 2. JSON形式の LFR_communities.txt を読み込み
#     with open(communities_path, "r", encoding="utf-8") as f:
#         node_to_community_raw = json.load(f)
#         # キーを文字列から int 型に変換
#         node_to_community = {int(k): v for k, v in node_to_community_raw.items()}

#     # コミュニティIDごとにノードをリスト化
#     comm_to_members = defaultdict(list)
#     for node_id, comm_id in node_to_community.items():
#         comm_to_members[comm_id].append(node_id)

#     communities_list = list(comm_to_members.values())
#     return G, communities_list