package ecs

import (
	"sync/atomic"
)

var (
	idInc uint64
)

// A BasicEntity is simply a set of components with a unique ID attached to it,
// nothing more. It belongs to any amount of Systems, and has a number of
// Components
type BasicEntity struct {
	// Entity ID.
	id       uint64
	parent   *BasicEntity
	children []*BasicEntity
}

// Identifier is an interface for anything that implements the basic ID() uint64,
// as the BasicEntity does.  It is useful as more specific interface for an
// entity registry than just the interface{} interface
type Identifier interface {
	ID() uint64
}

// IdentifierSlice implements the sort.Interface, so you can use the
// store entites in slices, and use the P=n*log n lookup for them
type IdentifierSlice []Identifier

// NewBasic creates a new Entity with a new unique identifier. It is safe for
// concurrent use.
func NewBasic() BasicEntity {
	return BasicEntity{id: atomic.AddUint64(&idInc, 1)}
}

// NewBasics creates an amount of new entities with a new unique identifiers. It
// is safe for concurrent use, and performs better than NewBasic for large
// numbers of entities.
func NewBasics(amount int) []BasicEntity {
	entities := make([]BasicEntity, amount)

	lastID := atomic.AddUint64(&idInc, uint64(amount))
	for i := 0; i < amount; i++ {
		entities[i].id = lastID - uint64(amount) + uint64(i) + 1
	}

	return entities
}

// ID returns the unique identifier of the entity.
func (e BasicEntity) ID() uint64 {
	return e.id
}

// GetBasicEntity returns a Pointer to the BasicEntity itself
// By having this method, All Entities containing a BasicEntity now automatically have a GetBasicEntity Method
// This allows system.Add functions to recieve a single interface
// EG:
// s.AddByInterface(a interface{GetBasicEntity()*BasicEntity, GetSpaceComponent()*SpaceComponent){
// s.Add(a.GetBasicEntity(),a.GetSpaceComponent())
//}
func (e *BasicEntity) GetBasicEntity() *BasicEntity {
	return e
}

// AppendChild appends a child to the BasicEntity
func (e *BasicEntity) AppendChild(child *BasicEntity) {
	child.parent = e
	e.children = append(e.children, child)
}

func (e *BasicEntity) RemoveChild(child *BasicEntity) {
	delete := -1
	for i, v := range e.children {
		if v.ID() == child.ID() {
			delete = i
			break
		}
	}
	if delete >= 0 {
		e.children = append(e.children[:delete], e.children[delete+1:]...)
	}
}

// Children returns the children of the BasicEntity
func (e *BasicEntity) Children() []BasicEntity {
	ret := []BasicEntity{}
	for _, child := range e.children {
		ret = append(ret, *child)
	}
	return ret
}

// var visited map[uint64]struct{}
// var descs []*BasicEntity

// Descendents returns the children and their children all the way down the tree.
func (e *BasicEntity) Descendents() []BasicEntity {
	return descendents([]BasicEntity{}, e, e)
}

func descendents(in []BasicEntity, this, top *BasicEntity) []BasicEntity {
	for _, child := range this.children {
		in = descendents(in, child, top)
	}
	if this.ID() == top.ID() {
		return in
	}
	return append(in, *this)
}

// Parent returns the parent of the BasicEntity
func (e *BasicEntity) Parent() *BasicEntity {
	return e.parent
}

// Len returns the length of the underlying slice
// part of the sort.Interface
func (is IdentifierSlice) Len() int {
	return len(is)
}

// Less will return true if the ID of element at i is less than j;
// part of the sort.Interface
func (is IdentifierSlice) Less(i, j int) bool {
	return is[i].ID() < is[j].ID()
}

// Swap the elements at positions i and j
// part of the sort.Interface
func (is IdentifierSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

// BasicFase is an interface that BasicEntity and entities containing
// a BasicEntity implement.
type BasicFace interface {
	GetBasicEntity() *BasicEntity
}
