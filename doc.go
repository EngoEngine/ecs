// Package ecs provides interfaces for the Entity Component System (ECS)
// paradigm used by engo.io/engo. It is predominately used by games, however
// will find use in other applications.
//
// The ECS paradigm aims to decouple distinct domains (e.g. rendering, input
// handling, AI) from one another, through a composition of independent
// components. The core concepts of ECS are described below.
//
//
// Entities
//
// An entity is simply a set of components with a unique ID attached to it,
// nothing more. In particular, an entity has no logic attached to it and stores
// no data explicitly (except for the ID).
//
// Each entity corresponds to a specific entity within the game, such as a
// character, an item, or a spell.
//
//
// Components
//
// A component stores the raw data related to a specific aspect of an entity,
// nothing more. In particular, a component has no logic attached to it.
//
// Different aspects may include the position, animation graphics, or input
// actions of an entity.
//
//
// Systems
//
// A system implements logic for processing entities possessing components of
// the same aspects as the system.
//
// For instance, an animation system may render entities possessing animation
// components.
package ecs
