# ecs

[![Build Status](https://github.com/EngoEngine/ecs/workflows/CI/badge.svg)](https://github.com/EngoEngine/ecs/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/EngoEngine/ecs)](https://goreportcard.com/report/github.com/EngoEngine/ecs)

This is our implementation of the "Entity Component System" model in Go. It was designed to be used in `engo`, however
it is not dependent on any other packages so is able to be used wherever!

## Basics
In the Entity Component System paradigm, you have three elements;

* Entities
* Components
* Systems.

In our implementation, we use the type `World` to work with those `System`s. Each `System` can have references to any number (including 0) of entities. And each `Entity` can have as many `Component`s as desired.

An example of creating a `World`, adding a `System` to it, and update all systems
```go
// Declare the world - you can also use "var world ecs.World"
world := ecs.World{}

// You can add as many Systems here as you like. The RenderSystem provided by `engo` is just an example.
world.AddSystem(&engo.RenderSystem{})

// This will usually be called within the game-loop, in order to update all Systems on every frame.  
world.Update(0.125) // 0.125 would be the time in seconds since the last update
```

## System
We've been talking about `System`s, but what are they? Anything that implements the interface, can be used as a `System`:

```go
type System interface {
	// Update is ran every frame, with `dt` being the time in seconds since the last frame
	Update(dt float32)

	// Delete should remove the entity from the system completely
	Remove(e BasicEntity)
}
```

What does this say? It needs to have an `Update` method (which is called from `world.Update`), and it needs to have a `Remove(ecs.BasicEntity)` method. Why require a Remove method, but not an Add method? Because there's no 'generic' `Add` method (the parameters may change), while in order to remove something, all you need it the unique identifier (as provided by the `BasicEntity`).

### Initialization
Optionally, your `System` may implement the `Initializer` interface, which allows you to do initialization for the given `World`. Basically, it allows you to initialize values, without having to call the function manually before adding it to the `World`. Whenever you add a `System` (one that implements the `Initializer` interface) to the world, the `New` method will be called.

```go
type Initializer interface {
	// New is the initialisation of the System, and may be used to initialize some values beforehand, like storing
	// a reference to the World
	New(*World)
}
```

### Priority
Optionally, your `System` may implement the `Prioritizer` interface, which allows the `World` to sort the `System`s based on that priority. If omitted, a value of `0` is assumed.

```go
type Prioritizer interface {
	// Priority indicates the order in which Systems should be executed per iteration, higher meaning sooner. Default is 0
	Priority() int
}
```

## Entities and Components
Where do the entities come in? All game-logic has to be done within `System`s (the `Update` method, to be precise)). `Component`s store data (which is used by those `System`s). An `Entity` is no more than a wrapper which combines multiple `Component`s and adds a unique identifier to the whole. This unique identifier is nothing magic: simply an incrementing integer value - nothing to worry about.

> Because the precise definition of those `Component`s can vary, this `ecs` package provides no `Component`s -- we only provide examples here. The `github.com/EngoEngine/engo/common` package offers lots of `Component`s and `System`s to work with, out of the box.

Let's view an example:

```go
type SpaceComponent struct {
    Width  float32
    Height float32
}

type HealthComponent struct {
    HealthPercentage float32
    ManaPercentage   float32
}

type Player struct {
    ecs.BasicEntity
    SpaceComponent
    HealthComponent
}
```

Here, the type `Player` is made out of three elements: the unique identifier (`ecs.BasicEntity`) and two `Component`s. A `System` may make use of one or more of those `Component`s. Which are required, is defined by the `Add` method on that `System`.

Let's view a few examples:

```go
func (MySystem1) Add(basic *ecs.BasicEntity, space *SpaceComponent) { /* ... */ }

func (MySystem2) Add(basic *ecs.BasicEntity, health *HealthComponent) { /* ... */ }

func (MySystem3) Add(basic *ecs.BasicEntity, space *SpaceComponent, health *HealthComponent) { /* ... */ }
```

These three different `Add` methods are all valid, and use different Components. But how can I add my `Entity` to the `System`, if I didn't save a reference to that `System`?

```go
// Initialize our custom Entity
// NOTE: we have to call `ecs.NewBasic` here, to give our Entity a new unique identifier
player := Player{BasicEntity: ecs.NewBasic()}

// Loop over all Systems
for _, system := range world.Systems() {

    // Use a type-switch to figure out which System is which
    switch sys := system.(type) {

        // Create a case for each System you want to use
        case *MySystem1:
            sys.Add(&player.BasicEntity, &player.SpaceComponent)
        case *MySystem3:
            sys.Add(&player.BasicEntity, &player.SpaceComponent, &player.Healthcomponent)
    }
}
```

That is all there is to it.

## Custom Systems - How to save Entities?

You more than likely will want to create `System`s yourself. We will now go in depth on what you should do when defining your own `Add` method for your `System`. As seen above, you can create any number (and type of) parameters you want.

> We do ask you to let *the first argument* be of type `*ecs.BasicEntity` - as a general rule.

Your `System` should include an array, slice or map in which to store those entities. Now it is important to note that you're not receiving entities per se -- you are receiving references to the `Component`s you need. The actual `Entity` (type `Player` in our example) may contain way more `Component`s. You will most-likely want to create a struct for you to store those pointers in. An example:

```go
type myAwesomeEntity struct {
    *ecs.BasicEntity
    *SpaceComponent
}

type MyAwesomeSystem struct {
    entities []myAwesomeEntity
}

func (m *MyAwesomeSystem) Add(basic *ecs.BasicEntity, space *SpaceComponent) {
    m.entities = append(m.entities, myAwesomeEntity{basic, space})
}
```

> ### NOTE
> As a convention, please include "System" in the name of your `System` -- at the end. When you define a struct (which contains pointers, as opposed to the `Player` struct we created earlier), please replace that `System` part with `Entity`. You should **only** use this newly-created struct in your similarly-named `System`. You will usually *never* want to export that `Entity` definition, as it is only being used in that `System`. If your system would be called `BallMovementSystem`, then your struct would be called `ballMovementEntity`.

### Removing Entities from your System
Your `System` must implement the `Remove` method as specified by the `System` interface. Whenever you start storing entities, you should define this method in such a way, that it removes the custom-created non-exported `Entity`-struct from the array, slice or map. An `ecs.BasicEntity` is given for you to figure out which element in the array, slice or map it is.

```go
// Remove removes the Entity from the System. This is what most Remove methods will look like
func (m *MyAwesomeSystem) Remove(basic ecs.BasicEntity) {
  	var delete int = -1
  	for index, entity := range m.entities {
    		if entity.ID() == basic.ID() {
    			delete = index
    			break
  		  }
  	}
  	if delete >= 0 {
    		m.entities = append(m.entities[:delete], m.entities[delete+1:]...)
  	}
}

// OR, if you were using a `map` instead of a `slice`:

// Remove removes the Entity from the System. As you see, removing becomes easier when using a `map`.
func (m *MyAwesomeSystem) Remove(basic ecs.BasicEntity) {
  	delete(m.entities, basic.ID())
}
//
```

> #### NOTE
> Even though that a `map` looks easier, if you want to loop over that `map` each frame, writing those additional lines to use a `slice` instead, is definitely worth it in terms of runtime performance. Iterating over a `map` is a lot slower.

## Custom Systems - The Update method
Whatever your `System` does on the `Update` method, is up to you. Each `System` is unique in that sense. If you're storing entities, then you might want to loop over them each frame. Again, this depends on your use-case.

```go
func (m *MyAwesomeSystem) Update(dt float32) {
    for _, entity := range m.entities {
        fmt.Println("I would like to tell you", entity.ID(), "that it has been", dt, "seconds since the last time we spoke. ")
    }
}
```

# Automatically add entities to systems
When your game gets *really* big, adding each entity to every system would be time consuming and buggy using the methods mentioned above. However, you can easily add entities to systems based solely on the interfaces that entity implements by
utilizing the `SystemAddByInterfacer`. This takes a bit of work up front, but makes things much easier if your number of systems and entities increases. We're going to start with an example `System` MySystem, with `Component` ComponentA

```go
type ComponentA struct {
    num int
}

type mySystemEntity struct {
    ecs.BasicEntity
    *ComponentA
}

type MySystem struct {
    entities []mySystemEntity
}

type (m *MySystem) Add(basic ecs.BasicEntity, a *ComponentA) { /* Add stuff goes here */ }
type (m *MySystem) Remove(basic ecs.BasicEntity) { /* Remove stuff here */ }
type (m *MySystem) Update(dt float32) { /* Update stuff here */ }
```

The components need to have corresponding Getters and Interfaces in order to be utilized. Let's add them

```go
func (a *ComponentA) GetComponentA() *ComponentA {
    reurn a
}

type AFace interface {
    GetComponentA() *ComponentA
}
```

### Note
The convention is that we add Face to the end of the component's name for the interface.

Now that we have interfaces for all the components, we need to add an interface to tell if we use the system or not. (BasicEntity already has this setup for you, as does any component or system that uses entities in `engo/common`)

```go
type Myable interface {
    ecs.BasicFace
    AFace
}
```

### Note
The convention is to add able to the end of the system's name for the interface

Finally, we have to add the AddByInterface function to the system. Don't worry about the casting, it can't panic as the world makes sure it implements the required interface befor passing entities to it.

```go
func (m *MySystem) AddByInterface(o ecs.Identifier) {
    obj := o.(Myable)
    m.Add(obj.GetBasicEntity(), obj.GetComponentA())
}
```

To use the system, instead of `w.AddSystem()` use

```go
var myable *Myable
w.AddSystemInterface(&MySystem{}, myable, nil)
```

### Note
This takes **a pointer to** the interface that the system needs implemented to use AddByInterface.

Finally, to add an entity, rather than looping through all the systems, you can just

```go
w.AddEntity(&entity)
```

## Exclude flags
You can also add an interface to the system for components that can act as flags to NOT add an entity to that system. First you'll have to make the component. It'll have to have a Getter and Interface as well.

```go
type NotMyComponent struct {}
type NotMyFace interface {
    GetNotMyComponent() *NotMyComponent
}
func (n *NotMyComponent) GetNotMyComponent() *NotMyComponent {
    return n
}
```

Then you can make the interface for the system

```go
type NotMyable interface {
    NotMyFace
}
```

Finally, we add it to the world

```go
var myable *Myable
var notMyable *NotMyable
w.AddSystemInterface(&MySystem{}, myable, notMyable)
```

Now our system can automatically, and it'll include all the entities that implement the Myable interface, except any entity that implements the NotMyable interface.
