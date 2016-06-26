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
