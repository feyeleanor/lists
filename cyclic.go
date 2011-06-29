package lists

import "github.com/feyeleanor/chain"
import "github.com/feyeleanor/raw"

/*
	A CycList is a circular list structure.
	Each node in the list may point to exactly one other node in the list.
	No node may be pointed to by more than one other node in the list.
	There are no nil links between nodes in the list.
*/

//	A declarative method for building CycLists
func Loop(items... interface{}) (c *CycList) {
	c = NewCycList(&chain.Cell{})
	c.Concatenate(items)
	return
}

type CycList struct {
	ListHeader
}

func NewCycList(n chain.Node) *CycList {
	return &CycList{ NewListHeader(n) }
}

func (c CycList) Clone() (r *CycList) {
	r = &CycList{ *c.ListHeader.Clone() }
	if r.end != nil {
		r.end.Link(chain.NEXT_NODE, r.start)
	}
	return
}

//	Iterate over all elements of the list indefinitely
//	The only way to terminate iteration is by raising a panic() in the applied function
func (c CycList) Cycle(f func(interface{})) {
	for n := c.start; ; n = chain.Next(n) {
		f(n.Content())
	}
}

func (c CycList) index(i int) (r int) {
	switch {
	case c.length == 0:		r = 0
	case i > 0:				r = i % c.length
	case i < 0:				r = c.length + (i % c.length)
	}
	return
}

// Return the value stored at the given offset from the start of the list
func (c CycList) At(i int) interface{} {
	return c.ListHeader.At(c.index(i))
}

// Set the value stored at the given offset from the start of the list
func (c CycList) Set(i int, v interface{}) {
	c.ListHeader.Set(c.index(i), v)
}

func (c *CycList) Rotate(i int) {
	if c != nil && c.end != nil {
		c.end = c.end.MoveTo(c.index(i))
		c.start = chain.Next(c.end)
	}
}

func (c *CycList) Append(v interface{}) {
	c.ListHeader.Append(v)
	if c.end == c.start {
		c.end.Link(chain.NEXT_NODE, c.start)
	}
}

func (c *CycList) Concatenate(i interface{}) {
	c.ListHeader.Concatenate(i)
	if c.end != nil {
		c.end.Link(chain.NEXT_NODE, c.start)
	}
}

//	Determines if another object is equivalent to the CycList
//	Two CycLists are identical if they both have the same number of nodes, and the head of each node is the same
func (c CycList) Equal(o interface{}) (r bool) {
	switch o := o.(type) {
	case *CycList:		r = o != nil && c.ListHeader.Equal(o.ListHeader)
	case CycList:		r = c.ListHeader.Equal(o.ListHeader)
	default:			raw.Catch(func() {
							r = c.start.(Equatable).Equal(o)
						})
	}
	return 
}

//	Reverses the order in which elements of a CycList are traversed
func (c *CycList) Reverse() {
	if c != nil {
		c.ListHeader.Reverse()
		c.end.Link(chain.NEXT_NODE, c.start)
	}
}

func (c *CycList) Tail() {
	c.ListHeader.Tail()
	c.end = c.start
}