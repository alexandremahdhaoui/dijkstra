# alexandremahdhaoui/dijkstra

A production-grade, zero-dependency Go implementation of Dijkstra’s algorithm featuring generic node types, customizable starting/ending criteria, and dynamic path reconstruction.

---

## Overview

Go developers frequently need optimal pathfinding and reachability analysis across custom domain models (e.g., geospatial networks, system dependency graphs, routing tables). Standard Dijkstra implementations rely on fixed int/string identifiers or concrete node structs, forcing repetitive boilerplate, manual type conversion, and fixed point-to-point interface logic. But how can developers compute shortest paths on arbitrary Go data structures without sacrificing type safety, performance, or custom search conditions?

Leverage `alexandremahdhaoui/dijkstra`. This is a modern generic library supporting custom node logic, multi-source/multi-target filtering (`NodeSelector`), zero-weight edge support, and instant path extraction.

---

## Core Architectural Pillars

### 1. Type-Safe Genericity

Operates on any custom type implementing `Node[T, W]` and any numerical weight constraint (`Number`). No type assertions or manual ID mappings required.

```go
type City struct {
    Name  string
    Roads []dijkstra.Edge[*City, int]
}

func (c *City) Edges() []dijkstra.Edge[*City, int] {
    return c.Roads
}
```

### 2. Predicate-Driven Traversal

Uses declarative `NodeSelector` closures (`start`, `end`) to evaluate start roots, multi-source seeds, or dynamic property-based stop conditions on the fly.

```go
// Multi-Source Start
    isEUHub := func(c *City) bool {
    return c.IsHub && c.Region == "EU"
}

// Dynamic Target Match
    hasAirport := func(c *City) bool {
    return c.HasAirport
}

optimum := dijkstra.Search(rootNode, isEUHub, hasAirport)
```

### 3. Robust Path Reconstruction

Tracks state nodes, start origins (`IsStart`), and parent linkages via `OptimumState[T, W]` to guarantee precise backtrack traces—even across zero-cost paths or cyclic graph topologies.

```go
optimum := dijkstra.Search(startCity, dijkstra.By(startCity), dijkstra.By(targetCity))
path, cost, ok := optimum.PathTo(targetCity)
```

---

## Structural Anatomy & API Reference

### Core Search Engine (`Search`)

* **Input Criteria:** Takes a generic root node `root`, a flexible `start` selector function, and a dynamic `end` selector function.
* **Execution Strategy:** Employs an internal priority queue (`container/heap`) with visited state tracking and distance relaxation.

```go
// Run full network traversal by passing nil for the end selector
optimum := dijkstra.Search(startCity, dijkstra.By(startCity), nil)
```

### Output State (`OptimumState[T, W]`)

* **Path Recovery:** `PathTo(end T)` returns ordered node slices `[]T`, total accumulated cost `W`, and reachability status `bool`.
* **Complete Dispatch:** Retains the shortest path tree for all reached nodes, allowing sub-microsecond lookup queries without re-executing the search algorithm.

---

## Getting Started

### Installation

```bash
go get github.com/alexandremahdhaoui/dijkstra
```

### Verification & Examples

Run the test suite across all edge cases and numeric type validations:

```bash
go test -v ./...
```

Run the interactive showcase CLI:

```bash
go run ./cmd/example/main.go
```

---

## License

Apache License 2
