from collections import Counter, defaultdict
import json
from pathlib import Path
import gzip
import networkx as nx
import numpy as np
import pandas as pd


def load_graph_data(data_dir: Path) -> tuple[nx.DiGraph, list[list[int]]]:
    """データファイル (.gzエッジリスト, JSON形式のコミュニティテキスト) を読み込む関数"""
    edgelist_path = data_dir / "LFR_edgelist.txt.gz"
    communities_path = data_dir / "LFR_communities.txt.gz"

    # 1. 有向グラフの読み込み
    G = nx.read_edgelist(
        edgelist_path,
        nodetype=int,
        create_using=nx.DiGraph
    )

    # 2. JSON形式の LFR_communities.txt を読み込み
    with gzip.open(communities_path, "rt", encoding="utf-8") as f:
        node_to_community_raw = json.load(f)
        node_to_community = {int(k): v for k, v in node_to_community_raw.items()}

    # コミュニティIDごとにノードをリスト化
    comm_to_members = defaultdict(list)
    for node_id, comm_id in node_to_community.items():
        comm_to_members[comm_id].append(node_id)

    communities_list = list(comm_to_members.values())
    return G, communities_list

