package ecs

import (
	"fmt"
	"testing"
)

const (
	benchmarkComponentCount = 1000
)

type MyComponent1 struct{ an int }

func (*MyComponent1) Type() string { return "MyComponent1" }

type MyComponent2 struct{ an int }

func (*MyComponent2) Type() string { return "MyComponent2" }

type getComponentSystem struct {
	LinearSystem
}

func (getComponentSystem) Type() string {
	return "getComponentSystem"
}

func (*getComponentSystem) New(*World) {}
func (*getComponentSystem) Pre()       {}
func (*getComponentSystem) Post()      {}

func (g *getComponentSystem) UpdateEntity(entity *Entity, dt float32) {
	var sp *MyComponent1
	if !entity.Component(&sp) {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if sp == nil {
		return
	}

	if len(entity.components) != 2 {
		return
	}

	var ren *MyComponent2
	if !entity.Component(&ren) {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if ren == nil {
		return
	}
}

func TestDuplication(t *testing.T) {
	e := NewEntity([]string{})
	amount := 20

	es := e.Duplicate(amount)

	if amount != len(es) {
		t.Errorf("Not same amount of entities %v %v", amount, len(es))
	}
}

func TestClone(t *testing.T) {
	e := NewEntity([]string{})
	e2 := e.Clone()

	if e.ID() == e2.ID() {
		t.Error("IDs are the same :/ (a statistical anomaly may have occurred.)")
	}
}

func BenchmarkComponent(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystem{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity("getComponentSystem")
			e.AddComponent(&MyComponent1{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

func BenchmarkComponentDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystem{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity("getComponentSystem")
			e.AddComponent(&MyComponent1{})
			e.AddComponent(&MyComponent2{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

type getComponentSystemFast struct {
	LinearSystem
}

func (getComponentSystemFast) Type() string {
	return "getComponentSystemFast"
}

func (*getComponentSystemFast) New(*World) {}
func (*getComponentSystemFast) Pre()       {}
func (*getComponentSystemFast) Post()      {}

func (g *getComponentSystemFast) UpdateEntity(entity *Entity, dt float32) {
	var sp *MyComponent1
	var ok bool
	if sp, ok = entity.ComponentFast(sp).(*MyComponent1); !ok {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if sp == nil {
		return
	}

	if len(entity.components) != 2 {
		return
	}

	var ren *MyComponent2
	if ren, ok = entity.ComponentFast(ren).(*MyComponent2); !ok {
		return
	}
	// Not needed, but we need to ensure it gets compiled correctly
	if ren == nil {
		return
	}
}

func BenchmarkComponentFast(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystemFast{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity("getComponentSystemFast")
			e.AddComponent(&MyComponent1{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

func BenchmarkComponentFastDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&getComponentSystemFast{})
		for i := 0; i < benchmarkComponentCount; i++ {
			e := NewEntity("getComponentSystemFast")
			e.AddComponent(&MyComponent1{})
			e.AddComponent(&MyComponent2{})
			w.AddEntity(e)
		}
	}
	Bench(b, preload, setup)
}

func BenchmarkComponentPure(b *testing.B) {
	e := NewEntity()
	e.AddComponent(&MyComponent1{1})

	b.ResetTimer()
	var comp1 *MyComponent1
	var ok bool

	for i := 0; i < b.N; i++ {
		ok = e.Component(&comp1)
	}

	fmt.Sprint(comp1, ok)
}

func BenchmarkComponentFastPure(b *testing.B) {
	e := NewEntity()
	e.AddComponent(&MyComponent1{1})

	b.ResetTimer()
	var comp1 *MyComponent1
	var ok bool

	for i := 0; i < b.N; i++ {
		comp1, ok = e.ComponentFast(comp1).(*MyComponent1)
	}

	fmt.Sprint(comp1, ok)
}
