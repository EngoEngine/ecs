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

func (s *simpleSystem) Remove(b BasicEntity) {
}

func (s *simpleSystem) Update(dt float32) {
}
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
