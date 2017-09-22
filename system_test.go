package ecs

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PrioritizedSystem struct {
	entities []*BasicEntity
	priority int
}

func (p *PrioritizedSystem) Priority() int { return p.priority }
func (*PrioritizedSystem) New(*World)      {}

func (sys *PrioritizedSystem) Update(float32) {
}
func (*PrioritizedSystem) Remove(BasicEntity) {}

// TestPrioritizedSystemsSort according to priorty (highest first)
func TestPrioritizedSystemsSort(t *testing.T) {
	assert := assert.New(t)

	var lowPriority, highPriority *PrioritizedSystem
	lowPriority = &PrioritizedSystem{nil, 0}
	highPriority = &PrioritizedSystem{nil, 1}
	var systems systems = []System{lowPriority, highPriority}
	sort.Sort(systems)
	assert.Equal(systems[0].(Prioritizer).Priority(), 1)
	assert.Equal(systems[1].(Prioritizer).Priority(), 0)
}

// TestUnPrioritizedSystemsStable when sorted
func TestUnPrioritizedSystemsStable(t *testing.T) {
	var first, second *PrioritizedSystem
	first = &PrioritizedSystem{nil, 0}
	second = &PrioritizedSystem{nil, 0}
	var systems systems = []System{first, second}
	sort.Sort(systems)
	sortedFirst := systems[0].(*PrioritizedSystem)
	sortedSecond := systems[1].(*PrioritizedSystem)
	// assert.Equal doesn't do reference comparison
	if first != sortedFirst || second != sortedSecond {
		t.Errorf("Sort of systems is not stable")
	}
}
