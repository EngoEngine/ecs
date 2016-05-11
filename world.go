package ecs

import (
	"sort"
)

// World contains a bunch of Entities, and a bunch of Systems. It is the
// recommended way to run ecs.
type World struct {
	systems systems
}

// AddSystem adds the given System to the World, sorted by priority.
func (w *World) AddSystem(system System) {
	if initializer, ok := system.(Initializer); ok {
		initializer.New(w)
	}

	w.systems = append(w.systems, system)
	sort.Sort(w.systems)
}

// Systems returns the list of Systems managed by the World.
func (w *World) Systems() []System {
	return w.systems
}

// Update updates each System managed by the World. It is invoked by the engine
// once every frame, with dt being the duration since the previous update.
func (w *World) Update(dt float32) {
	for _, system := range w.Systems() {
		system.Update(dt)
	}
}

// RemoveEntity removes the entity across all systems.
func (w *World) RemoveEntity(e BasicEntity) {
	for _, sys := range w.systems {
		sys.Remove(e)
	}
}
