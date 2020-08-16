package ecs

import (
	"testing"
)

func TestWorld_AddSystemInterface(t *testing.T) {
	type foo interface {
		a() string
	}
	var fooInstance *foo
	type args struct {
		sys SystemAddByInterfacer
		in  interface{}
		ex  interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// {"adds individual interface", args{}},
		{"adds multiple interfaces", args{nil, fooInstance, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := new(World)
			w.AddSystemInterface(tt.args.sys, tt.args.in, tt.args.ex)
		})
	}
}

type simpleEntity struct {
	BasicEntity
}
type entent struct {
	*BasicEntity
}
type simpleSystem struct {
	entities []entent
}

func (s *simpleSystem) Add(b *BasicEntity) {
	s.entities = append(s.entities, entent{b})
}

func (s *simpleSystem) AddByInterface(i Identifier) {
	obj, ok := i.(BasicFace)
	if ok {
		s.Add(obj.GetBasicEntity())
	}
}

func (s *simpleSystem) Remove(b BasicEntity) {}

func (s *simpleSystem) Update(dt float32) {}

func TestWorld_AddEntity(t *testing.T) {

	type args struct {
		systems []SystemAddByInterfacer
		e       Identifier
	}
	tests := []struct {
		name string
		args args
	}{
		{"works with multiple interfaces", args{
			systems: []SystemAddByInterfacer{&simpleSystem{}},
			e:       &simpleEntity{NewBasic()},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := new(World)
			sys := new(simpleSystem)
			var face *BasicFace
			w.AddSystemInterface(sys, []interface{}{face}, nil)
			w.AddEntity(tt.args.e)
			if len(sys.entities) == 0 {
				t.Error(len(sys.entities))
			}
		})
	}
}

type priorityChangeSystem struct {
	Rank int
}

func (s *priorityChangeSystem) Priority() int {
	return s.Rank
}

func (s *priorityChangeSystem) Add(b *BasicEntity) {}

func (s *priorityChangeSystem) AddByInterface(i Identifier) {
	obj, ok := i.(BasicFace)
	if ok {
		s.Add(obj.GetBasicEntity())
	}
}

func (s *priorityChangeSystem) Remove(b BasicEntity) {}

func (s *priorityChangeSystem) Update(dt float32) {}

func TestWorld_SortSystems(t *testing.T) {
	w := new(World)
	one := priorityChangeSystem{Rank: 1}
	w.AddSystem(&one)
	two := priorityChangeSystem{Rank: 2}
	w.AddSystem(&two)
	expected := []System{
		&two,
		&one,
	}
	for idx, sys := range w.Systems() {
		p := sys.(Prioritizer)
		exp := expected[idx].(Prioritizer)
		if p.Priority() != exp.Priority() {
			t.Error("Systems were not in the correct order")
		}
	}
	one.Rank = 5
	for idx, sys := range w.Systems() {
		p := sys.(Prioritizer)
		exp := expected[idx].(Prioritizer)
		if p.Priority() != exp.Priority() {
			t.Error("Systems were switched before sort wwas called")
		}
	}
	w.SortSystems()
	for idx, sys := range w.Systems() {
		p := sys.(Prioritizer)
		exp := expected[len(expected)-1-idx].(Prioritizer)
		if p.Priority() != exp.Priority() {
			t.Error("Systems were switched before sort wwas called")
		}
	}
}
