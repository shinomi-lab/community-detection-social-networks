import json
import os
from pathlib import Path

import networkx as nx


def write_graph(p: Path):
    edges = [
        (0, 1),
        (0, 2),
        (1, 2),
        (1, 3),
        (2, 3),
    ]

    g = nx.DiGraph()
    for (i, j) in edges:
      g.add_edge(i, j)

    # `data=False`で容量削減
    nx.write_edgelist(g, p.joinpath("edgelist.txt"), data=False)

    # 末尾に".gz"をつけてgzip圧縮をしてさらに容量削減
    nx.write_edgelist(g, p.joinpath("edgelist.txt.gz"), data=False)

    # 検証用の隣接行列JSON
    n = g.number_of_nodes()
    adj: dict[int, dict[int, int]] = {}
    for i in range(n):
      for j in range(n):
        a = adj.setdefault(i, {})
        a[j] = 1 if g.has_edge(i, j) else 0

    with open(p.joinpath("adj.json"), "w") as f:
        json.dump(adj, f)


def read_graph(p: Path):
    g: nx.DiGraph = nx.read_edgelist(
        p.joinpath("edgelist.txt"), create_using=nx.DiGraph
    )

    with open(p.joinpath("adj.json")) as f:
      adj = json.load(f)

    for i in g.nodes():
      for j in g.nodes():
        # print(i, j, g.has_edge(i, j), adj[i][j])
        if g.has_edge(i, j):
          assert adj[i][j] == 1
        else:
          assert adj[i][j] == 0


if __name__ == "__main__":
    dir = Path(os.getcwd()).parent.joinpath("difftools/network")
    print(dir)
    write_graph(dir)
    read_graph(dir)
