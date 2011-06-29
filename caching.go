package lists

import "github.com/feyeleanor/chain"

type cachedNode struct {
	chain.Node
	index	int
}

func (c cachedNode) Update(i int, node chain.Node) {
	c.index = i
	c.Node = node
}

func (c cachedNode) Clear() {
	c.index = 0
	c.Node = nil
}

func (c cachedNode) ClosestNode(i int) (node chain.Node, offset int) {
	if i > c.index {
		return c.Node, c.index
	}
	return
}