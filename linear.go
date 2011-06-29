package lists

import "github.com/feyeleanor/chain"

/*
	A LinearList is a finitely-terminated list structure.
	Each node in the list may point to exactly one other node in the list.
	The terminating node does not point to any other node.
	No node may be pointed to by more than one other node in the list.
	There are no nil links between nodes in the list.
*/

func List(items... interface{}) (l *LinearList) {
	l = NewLinearList(&chain.Cell{})
	l.Concatenate(items)
	return
}

type LinearList struct {
	ListHeader
}

func NewLinearList(n chain.Node) *LinearList {
	return &LinearList{ NewListHeader(n) }
}

func (l LinearList) End() chain.Node {
	return l.end
}

func (l LinearList) Clone() *LinearList {
	return &LinearList{ *l.ListHeader.Clone() }
}

//	Determines if another object is equivalent to the LinearList
//	Two LinearLists are identical if they both have the same number of nodes, and the head of each node is the same
func (l LinearList) Equal(o interface{}) (r bool) {
	switch o := o.(type) {
	case *LinearList:	r = o != nil && l.ListHeader.Equal(o.ListHeader)
	case LinearList:	r = l.ListHeader.Equal(o.ListHeader)
	}
	return 
}

//	Removes all elements in the range from the list.
func (l *LinearList) Delete(from, to int) {
	if l != nil && l.EnforceBounds(&from, &to) {
		last_element_index := l.length - 1
		switch {
		case from == 0:						switch {
											case to == 0:					l.start = chain.Next(l.start)
																			l.length--

											case to == last_element_index:	l.start = nil
																			l.end = nil
																			l.length = 0

											default:						l.start = l.findNode(to + 1)
																			l.length -= to + 1
											}

		case from == last_element_index:	l.end = l.findNode(from - 1)
											l.end.Link(chain.NEXT_NODE, nil)
											l.length--

		case to == last_element_index:		l.end = l.findNode(from - 1)
											l.end.Link(chain.NEXT_NODE, nil)
											l.length = from


		case from == to:					s := l.findNode(from - 1)
											s.Link(chain.NEXT_NODE, s.MoveTo(2))
											l.length--

		default:							e := l.findNode(from - 1)
											e.Link(chain.NEXT_NODE, e.MoveTo(to - from + 2))
											l.length -= to - from + 1
		}
	}
}

//	Removes the elements in the range from the current list and returns a new list containing them.
func (l *LinearList) Cut(start, end int) (r LinearList) {
	if l != nil {
		r.nodeType = l.nodeType
		if ok := l.EnforceBounds(&start, &end); ok {
			last_element_index := l.length - 1
			switch {
			case start == 0:					r.start = l.start
												if end == last_element_index {
													r.end = l.end
													l.start = nil
													l.end = nil
												} else {
													r.end = l.findNode(end)
													l.start = l.findNode(end + 1)
												}

			case start == last_element_index:	r.start = l.end
												l.end = l.findNode(start - 1)

			case end == last_element_index:		l.end = l.findNode(start - 1)
												r.start = l.findNode(start)

			case start == end:					s := l.findNode(start - 1)
												r.start = l.findNode(start)
												r.end = r.start
												s.Link(chain.NEXT_NODE, l.findNode(start + 1))

			default:							s := l.findNode(start - 1)
												r.start = l.findNode(start)
												r.end = l.findNode(end)
												s.Link(chain.NEXT_NODE, l.findNode(end + 1))
			}
			l.cache.Clear()
			if r.end != nil {
				r.end.Link(chain.NEXT_NODE, nil)
			}
			r.length = end - start + 1

			if l.end != nil {
				l.end.Link(chain.NEXT_NODE, nil)
			}
			l.length -= r.length
		}
	}
	return
}

//	Insert an item into the list at the given location.
func (l *LinearList) Insert(i int, o interface{}) {
	if i > -1 && i <= l.length {
		switch {
		case l == nil:					fallthrough
		case i == l.length:				l.Append(o)

		case i == 0:					n := l.NewListNode(o)
										n.Link(chain.NEXT_NODE, l.start)
										l.start = n
										l.length++

		default:						n1 := l.findNode(i - 1)
										n2 := l.findNode(i)
										n1.Link(chain.NEXT_NODE, l.NewListNode(o))
										chain.Next(n1).Link(chain.NEXT_NODE, n2)
										l.length++
		}
	}
}

//	Take all the elements from another list and insert them into this list, destroying the other list if successful.
func (l *LinearList) Absorb(i int, o *LinearList) (ok bool) {
	if o != nil && i > -1 && i <= l.length {
		switch {
		case l == nil:					*l = *o

		case i == 0:					o.end.Link(chain.NEXT_NODE, l.start)
										l.start = o.start
										l.length += o.length

		case i == l.length:				l.end.Link(chain.NEXT_NODE, o.start)
										l.end = o.end
										l.length += o.length

		default:						n := l.findNode(i - 1)
										o.end.Link(chain.NEXT_NODE, l.findNode(i))
										n.Link(chain.NEXT_NODE, o.start)
										l.length += o.length
		}
		o.Erase()
	 	ok = true
	}
	return
}