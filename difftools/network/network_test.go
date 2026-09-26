package network_test

import (
	"difftools/network"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNumFollowersList(t *testing.T) {
	adj := [][]int{
		{1, 1, 1, 0},
		{0, 0, 0, 0},
		{1, 1, 0, 0},
		{1, 0, 0, 0},
	}
	numFollowersList := network.GetFollowerNums(adj)

	for i := range adj {
		assert.Equal(t, numFollowersList[i], network.FollowerSize(adj, i))
	}
	assert.Equal(t, numFollowersList[0], 3)
	assert.Equal(t, numFollowersList[1], 0)
	assert.Equal(t, numFollowersList[2], 2)
	assert.Equal(t, numFollowersList[3], 1)
}

func TestReadEdgeList(t *testing.T) {
	net1 := network.ReadEdgelist("testdata/edgelist.txt")
	net2 := network.ReadAdjJson("testdata/adj.json")

	assert.Equal(t, net1.N, net2.N)
	for i := range net1.N {
		for j := range net2.N {
			assert.Equal(t, net1.Adj[i][j], net2.Adj[i][j])
		}
		assert.Equal(t, net1.Followers[i], net2.Followers[i])
	}
}

func TestReadEdgeListGz(t *testing.T) {
	net1 := network.ReadEdgelist("testdata/edgelist.txt")
	net2 := network.ReadEdgelist("testdata/edgelist.txt.gz")

	assert.Equal(t, net1.N, net2.N)
	for i := range net1.N {
		for j := range net2.N {
			assert.Equal(t, net1.Adj[i][j], net2.Adj[i][j])
		}
		assert.Equal(t, net1.Followers[i], net2.Followers[i])
	}
}
