package ecs

import (
	"reflect"
)

// Component is a piece of data which belongs to an Entity
type Component interface {
	Type() string
}

// Entity is the E in Entity Component System. It belongs to any amount of
// Systems, and has a number of Components
type Entity struct {
	id         string
	components map[string]Component
	requires   map[string]bool
	Pattern    string
}

// NewEntity creates a new Entity given an array of Systems which should be
// required
func NewEntity(requires ...string) *Entity {
	e := &Entity{
		id:         generateUUID(),
		requires:   make(map[string]bool),
		components: make(map[string]Component),
	}
	for _, req := range requires {
		e.requires[req] = true
	}
	return e
}

// DoesRequire checks if the Entity requires a system
func (e *Entity) DoesRequire(name string) bool {
	return e.requires[name]
}

// Requirements returns a list of all system requirements for this entity
func (e *Entity) Requirements() []string {
	keys := make([]string, len(e.requires))

	i := 0
	for k := range e.requires {
		keys[i] = k
		i++
	}
	return keys
}

// AddComponent adds a new Component to the Entity
func (e *Entity) AddComponent(component Component) {
	e.components[component.Type()] = component
}

// RemoveComponent removes a Component from the Entity
func (e *Entity) RemoveComponent(component Component) {
	delete(e.components, component.Type())
}

// Component takes a double pointer to a Component,
// and populates it with the value of the right type.
func (e *Entity) Component(x interface{}) bool {
	v := reflect.ValueOf(x).Elem() // *T
	c, ok := e.components[v.Interface().(Component).Type()]
	if !ok {
		return false
	}
	v.Set(reflect.ValueOf(c))
	return true
}

// ComponentFast returns the same object as Component
// but without using reflect (and thus faster).
// Be sure to define the .Type() such that it takes a pointer receiver
func (e *Entity) ComponentFast(c Component) interface{} {
	return e.components[c.Type()]
}

// ID returns the string ID of the Entity
func (e *Entity) ID() string {
	return e.id
}

// Duplicate populates an array with the specified amounts of clones from an Entity.
// NOTE: If you provide your components as pointers, the pointer will be duplicated.
// Hence, the entities will "share" a component.
func (e Entity) Duplicate(amount int) []Entity {
	var es []Entity

	for i := 0; i < amount; i++ {
		es = append(es, e.Clone())
	}

	return es
}

// Clone creates a duplication of an entity.
func (e Entity) Clone() Entity {
	clone := e
	clone.id = generateUUID()

	return clone
}
