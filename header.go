package lists

import "github.com/feyeleanor/chain"
import "fmt"
import "reflect"
import "strings"


type ListHeader struct {
	nodeType	reflect.Type
	start 		chain.Node
	end			chain.Node
	cache		cachedNode
	length		int
}

func NewListHeader(n chain.Node) ListHeader {
	t := reflect.TypeOf(n)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}
	return ListHeader{ nodeType: t }
}

func (l ListHeader) newListNode() chain.Node {
	return reflect.New(l.nodeType.Elem()).Interface().(chain.Node)
}

func (l ListHeader) NewListNode(value interface{}) (n chain.Node) {
	n = l.newListNode()
	n.Set(chain.CURRENT_NODE, value)
	return
}

func (l ListHeader) EnforceBounds(start, end *int) (ok bool) {
	if *start < 0 {
		*start = 0
	}

	if *end > l.length - 1 {
		*end = l.length - 1
	}

	if *end >= *start {
		ok = true
	}
	return
}

func (l *ListHeader) Erase() {
	l.start = nil
	l.end = nil
	l.length = 0
	l.cache.Clear()
}

func (l ListHeader) String() (t string) {
	terms := []string{}
	l.Each(func(term interface{}) {
		terms = append(terms, fmt.Sprintf("%v", term))
	})
	if l.length > 0 && l.start == chain.Next(l.end) {
		terms = append(terms, "...")
	}
	t = strings.Join(terms, " ")
	t = strings.Replace(t, "()", "nil", -1)
	t = strings.Replace(t, "<nil>", "nil", -1)
	return "(" + t + ")"
}

func (l ListHeader) Len() (c int) {
	return l.length
}

func (l ListHeader) Start() chain.Node {
	return l.start
}

func (l ListHeader) End() chain.Node {
	return l.end
}

func (l ListHeader) Clone() (r *ListHeader) {
	r = &ListHeader{ nodeType: l.nodeType }
	l.Each(func(v interface{}) { r.Append(v) })
	return
}

func (l *ListHeader) Expand(i, n int) {
	if i > -1 && i <= l.length {
		switch {
		case l == nil:					fallthrough
		case i == l.length:				for ; n > 0; n-- {
											l.Append(l.newListNode())
										}

		case i == 0:					l.length = n
										for ; n > 0; n-- {
											x := l.newListNode()
											x.Link(chain.NEXT_NODE, l.start)
											l.start = x
										}

		default:						x1 := l.findNode(i - 1)
										x2 := l.findNode(i)
										l.length += n
										for ; n > 0; n-- {
											x1.Link(chain.NEXT_NODE, l.newListNode())
											x1 = chain.Next(x1)
										}
										x1.Link(chain.NEXT_NODE, x2)
		}
	}
}

func (l ListHeader) Each(f func(interface{})) {
	n := l.start
	for i := l.length; i > 0; i-- {
		f(n.Content())
		n = chain.Next(n)
	}
}

func (l ListHeader) equal(o ListHeader) (r bool) {
	if l.length == o.length {
		var e	Equatable

		r = true
		n := l.start
		x := o.start
		for i := l.length; r && i > 0; i-- {
			if e, r = n.(Equatable); r && e.Equal(x) {
				n = chain.Next(n)
				x = chain.Next(x)
			}
		}
	}
	return
}

func (l ListHeader) Equal(o interface{}) (r bool) {
	switch o := o.(type) {
	case *ListHeader:	r = o != nil && l.equal(*o)
	case ListHeader:	r = l.equal(o)
	}
	return
}

func (l *ListHeader) eachNode(f func(int, chain.Node)) {
	n := l.start
	for i := 0; i < l.length; i++ {
		f(i, n)
		n = chain.Next(n)
	}
}

func (l ListHeader) findNode(i int) (n chain.Node) {
	switch {
	case i == 0:				n = l.start
	case i == l.length - 1:		n = l.end
	default:					start, offset := l.cache.ClosestNode(i)
								if start == nil {
									start = l.start
								}
								if n = start.MoveTo(i - offset); n != nil {
									l.cache.Update(i, n)
								}
	}
	return
}

func (l ListHeader) At(i int) (r interface{}) {
	if n := l.findNode(i); n != nil {
		r = n.Content()
	}
	return
}

func (l ListHeader) Set(i int, v interface{}) {
	if n := l.findNode(i); n != nil {
		n.Set(chain.CURRENT_NODE, v)
	}
}

func (l ListHeader) Clear(i int) {
	l.Set(i, nil)
}

func (l *ListHeader) Append(v interface{}) {
	switch {
	case l.start == nil:	l.start = l.NewListNode(v)
							l.end = l.start

	default:				tail := chain.Next(l.end)
							l.end.Link(chain.NEXT_NODE, l.NewListNode(v))
							l.end = chain.Next(l.end)
							l.end.Link(chain.NEXT_NODE, tail)
	}
	l.length++
}

func (l *ListHeader) Concatenate(s interface{}) {
	switch s := s.(type) {
	case []interface{}:		if length := len(s); length > 0 {
								l.Append(s[0])
								if length > 1 {
									tail := chain.Next(l.end)
									for _, v := range s[1:] {
										l.end.Link(chain.NEXT_NODE, l.NewListNode(v))
										l.end = chain.Next(l.end)
									}
									l.end.Link(chain.NEXT_NODE, tail)
									l.length += length - 1
								}
							}

	case Sequence:			if length := s.Len(); length > 0 {
								l.Append(s.At(0))
								if length > 1 {
									tail := chain.Next(l.end)
									for i := 1; i < length; i++ {
										l.end.Link(chain.NEXT_NODE, l.NewListNode(s.At(i)))
										l.end = chain.Next(l.end)
									}
									l.end.Link(chain.NEXT_NODE, tail)
									l.length += length - 1
								}
							}

	default:				switch s := reflect.ValueOf(s); s.Kind() {
							case reflect.Slice:				if length := s.Len(); length > 0 {
																l.Append(s.Index(0).Interface())
																if length > 1 {
																	tail := chain.Next(l.end)
																	for i := 1; i < length; i++ {
																		l.end.Link(chain.NEXT_NODE, l.NewListNode(s.Index(i).Interface()))
																		l.end = chain.Next(l.end)
																	}
																	l.end.Link(chain.NEXT_NODE, tail)
																	l.length += length - 1
																}
															}
							}
	}
}

//	Iterates through the list reducing the nesting of each element which can be flattened.
//	Elements which are themselves LinearLists will be inlined as part of the containing list and their contained list destroyed.
func (l *ListHeader) Flatten() {
	l.eachNode(func(i int, n chain.Node) {
		value := n.Content()
		if h, ok := value.(Flattenable); ok {
			h.Flatten()
		}

		if h, ok := value.(Linkable); ok {
			switch length := h.Len(); {
			case length == 0:		n.Set(chain.CURRENT_NODE, nil)

			case length == 1:		n.Set(chain.CURRENT_NODE, h.Start().Content())

			default:				l.length += length - 1
									h.End().Link(chain.NEXT_NODE, chain.Next(n))
									n.Link(chain.CURRENT_NODE, h.Start())
									if n == l.start {
										l.start = h.Start()
									}

									if n == l.end {
										l.end = h.End()
									}
			}
		} else {
			n.Set(chain.CURRENT_NODE, value)
		}
	})
}

func (l ListHeader) Compact() []interface{} {
	s := make([]interface{}, l.Len(), l.Len())
	i := 0
	l.Each(func(v interface{}) {
		s[i] = v
		i++
	})
	return s
}

//	Reverses the order in which elements of a List are traversed
func (l *ListHeader) Reverse() {
	if l != nil {
		current := l.start
		l.end = current

		for i := l.length; i > 0; i-- {
			next := chain.Next(current)
			current.Link(chain.NEXT_NODE, l.start)
			l.start = current
			current = next				
		}
	}
	return
}

func (l ListHeader) Head() (r interface{}) {
	if l.start != nil {
		l.start.Content()
	}
	return
}

func (l *ListHeader) Tail() {
	if n := l.start; n != nil {
		l.start = chain.Next(n)
		n.Link(chain.NEXT_NODE, nil)
		l.length--
	}
}