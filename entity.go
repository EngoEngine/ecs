package ecs

import (
	"log"
	"math"
	"sync"
)

var (
	counterLock sync.Mutex
	id_incr     uint64
)

// Entity is the E in Entity Component System. It belongs to any amount of
// Systems, and has a number of Components
type BasicEntity struct {
	id uint64
}

// NewBasic creates a new Entity with a new unique identifier - can be called across multiple goroutines
func NewBasic() BasicEntity {
	counterLock.Lock()
	id_incr++
	if id_incr >= math.MaxUint64 {
		log.Println("Warning: id overload")
		id_incr = 1
	}
	counterLock.Unlock()
	return BasicEntity{id_incr}
}

// NewBasics creates an amount of new entities with a new unique identifier - can be called across multiple goroutines
// Performs better than NewBasic for large numbers of entities.
func NewBasics(amount int) []BasicEntity {
	entities := make([]BasicEntity, amount)

	counterLock.Lock()
	for i := 0; i < amount; i++ {
		id_incr++
		if id_incr >= math.MaxUint64 {
			log.Println("Warning: id overload")
			id_incr = 1
		}
		entities[i] = BasicEntity{id_incr}
	}
	counterLock.Unlock()

	return entities
}

func (e BasicEntity) ID() uint64 {
	return e.id
}
