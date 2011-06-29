package lists

import "github.com/feyeleanor/chain"
import "fmt"
import "testing"

func TestNewListNode(t *testing.T) {
	l := NewListHeader(&chain.Cell{})
	if n, ok := l.NewListNode(-1).(*chain.Cell); !ok {
		t.Fatalf("node should be of type *chain.Cell")
	} else if n.Content() != -1 {
		t.Fatalf("node should be 0 but is %v", n.Content())
	}
}

func TestEnforceBounds(t *testing.T) {
	ConfirmEnforceBounds := func(l *LinearList, start, end, expected_start, expected_end int) {
		title := fmt.Sprintf("%v.EnforcedBounds(%v, %v)", l, start, end)
		ok := l.EnforceBounds(&start, &end)
		switch {
		case !ok:						t.Fatalf("%v should be true", title)
		case start != expected_start:	t.Fatalf("%v start should be %v but is %v", title, expected_start, start)
		case end != expected_end:		t.Fatalf("%v end should be %v but is %v", title, expected_end, end)
		}
	}
	RefuteEnforceBounds := func(l *LinearList, start, end int) {
		title := fmt.Sprintf("%v.EnforcedBounds(%v, %v)", l, start, end)
		ok := l.EnforceBounds(&start, &end)
		if ok {
			t.Fatalf("%v should be false", title)
		}
	}

	RefuteEnforceBounds(List(), -1, -1)
	RefuteEnforceBounds(List(), -1, 0)
	RefuteEnforceBounds(List(), -1, 1)

	RefuteEnforceBounds(List(), 0, -1)
	RefuteEnforceBounds(List(), 0, 0)
	RefuteEnforceBounds(List(), 0, 1)

	RefuteEnforceBounds(List(), 1, -1)
	RefuteEnforceBounds(List(), 1, 0)
	RefuteEnforceBounds(List(), 1, 1)

	RefuteEnforceBounds(List(0), -1, -1)
	ConfirmEnforceBounds(List(0), -1, 0, 0, 0)
	ConfirmEnforceBounds(List(0), -1, 1, 0, 0)

	RefuteEnforceBounds(List(0), 0, -1)
	ConfirmEnforceBounds(List(0), 0, 0, 0, 0)
	ConfirmEnforceBounds(List(0), 0, 1, 0, 0)

	RefuteEnforceBounds(List(0), 1, -1)
	RefuteEnforceBounds(List(0), 1, 0)
	RefuteEnforceBounds(List(0), 1, 1)

	RefuteEnforceBounds(List(0, 1), -1, -1)
	ConfirmEnforceBounds(List(0, 1), -1, 0, 0, 0)
	ConfirmEnforceBounds(List(0, 1), -1, 1, 0, 1)
	ConfirmEnforceBounds(List(0, 1), -1, 2, 0, 1)

	RefuteEnforceBounds(List(0, 1), 0, -1)
	ConfirmEnforceBounds(List(0, 1), 0, 0, 0, 0)
	ConfirmEnforceBounds(List(0, 1), 0, 1, 0, 1)
	ConfirmEnforceBounds(List(0, 1), 0, 2, 0, 1)

	RefuteEnforceBounds(List(0, 1), 1, -1)
	RefuteEnforceBounds(List(0, 1), 1, 0)
	ConfirmEnforceBounds(List(0, 1), 1, 1, 1, 1)
	ConfirmEnforceBounds(List(0, 1), 1, 2, 1, 1)

	RefuteEnforceBounds(List(0, 1), 2, -1)
	RefuteEnforceBounds(List(0, 1), 2, 0)
	RefuteEnforceBounds(List(0, 1), 2, 1)
	RefuteEnforceBounds(List(0, 1), 2, 2)
}