package main

import (
	"difftools/network"
)

func main() {
	net := network.ReadJson("cmd/test/sample_adj.json")

	for i, v := range net.Followers {
		println(i, v)
	}
	// println(n)
	// // _ = adj
	// fmt.Println(adj)
}
