package ecs

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MySystemOneEntity struct {
	e  *BasicEntity
	c1 *MyComponent1
}

type MySystemOneable interface {
	BasicFace
	MyComponent1Face
}

type MySystemOne struct {
	entities []MySystemOneEntity
}

func (*MySystemOne) Priority() int { return 0 }
func (*MySystemOne) New(*World)    {}
func (sys *MySystemOne) Update(dt float32) {
	for _, e := range sys.entities {
		e.c1.A++
	}
}
func (*MySystemOne) Remove(e BasicEntity) {}
func (sys *MySystemOne) Add(e *BasicEntity, c1 *MyComponent1) {
	sys.entities = append(sys.entities, MySystemOneEntity{e, c1})
}
func (sys *MySystemOne) AddByInterface(o Identifier) {
	obj := o.(MySystemOneable)
	sys.Add(obj.GetBasicEntity(), obj.GetMyComponent1())
}

type MySystemOneTwoEntity struct {
	e  *BasicEntity
	c1 *MyComponent1
	c2 *MyComponent2
}

type MySystemOneTwoable interface {
	BasicFace
	MyComponent1Face
	MyComponent2Face
}

type NotMySystemOneTwoable interface {
	NotMyComponent12Face
}

type MySystemOneTwo struct {
	entities []MySystemOneTwoEntity
}

func (*MySystemOneTwo) Priority() int { return 0 }
func (*MySystemOneTwo) New(*World)    {}
func (sys *MySystemOneTwo) Update(dt float32) {
	for _, e := range sys.entities {
		e.c1.B++
		e.c2.D++
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
func (sys *MySystemOneTwo) AddByInterface(o Identifier) {
	obj := o.(MySystemOneTwoable)
	sys.Add(obj.GetBasicEntity(), obj.GetMyComponent1(), obj.GetMyComponent2())
}

type MySystemTwoEntity struct {
	e  *BasicEntity
	c2 *MyComponent2
}

type MySystemTwoable interface {
	BasicFace
	MyComponent2Face
}

type NotMySystemTwoable interface {
	NotMyComponent2Face
}

type MySystemTwo struct {
	entities []MySystemTwoEntity
}

func (*MySystemTwo) Priority() int { return 0 }
func (*MySystemTwo) New(*World)    {}
func (sys *MySystemTwo) Update(dt float32) {
	for _, e := range sys.entities {
		e.c2.C++
	}
}
func (sys *MySystemTwo) Remove(e BasicEntity) {
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
func (sys *MySystemTwo) Add(e *BasicEntity, c2 *MyComponent2) {
	sys.entities = append(sys.entities, MySystemTwoEntity{e, c2})
}
func (sys *MySystemTwo) AddByInterface(o Identifier) {
	obj := o.(MySystemTwoable)
	sys.Add(obj.GetBasicEntity(), obj.GetMyComponent2())
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
type MyComponent1Face interface {
	GetMyComponent1() *MyComponent1
}

func (c *MyComponent1) GetMyComponent1() *MyComponent1 {
	return c
}

type MyComponent2 struct {
	C, D int
}
type MyComponent2Face interface {
	GetMyComponent2() *MyComponent2
}

func (c *MyComponent2) GetMyComponent2() *MyComponent2 {
	return c
}

type NotMyComponent2 struct{}
type NotMyComponent2Face interface {
	GetNotMyComponent2() *NotMyComponent2
}

func (n *NotMyComponent2) GetNotMyComponent2() *NotMyComponent2 {
	return n
}

type NotMyComponent12 struct{}
type NotMyComponent12Face interface {
	GetNotMyComponent12() *NotMyComponent12
}

func (n *NotMyComponent12) GetNotMyComponent12() *NotMyComponent12 {
	return n
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

// TestIdentifierInterface makes sure that my entity can be stored as an Identifier interface
func TestIdentifierInterface(t *testing.T) {
	e1 := MyEntity1{}
	e1.BasicEntity = NewBasic()

	var slice []Identifier = []Identifier{e1}

	_, ok := slice[0].(MyEntity1)
	assert.True(t, ok, "MyEntity1 should have been recoverable from the Identifier interface")
}

func TestSortableIdentifierSlice(t *testing.T) {
	e1 := MyEntity1{}
	e1.BasicEntity = NewBasic()
	e2 := MyEntity1{}
	e2.BasicEntity = NewBasic()

	var entities IdentifierSlice = []Identifier{e2, e1}
	sort.Sort(entities)
	assert.ObjectsAreEqual(e1, entities[0])
	assert.ObjectsAreEqual(e2, entities[1])
}

// TestSystemEntityFiltering checks that entities go into the right systems and the flags are obeyed
func TestSystemEntityFiltering(t *testing.T) {
	w := &World{}

	var sys1in *MySystemOneable
	w.AddSystemInterface(&MySystemOne{}, sys1in, nil)

	var sys2in *MySystemTwoable
	var sys2out *NotMySystemTwoable
	w.AddSystemInterface(&MySystemTwo{}, sys2in, sys2out)

	var sys12in *MySystemOneTwoable
	var sys12out *NotMySystemOneTwoable
	w.AddSystemInterface(&MySystemOneTwo{}, sys12in, sys12out)

	e1 := struct {
		BasicEntity
		*MyComponent1
	}{
		NewBasic(),
		&MyComponent1{},
	}
	w.AddEntity(&e1)

	e2 := struct {
		BasicEntity
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent2{},
	}
	w.AddEntity(&e2)

	e12 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
	}
	w.AddEntity(&e12)

	e12x2 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotMyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotMyComponent2{},
	}
	w.AddEntity(&e12x2)

	e12x12 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotMyComponent12
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotMyComponent12{},
	}
	w.AddEntity(&e12x12)

	e12x12x2 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotMyComponent12
		*NotMyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotMyComponent12{},
		&NotMyComponent2{},
	}
	w.AddEntity(&e12x12x2)

	w.Update(0.125)

	assert.Equal(t, 1, e1.A, "e1 was not updated by system 1")
	assert.Equal(t, 0, e1.B, "e1 was updated by system 12")

	assert.Equal(t, 1, e2.C, "e2 was not updated by system 2")
	assert.Equal(t, 0, e2.D, "e2 was updated by system 12")

	assert.Equal(t, 1, e12.A, "e12 was not updated by system 1")
	assert.Equal(t, 1, e12.B, "e12 was not updated by system 12")
	assert.Equal(t, 1, e12.C, "e12 was not updated by system 2")
	assert.Equal(t, 1, e12.D, "e12 was not updated by system 12")

	assert.Equal(t, 1, e12x2.A, "e12x2 was not updated by system 1")
	assert.Equal(t, 1, e12x2.B, "e12x2 was not updated by system 12")
	assert.Equal(t, 0, e12x2.C, "e12x2 was updated by system 2")
	assert.Equal(t, 1, e12x2.D, "e12x2 was not updated by system 12")

	assert.Equal(t, 1, e12x12.A, "e12x12 was not updated by system 1")
	assert.Equal(t, 0, e12x12.B, "e12x12 was updated by system 12")
	assert.Equal(t, 1, e12x12.C, "e12x12 was not updated by system 2")
	assert.Equal(t, 0, e12x12.D, "e12x12 was updated by system 12")

	assert.Equal(t, 1, e12x12x2.A, "e12x12x2 was not updated by system 1")
	assert.Equal(t, 0, e12x12x2.B, "e12x12x2 was updated by system 12")
	assert.Equal(t, 0, e12x12x2.C, "e12x12x2 was updated by system 2")
	assert.Equal(t, 0, e12x12x2.D, "e12x12x2 was updated by system 12")
}

// TestParentChild tests parenting of BasicEntities
func TestParentChild(t *testing.T) {
	parent := NewBasic()
	children := NewBasics(3)
	if len(parent.Children()) != 0 {
		t.Errorf("Children did not initalize to a zero value")
	}
	if parent.Parent() != nil {
		t.Errorf("Parent did not initalize as nil")
	}
	parent.AppendChild(&children[0])
	parent.AppendChild(&children[1])
	parent.AppendChild(&children[2])
	if len(parent.Children()) != 3 {
		t.Errorf("Failed to add all three children to parent")
	}
	for i := 0; i < 3; i++ {
		if children[i].Parent() != &parent {
			t.Errorf("Parent was not updated properly for children.")
		}
	}
}

// TestRemoveChild tests removing a child
func TestRemoveChild(t *testing.T) {
	parent := NewBasic()
	children := NewBasics(3)
	parent.AppendChild(&children[0])
	parent.AppendChild(&children[1])
	parent.AppendChild(&children[2])
	if len(parent.Children()) != 3 {
		t.Errorf("Parent did not successfully add all three children.")
	}
	parent.RemoveChild(&children[1])
	if len(parent.Children()) != 2 {
		t.Errorf("Parent didd not successfully remove a child")
	}
	for i := 0; i < 2; i++ {
		if children[1].ID() == parent.Children()[i].ID() {
			t.Errorf("Found removed child in parent still")
		}
	}
}

func TestDescendents(t *testing.T) {
	parent := NewBasic()
	children := NewBasics(7)
	parent.AppendChild(&children[0])
	parent.AppendChild(&children[1])
	parent.AppendChild(&children[2])
	children[0].AppendChild(&children[3])
	children[0].AppendChild(&children[4])
	children[1].AppendChild(&children[5])
	children[3].AppendChild(&children[6])

	if len(parent.Descendents()) != 7 {
		t.Errorf("Parent did not have all descendents.")
	}
	testMap := map[uint64]struct{}{}
	testMap[children[0].ID()] = struct{}{}
	testMap[children[1].ID()] = struct{}{}
	testMap[children[2].ID()] = struct{}{}
	testMap[children[3].ID()] = struct{}{}
	testMap[children[4].ID()] = struct{}{}
	testMap[children[5].ID()] = struct{}{}
	testMap[children[6].ID()] = struct{}{}
	for _, d := range parent.Descendents() {
		delete(testMap, d.ID())
	}
	if len(testMap) != 0 {
		t.Errorf("Expected children not found in parent. Did not find %v", testMap)
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

func BenchmarkNewBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasic()
	}
}

func BenchmarkNewBasics1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(1)
	}
}

func BenchmarkNewBasic10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(10)
	}
}

func BenchmarkNewBasic100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(100)
	}
}

func BenchmarkNewBasic1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(1000)
	}
}
