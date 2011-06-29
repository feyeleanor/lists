package lists

import "github.com/feyeleanor/chain"

type Equatable interface {
	chain.Equatable
}

type Flattenable interface {
	Flatten()
}

type Linear interface {
	chain.Linear
}

type Iterable interface {
	Each(func(interface{}))
}

type Linkable interface {
	Linear
	Start() chain.Node
	End() chain.Node
}

type Sequence interface {
	Linear
	Iterable
	At(int) interface{}
	Set(int, interface{})
}