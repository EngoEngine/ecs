package ecs

import "testing"

type PriorityComponent struct {
	throughFirstSystem bool
}

type PriorityEntity struct {
	BasicEntity
	*PriorityComponent
}

type SystemPriorityFirst struct {
	e PriorityEntity
}

func (s *SystemPriorityFirst) Priority() int { return 500 }

func (s *SystemPriorityFirst) Add(b BasicEntity, p *PriorityComponent) {
	s.e = PriorityEntity{b, p}
}

func (s *SystemPriorityFirst) Remove(basic BasicEntity) {}

func (s *SystemPriorityFirst) Update(dt float32) {
	s.e.throughFirstSystem = true
}

type SystemPrioritySecond struct {
	e    PriorityEntity
	pass bool
}

func (s *SystemPrioritySecond) Priority() int { return 1 }

func (s *SystemPrioritySecond) Add(b BasicEntity, p *PriorityComponent) {
	s.e = PriorityEntity{b, p}
}

func (s *SystemPrioritySecond) Remove(basic BasicEntity) {}

func (s *SystemPrioritySecond) Update(dt float32) {
	if s.e.throughFirstSystem {
		s.pass = true
	}
}

// TestSystemPriority tests if systems get added based on priority
func TestSystemPriority(t *testing.T) {
	w := &World{}
	sys2 := SystemPrioritySecond{}
	w.AddSystem(&sys2)
	w.AddSystem(&SystemPriorityFirst{})
	ent := struct {
		BasicEntity
		PriorityComponent
	}{
		BasicEntity: NewBasic(),
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemPriorityFirst:
			sys.Add(ent.BasicEntity, &ent.PriorityComponent)
		case *SystemPrioritySecond:
			sys.Add(ent.BasicEntity, &ent.PriorityComponent)
		}
	}
	w.Update(1)
	if !sys2.pass {
		t.Error("Systems were not run in order after updated by the world.")
	}
}

type SystemAddRemove struct {
	entities []PriorityEntity
}

func (s *SystemAddRemove) Add(b BasicEntity, p *PriorityComponent) {
	s.entities = append(s.entities, PriorityEntity{b, p})
}

func (s *SystemAddRemove) Remove(basic BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *SystemAddRemove) Update(dt float32) {}

// TestAddRemove tests adding and removing entities in systems that don't implement
// SystemAddByInterfacer to the world.
func TestAddRemove(t *testing.T) {
	w := &World{}
	sys := SystemAddRemove{}
	w.AddSystem(&sys)
	ent := struct {
		BasicEntity
		PriorityComponent
	}{
		BasicEntity: NewBasic(),
	}
	w.AddEntity(ent)
	if len(sys.entities) != 0 {
		t.Error("Entity was added even though the system does not implement SystemAddByInterfacer")
	}
	sys.Add(ent.BasicEntity, &ent.PriorityComponent)
	if len(sys.entities) != 1 {
		t.Error("Failed to add entity to system")
	}
	w.RemoveEntity(ent.BasicEntity)
	if len(sys.entities) != 0 {
		t.Error("Removing the entity from the world did not remove it from the system.")
	}
}
