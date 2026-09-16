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

	for i := 0; i < len(adj); i++ {
		assert.Equal(t, numFollowersList[i], network.FollowerSize(adj, i))
	}
	assert.Equal(t, numFollowersList[0], 3)
	assert.Equal(t, numFollowersList[1], 0)
	assert.Equal(t, numFollowersList[2], 2)
	assert.Equal(t, numFollowersList[3], 1)
}
