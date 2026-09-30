from collections import Counter, defaultdict
import json
from pathlib import Path

import networkx as nx
import numpy as np
import pandas as pd


def load_data(data_dir: Path):
    """データファイル (.gzエッジリスト, JSON形式のコミュニティテキスト) を読み込む関数"""
    edgelist_path = data_dir / "LFR_edgelist.txt.gz"
    communities_path = data_dir / "LFR_communities.txt"

    # 1. 有向グラフの読み込み
    G = nx.read_edgelist(
        edgelist_path,
        nodetype=int,
        create_using=nx.DiGraph
    )

    # 2. JSON形式の LFR_communities.txt を読み込み
    with open(communities_path, "r", encoding="utf-8") as f:
        node_to_community_raw = json.load(f)
        node_to_community = {int(k): v for k, v in node_to_community_raw.items()}

    # コミュニティIDごとにノードをリスト化
    comm_to_members = defaultdict(list)
    for node_id, comm_id in node_to_community.items():
        comm_to_members[comm_id].append(node_id)

    communities_list = list(comm_to_members.values())
    return G, communities_list


def analyze_competitor_scenarios(G, communities, global_pr_scores, target_id=2):
    """
    Community target_id と他のコミュニティとの繋がり、および相手の構造特性を同時に解析する
    """
    if target_id >= len(communities):
        raise ValueError(f"target_id ({target_id}) がコミュニティ総数 ({len(communities)}) を超えています。")

    target_nodes = set(communities[target_id])
    scenario_data = []
    bridge_details = {}

    header_text = f"--- 💡 コミュニティ {target_id} の隣人構造解析 ---\n"

    for i, comm_members in enumerate(communities):
        if i == target_id:
            continue

        current_nodes = set(comm_members)

        # 1. 境界エッジの抽出（パイプの太さ）
        inter_edges = list(nx.edge_boundary(G, target_nodes, current_nodes))
        edge_count = len(inter_edges)

        if edge_count == 0:
            continue  # 完全に独立しているコミュニティは除外

        # 2. 両側のブリッジノード（架け橋）を特定
        target_side_bridges = set(u for u, v in inter_edges)
        partner_side_bridges = set(v for u, v in inter_edges)

        # 3. 相手コミュニティ自体の「構造の強さ」を全体PRから計算
        partner_pr_values = [global_pr_scores[n] for n in comm_members if n in G]
        partner_mean = np.mean(partner_pr_values) if partner_pr_values else 0.0

        # 4. 最大のゲートキーパーの特定 (Target側)
        my_nodes_in_edges = [u for u, v in inter_edges]
        if my_nodes_in_edges:
            gatekeeper, gk_edges = Counter(my_nodes_in_edges).most_common(1)[0]
        else:
            gatekeeper, gk_edges = None, 0

        scenario_data.append({
            "Opponent ID": i,
            "Opponent Size (規模)": len(current_nodes),
            "Opponent Mean (構造の強さ)": partner_mean,
            "Inter Edges (パイプの太さ)": edge_count,
            "Target側ブリッジ数": len(target_side_bridges),
            "Opponent側ブリッジ数": len(partner_side_bridges),
            "Target側最大ゲートキーパーノード": gatekeeper,
            "その1人が持つ接続エッジ数": gk_edges
        })

        bridge_details[i] = inter_edges

    df_scenarios = pd.DataFrame(scenario_data)
    if not df_scenarios.empty:
        df_scenarios = df_scenarios.sort_values(by="Inter Edges (パイプの太さ)", ascending=False).reset_index(drop=True)

    return df_scenarios, bridge_details, header_text


def main():
    # 実行ファイルの位置を基準に相対パスで data ディレクトリを指定
    base_dir = Path(__file__).resolve().parent.parent.parent
    data_dir = base_dir / "data"

    print(f"📁 データ参照元: {data_dir}")

    # データ読み込み
    G, communities = load_data(data_dir)

    # 全体 PageRank 計算
    full_pr = nx.pagerank(G)

    # 構造解析実行 (target_id=2)
    scenarios_df, _, header_text = analyze_competitor_scenarios(G, communities, full_pr, target_id=2)

    # 表示・出力用テキストの整形
    output_text = (
        f"{header_text}\n"
        "📊 【コミュニティ間影響力・隣接構造解析テーブル】\n"
        f"{scenarios_df.to_string(index=False)}\n"
    )

    # テキストファイルへの保存処理 (data/ ディレクトリへ)
    output_file_path = data_dir / "analyze_inter_community_edges_result.txt"
    with open(output_file_path, "w", encoding="utf-8") as f:
        f.write(output_text)

    print(f"💾 結果をファイルに保存しました: {output_file_path}")


if __name__ == "__main__":
    main()