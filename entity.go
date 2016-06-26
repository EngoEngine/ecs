package ecs

import (
	"sync"
	"sync/atomic"
)

var (
	counterLock sync.Mutex
	idInc       uint64
)

// Entity is the E in Entity Component System. It belongs to any amount of
// Systems, and has a number of Components
type BasicEntity struct {
	id uint64
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

// NewBasic creates a new Entity with a new unique identifier - can be called across multiple goroutines
func NewBasic() BasicEntity {
	return BasicEntity{id: atomic.AddUint64(&idInc, 1)}
}

// NewBasics creates an amount of new entities with a new unique identifier - can be called across multiple goroutines
// Performs better than NewBasic for large numbers of entities.
func NewBasics(amount int) []BasicEntity {
	entities := make([]BasicEntity, amount)

	counterLock.Lock()
	for i := 0; i < amount; i++ {
		idInc++
		entities[i] = BasicEntity{id: idInc}
	}
	counterLock.Unlock()

	return entities
}

func (e BasicEntity) ID() uint64 {
	return e.id
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
