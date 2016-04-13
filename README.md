# ECS

[![Build Status](https://travis-ci.org/EngoEngine/ecs.svg?branch=master)](https://travis-ci.org/EngoEngine/ecs)

## What's ECS?
ECS stands for Entity-Component-System paradigm. More information can be found [here](https://en.wikipedia.org/wiki/Entity_component_system). 

Basically, you have a `World`, which consists of a series of different `System`s and `Entity`s. The `Entity`s are 
modular, and can store values within `Component`s. All the logic happens within `System`s: so whenever you need to
add / alter / update an `Entity`, you will usually do so from within a `System`. If you need to pass information
from one `System` to another, you can store intermediary values within a `Component` within that `Entity`, and read
that value back from within the other `System`. 

## So how does it work?
This currently forms the core of the game-engine [engo](https://github.com/EngoEngine/engo). 

If you want to use it without using `engo`, you can continue reading. 


### Globally 

```go
// Defining it, and calling the required initialization function
world := ecs.World{}
world.New()

// Then you add a series of Systems - you can do this at any point
world.AddSystem(mySystem1)
world.AddSystem(mySystem2)

// After that, you can add a few Entities to the World. The first is being added to mySystem1, 
// the second to mySystem2 as well
entity1 := ecs.NewEntity([]string{"mySystem1"})
world.AddEntity(entity1)

entity2 := ecs.NewEntity([]string{"mySystem1", "mySystem2"})
world.AddEntity(entity2)
```

### Systems
Anything that implements the `System` interface can be used within ECS:
```go
// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way suitable
// for the current implementation.
type System interface {
	// Type returns a unique string identifier, usually the struct name
	// eg. "RenderSystem", "CollisionSystem"...
	Type() string
	// Priority is used to create the order in which Systems (in the World) are processed
	Priority() int

	// New is the initialisation of the System
	New(*World)
	// Update is ran every frame, with `dt` being the time in seconds since the last frame
	Update(dt float32)

	// AddEntity adds a new Entity to the System
	AddEntity(entity *Entity)
	// RemoveEntity removes an Entity from the System
	RemoveEntity(entity *Entity)
}
```

So we have a `Type() string` function, which should uniquely identify the `System` type, to distinguish between
the different `System`s. The output of this is also used when adding an `Entity` to the `World`, to see on which
`System`s the `Entity`s depend. 

The `Priority() int` is being used to compute the order in which to process all `System`s in a world-loop. 

The `New(*World)` function is being called to initialize the `System` for a given world. This makes it easy to 
reinitialize it whenever a World changes. This often happens in `engo` whenever the Scene changes (whenever you open
a menu for example). 

The key function here is `Update(dt float32)`. Each time the world loops, one call to the `System` is made through
this `Update` function. The `dt` parameter is the time in seconds since the last loop. This value is thus the same
for all `System`s, but changes each world-loop. Here you can do most of your logic: path finding, movement, animation, 
etc. etc.

Since you might want to know which `Entity`s of the entire `World` your `System` needs to work with (as defined by
the `Entity`s), there are the functions `AddEntity` and `RemoveEntity`. These allow you to optionally keep track of
which `Entity`s depend on your `System`, and which you should be updating. Example: your `System` computes the path
minions should walk in your game. It would make sense it only does that for minions, and not for static things like
rocks or other scenery. By stating that those `MinionEntity`s depend on your `PathSystem`, you could use these
`AddEntity` and `RemoveEntity` functions to keep track of just those `Entity`s you need to worry about. 

### Entities
An `Entity` is just a glorified `[]Component`, and it's easier to look at it like that. 

### Components
At any time, we can add `Component`s to those `Entity`s:

```go
entity1 := ecs.NewEntity([]string{"mySystem1"})
entity1.AddComponent(&myComponent1{})
world.AddEntity(entity1)

entity2 := ecs.NewEntity([]string{"mySystem2"})
world.AddEntity(entity2)
entity2.AddComponent(&myComponent2{})
```

Note that you can only add one `Component` per `Entity`.  
