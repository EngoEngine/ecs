package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MySystemOneEntity struct {
	e  *BasicEntity
	c1 *MyComponent1
}

type MySystemOne struct {
	entities []MySystemOneEntity
}

func (*MySystemOne) Priority() int { return 0 }
func (*MySystemOne) New(*World)    {}
func (sys *MySystemOne) Update(dt float32) {
	for _, e := range sys.entities {
		e.c1.A = 5
	}
}
func (*MySystemOne) Remove(e BasicEntity) {}
func (sys *MySystemOne) Add(e *BasicEntity, c1 *MyComponent1) {
	sys.entities = append(sys.entities, MySystemOneEntity{e, c1})
}

type MySystemOneTwoEntity struct {
	e  *BasicEntity
	c1 *MyComponent1
	c2 *MyComponent2
}

type MySystemOneTwo struct {
	entities []MySystemOneTwoEntity
}

func (*MySystemOneTwo) Priority() int { return 0 }
func (*MySystemOneTwo) New(*World)    {}
func (sys *MySystemOneTwo) Update(dt float32) {
	for _, e := range sys.entities {
		if e.c1 == nil {
			return
		}

		if e.c2 == nil {
			return
		}
	}
}
func (sys *MySystemOneTwo) Remove(e BasicEntity) {
	delete := -1
	for index, entity := range sys.entities {
		if entity.e.ID() == e.ID() {
			delete = index
		}
	}
	if delete >= 0 {
		sys.entities = append(sys.entities[:delete], sys.entities[delete+1:]...)
	}
}
func (sys *MySystemOneTwo) Add(e *BasicEntity, c1 *MyComponent1, c2 *MyComponent2) {
	sys.entities = append(sys.entities, MySystemOneTwoEntity{e, c1, c2})
}

type MyEntity1 struct {
	BasicEntity
	MyComponent1
}

type MyEntity2 struct {
	BasicEntity
	MyComponent2
}

type MyEntity12 struct {
	BasicEntity
	MyComponent1
	MyComponent2
}

type MyComponent1 struct {
	A, B int
}

type MyComponent2 struct {
	C, D int
}

// TestCreateEntity ensures IDs which are created, are unique
func TestCreateEntity(t *testing.T) {
	e1 := MyEntity1{}
	e1.BasicEntity = NewBasic()

	e2 := MyEntity1{}
	e2.BasicEntity = NewBasic()

	assert.NotEqual(t, e1.id, e2.id, "BasicEntity IDs should be unique")
}

// TestChangeableComponents ensures that Components which are being referenced, are changeable
func TestChangeableComponents(t *testing.T) {
	w := &World{}

	sys1 := &MySystemOne{}
	w.AddSystem(sys1)

	e1 := MyEntity1{}
	e1.BasicEntity = NewBasic()

	sys1.Add(&e1.BasicEntity, &e1.MyComponent1)

	sys1.Update(0.125)

	assert.NotZero(t, e1.MyComponent1.A, "MySystemOne should have been able to change the value of MyComponent1.A")
}

// TestDelete tests a commonly used method for removing an entity from the list of entities
func TestDelete(t *testing.T) {
	const maxEntities = 10

	for j := 1; j < maxEntities; j++ {
		w := &World{}

		sys12 := &MySystemOneTwo{}
		w.AddSystem(sys12)

		var entities []BasicEntity

		// Add all of them
		for i := 0; i < maxEntities; i++ {
			e := MyEntity12{BasicEntity: NewBasic()}
			sys12.Add(&e.BasicEntity, &e.MyComponent1, &e.MyComponent2)
			entities = append(entities, e.BasicEntity) // in order to remove it without having a reference to e
		}

		before := len(sys12.entities)

		// Attempt to remove j
		sys12.Remove(entities[j])

		assert.Len(t, sys12.entities, before-1, "MySystemOne should now have exactly one less Entity")
	}
}

func BenchmarkIdiomatic(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		sys12 := &MySystemOneTwo{}
		w.AddSystem(sys12)

		e1 := MyEntity1{}
		e1.BasicEntity = NewBasic()

		sys12.Add(&e1.BasicEntity, &e1.MyComponent1, nil)
	}

	Bench(b, preload, setup)
}

func BenchmarkIdiomaticDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		sys12 := &MySystemOneTwo{}
		w.AddSystem(sys12)

		e12 := MyEntity12{}
		e12.BasicEntity = NewBasic()

		sys12.Add(&e12.BasicEntity, &e12.MyComponent1, &e12.MyComponent2)
	}

	Bench(b, preload, setup)
}

// Bench is a helper-function to easily benchmark one frame, given a preload / setup function
func Bench(b *testing.B, preload func(), setup func(w *World)) {
	w := &World{}

	preload()
	setup(w)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w.Update(1 / 120) // 120 fps
	}
}
