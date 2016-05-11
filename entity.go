package ecs

import (
	"log"
	"math"
	"sync"
)

var (
	counterLock sync.Mutex
	idInc       uint64
)

// A BasicEntity is simply a set of components with a unique ID attached to it,
// nothing more. It belongs to any amount of Systems, and has a number of
// Components
type BasicEntity struct {
	// Entity ID.
	id uint64
}

// NewBasic creates a new Entity with a new unique identifier. It is safe for
// concurrent use.
func NewBasic() BasicEntity {
	counterLock.Lock()
	idInc++
	if idInc >= math.MaxUint64 {
		log.Println("Warning: id overload")
		idInc = 1
	}
	counterLock.Unlock()
	return BasicEntity{idInc}
}

// NewBasics creates an amount of new entities with a new unique identifiers. It
// is safe for concurrent use, and performs better than NewBasic for large
// numbers of entities.
func NewBasics(amount int) []BasicEntity {
	entities := make([]BasicEntity, amount)

	counterLock.Lock()
	for i := 0; i < amount; i++ {
		idInc++
		if idInc >= math.MaxUint64 {
			log.Println("Warning: id overload")
			idInc = 1
		}
		entities[i] = BasicEntity{idInc}
	}
	counterLock.Unlock()

	return entities
}

// ID returns the unique identifier of the entity.
func (e BasicEntity) ID() uint64 {
	return e.id
}
