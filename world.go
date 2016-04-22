package ecs

import (
	"sort"
)

// World contains a bunch of Entities, and a bunch of Systems.
// It is the recommended way to run ecs
type World struct {
	systems systems
}

// AddSystem adds a new System to the World, and then sorts them based on Priority
func (w *World) AddSystem(system System) {
	if initializer, ok := system.(Initializer); ok {
		initializer.New(w)
	}

	w.systems = append(w.systems, system)
	sort.Sort(w.systems)
}

// Systems returns a list of Systems
func (w *World) Systems() []System {
	return w.systems
}

// Update is called on each frame, with `dt` being the time difference in seconds since the last `Update` call
func (w *World) Update(dt float32) {
	for _, system := range w.Systems() {
		system.Update(dt)
	}
}

// RemoveEntity removes the entity across all systems
func (w *World) RemoveEntity(e BasicEntity) {
	for _, sys := range w.systems {
		sys.Remove(e)
	}
}
