package ecs

// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way suitable
// for the current implementation.
type System interface {
	// Update is ran every frame, with `dt` being the time in seconds since the last frame
	Update(dt float32)

	// Delete should remove the entity from the system completely
	Remove(e BasicEntity)
}

type Prioritizer interface {
	// Priority indicates the order in which Systems should be executed per iteration, higher meaning sooner. Default is 0
	Priority() int
}

type Initializer interface {
	// New is the initialisation of the System, and may be used to initialize some values beforehand, like storing
	// a reference to the World
	New(*World)
}

// Systems implements a sortable list of `System`. It is indexed on `System.Priority()`.
type Systems []System

func (s Systems) Len() int {
	return len(s)
}

func (s Systems) Less(i, j int) bool {
	var prio1, prio2 int

	if prior1, ok := s[i].(Prioritizer); ok {
		prio1 = prior1.Priority()
	}
	if prior2, ok := s[j].(Prioritizer); ok {
		prio2 = prior2.Priority()
	}

	return prio1 > prio2
}

func (s Systems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
