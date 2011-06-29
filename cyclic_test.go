package lists

import "github.com/feyeleanor/chain"
import "testing"

func TestCycListLen(t *testing.T) {
	ConfirmLen := func(c *CycList, x int) {
		if i := c.Len(); i != x {
			t.Fatalf("'%v' length should be %v but is %v", c.String(), x, i)
		}
	}
	ConfirmLen(Loop(), 0)
	ConfirmLen(Loop(4), 1)
	ConfirmLen(Loop(4, 3, 2, 1), 4)
	ConfirmLen(Loop(4, Loop(3), 2, 1), 4)
}

func TestCycListClone(t *testing.T) {
	ConfirmClone := func(c, r *CycList) {
		x := c.Clone()
		if !x.Equal(r) {
			t.Fatalf("%v should be %v", x, r)
		}
	}
	ConfirmClone(Loop(), Loop())
	ConfirmClone(Loop(0), Loop(0))
	ConfirmClone(Loop(0, 1), Loop(0, 1))
	ConfirmClone(Loop(0, Loop(0, 1)), Loop(0, Loop(0, 1)))
}

func TestCycListEach(t *testing.T) {
	c := Loop(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	count := 0
	c.Each(func(i interface{}) {
		if i != count {
			t.Fatalf("element %v erroneously reported as %v", count, i)
		}
		count++
	})
	if count != c.length {
		t.Fatalf("loop length %v erroneously reported iterations as %v", c.length, count)
	}
}

func TestCycListCycle(t *testing.T) {
	c := Loop(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	count := 0
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("element %v erroneously reported as %v", count, x)
		}
	}()
	c.Cycle(func(i interface{}) {
		if i != count {
			panic(i)
		}
		count++
		if count == c.Len() {
			panic(nil)
		}
	})
}

func TestCycListAt(t *testing.T) {
	ConfirmAt := func(c *CycList, i int, v interface{}) {
		if x := c.At(i); x != v {
			t.Fatalf("%v.At(%v) should be %v but is %v", c, i, v, x)
		}
	}
	c := Loop(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	ConfirmAt(c, -32, 18)
	ConfirmAt(c, -21, 19)
	ConfirmAt(c, -10, 10)
	ConfirmAt(c, -1, 19)
	ConfirmAt(c, 0, 10)
	ConfirmAt(c, 9, 19)
	ConfirmAt(c, 10, 10)
	ConfirmAt(c, 21, 11)
	ConfirmAt(c, 32, 12)
}

func TestCycListSet(t *testing.T) {
	ConfirmSet := func(c *CycList, i int, v interface{}) {
		c.Set(i, v)
		if x := c.At(i); x != v {
			t.Fatalf("%v.Set(%v) should be %v but is %v", c, i, v, x)
		}
	}
	c := Loop(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	ConfirmSet(c, -21, 0)
	ConfirmSet(c, -10, -10)
	ConfirmSet(c, -1, -1)
	ConfirmSet(c, 0, 10)
	ConfirmSet(c, 9, 10)
	ConfirmSet(c, 11, 11)
	ConfirmSet(c, 22, 12)
	ConfirmSet(c, 33, 13)
}

func TestCycListRotate(t *testing.T) {
	ConfirmRotate := func(c *CycList, i int, r *CycList) {
		c.Rotate(i)
		if !c.Equal(r) {
			t.Fatalf("%v should be %v", c, r)
		}
	}
	ConfirmRotate(Loop(), 0, Loop())
	ConfirmRotate(Loop(), 1, Loop())
	ConfirmRotate(Loop(), 2, Loop())
	ConfirmRotate(Loop(0), 0, Loop(0))
	ConfirmRotate(Loop(0), 1, Loop(0))
	ConfirmRotate(Loop(0), 2, Loop(0))

	ConfirmRotate(Loop(0, 1, 2, 3), 0, Loop(0, 1, 2, 3))

	ConfirmRotate(Loop(0, 1, 2, 3), 1, Loop(1, 2, 3, 0))
	ConfirmRotate(Loop(0, 1, 2, 3), -3, Loop(1, 2, 3, 0))

	ConfirmRotate(Loop(0, 1, 2, 3), 2, Loop(2, 3, 0, 1))
	ConfirmRotate(Loop(0, 1, 2, 3), -2, Loop(2, 3, 0, 1))

	ConfirmRotate(Loop(0, 1, 2, 3), 3, Loop(3, 0, 1, 2))
	ConfirmRotate(Loop(0, 1, 2, 3), -1, Loop(3, 0, 1, 2))

	ConfirmRotate(Loop(0, 1, 2, 3), 4, Loop(0, 1, 2, 3))
	ConfirmRotate(Loop(0, 1, 2, 3), -4, Loop(0, 1, 2, 3))
}

func TestCycListAppend(t *testing.T) {
	ConfirmAppend := func(c *CycList, v interface{}, r *CycList) {
		c.Append(v)
		if !c.Equal(r) {
			t.Fatalf("%v should be %v", c, r)
		}
	}
	ConfirmAppend(Loop(), 1, Loop(1))
	ConfirmAppend(Loop(1), 2, Loop(1, 2))
}

func TestCycListConcatenate(t *testing.T) {
	ConfirmConcatenate := func(c *CycList, s interface{}, r *CycList) {
		c.Concatenate(s)
		if !c.Equal(r) {
			t.Fatalf("%v should be %v", c, r)
		}
	}
	ConfirmConcatenate(Loop(), []interface{}{}, Loop())
	ConfirmConcatenate(Loop(), []interface{}{1}, Loop(1))
	ConfirmConcatenate(Loop(1), []interface{}{2, 3}, Loop(1, 2, 3))
}

func TestCycListString(t *testing.T) {
	ConfirmFormat := func(c *CycList, x string) {
		if s := c.String(); s != x {
			t.Fatalf("'%v' erroneously serialised as '%v'", x, s)
		}
	}

	ConfirmFormat(Loop(), "()")
	ConfirmFormat(Loop(0), "(0 ...)")
	ConfirmFormat(Loop(0, nil), "(0 nil ...)")
	ConfirmFormat(Loop(0, Loop(0)), "(0 (0 ...) ...)")
	ConfirmFormat(Loop(1, Loop(0, nil)), "(1 (0 nil ...) ...)")

	ConfirmFormat(Loop(1, 0, nil), "(1 0 nil ...)")

	r := Loop(10, Loop(0, Loop(0)))
	ConfirmFormat(r, "(10 (0 (0 ...) ...) ...)")
	r.Rotate(chain.NEXT_NODE)
	ConfirmFormat(r, "((0 (0 ...) ...) 10 ...)")
	ConfirmFormat(r.start.Content().(*CycList), "(0 (0 ...) ...)")

	r = Loop(r, 0, Loop(-1, -2, r))
	ConfirmFormat(r, "(((0 (0 ...) ...) 10 ...) 0 (-1 -2 ((0 (0 ...) ...) 10 ...) ...) ...)")
}

func TestLoop(t *testing.T) {
	ConfirmFormat := func(c *CycList, x string) {
		if s := c.String(); s != x {
			t.Fatalf("'%v' erroneously serialised as '%v'", x, s)
		}
	}
	ConfirmFormat(Loop(), "()")
	ConfirmFormat(Loop(1), "(1 ...)")
	ConfirmFormat(Loop(2, 1), "(2 1 ...)")
	ConfirmFormat(Loop(3, 2, 1), "(3 2 1 ...)")
	ConfirmFormat(Loop(4, 3, 2, 1), "(4 3 2 1 ...)")

	c := Loop(4, 3, 2, 1)
	ConfirmFormat(c, "(4 3 2 1 ...)")
	ConfirmFormat(Loop(5, c, 0), "(5 (4 3 2 1 ...) 0 ...)")
	c = Loop(5, c, 0)
	ConfirmFormat(c, "(5 (4 3 2 1 ...) 0 ...)")
}

func TestCycListReverse(t *testing.T) {
	ConfirmReverse := func(c, r *CycList) {
		c.Reverse()
		if !c.Equal(r) {
			t.Fatalf("%v should be %v", c, r)
		}
	}

	c := Loop(1)
	ConfirmReverse(c, Loop(1))
	ConfirmReverse(c, Loop(1))

	ConfirmReverse(Loop(1, 2), Loop(2, 1))
	ConfirmReverse(Loop(2, 1), Loop(1, 2))

	c = Loop(1, 2)
	ConfirmReverse(c, Loop(2, 1))
	ConfirmReverse(c, Loop(1, 2))

	ConfirmReverse(Loop(1, 2, 3), Loop(3, 2, 1))
	ConfirmReverse(Loop(3, 2, 1), Loop(1, 2, 3))

	c = Loop(1, 2, 3)
	ConfirmReverse(c, Loop(3, 2, 1))
	ConfirmReverse(c, Loop(1, 2, 3))

	ConfirmReverse(Loop(1, 2, 3, 4), Loop(4, 3, 2, 1))
	ConfirmReverse(Loop(4, 3, 2, 1), Loop(1, 2, 3, 4))
}

func TestCycListFlatten(t *testing.T) {
	ConfirmFlatten := func(c, r *CycList) {
		c.Flatten()
		if !c.Equal(r) {
			t.Fatalf("%v should be %v", c, r)
		}
	}
	ConfirmFlatten(Loop(), Loop())
	ConfirmFlatten(Loop(1), Loop(1))
	ConfirmFlatten(Loop(1, Loop(2)), Loop(1, 2))
	ConfirmFlatten(Loop(1, Loop(2, Loop(3))), Loop(1, 2, 3))

	ConfirmFlatten(Loop(0, List(1)), Loop(0, 1))
	ConfirmFlatten(Loop(0, List(1), 2), Loop(0, 1, 2))
	ConfirmFlatten(Loop(0, List(1, 2), 3), Loop(0, 1, 2, 3))
	ConfirmFlatten(Loop(0, List(1, List(2, 3), 4, List(5, List(6, 7)))), Loop(0, 1, 2, 3, 4, 5, 6, 7))

	ConfirmFlatten(Loop(0, List(1, Loop(2, 3))), Loop(0, 1, 2, 3))
	ConfirmFlatten(Loop(0, List(1, Loop(2, 3, List(4, Loop(5))))), Loop(0, 1, 2, 3, 4, 5))
}

func TestCycListCompact(t *testing.T) {
	t.Fatal()
	ConfirmCompact := func(c *CycList, r interface{}) {
//		if x := c.Compact(); !r.Equal(x) {
//			t.Fatalf("%v.Compact() should be %v but is %v", c, r, x)
//		}
	}

	ConfirmCompact(Loop(), []interface{}{})
	ConfirmCompact(Loop(1), []interface{}{ 1 })
	ConfirmCompact(Loop(0, 1, 2, 3, 4, 5, 6, 7, 8, 9), []interface{}{ 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 })
}