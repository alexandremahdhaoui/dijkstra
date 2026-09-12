package dijkstra_test

import (
	"reflect"
	"testing"

	"github.com/alexandremahdhaoui/dijkstra"
)

type CustomInt int
type CustomFloat float64

type TestNode[W dijkstra.Number] struct {
	Name  string
	Roads []dijkstra.Edge[*TestNode[W], W]
}

func (n *TestNode[W]) Edges() []dijkstra.Edge[*TestNode[W], W] {
	return n.Roads
}

func testNumericTypeHelper[W dijkstra.Number](t *testing.T, w1, w2, w3 W) {
	t.Helper()

	a := &TestNode[W]{Name: "A"}
	b := &TestNode[W]{Name: "B"}
	c := &TestNode[W]{Name: "C"}

	a.Roads = []dijkstra.Edge[*TestNode[W], W]{
		{To: b, Weight: w1},
		{To: c, Weight: w1 + w2 + w3},
	}
	b.Roads = []dijkstra.Edge[*TestNode[W], W]{
		{To: c, Weight: w2},
	}

	opt := dijkstra.Search(a, dijkstra.By(a), dijkstra.By(c))
	path, cost, ok := opt.PathTo(c)

	if !ok {
		t.Fatalf("expected path to exist")
	}
	if cost != w1+w2 {
		t.Errorf("expected cost %v, got %v", w1+w2, cost)
	}
	wantPath := []*TestNode[W]{a, b, c}
	if !reflect.DeepEqual(path, wantPath) {
		t.Errorf("expected path %v, got %v", wantPath, path)
	}
}

func TestSearch_AllNumericTypes(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[int](t, 10, 20, 5)
	})
	t.Run("int8", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[int8](t, 5, 10, 2)
	})
	t.Run("int16", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[int16](t, 100, 200, 50)
	})
	t.Run("int32", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[int32](t, 1000, 2000, 500)
	})
	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[int64](t, 10000, 20000, 5000)
	})
	t.Run("uint", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[uint](t, 10, 20, 5)
	})
	t.Run("uint8", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[uint8](t, 5, 10, 2)
	})
	t.Run("uint16", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[uint16](t, 100, 200, 50)
	})
	t.Run("uint32", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[uint32](t, 1000, 2000, 500)
	})
	t.Run("uint64", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[uint64](t, 10000, 20000, 5000)
	})
	t.Run("float32", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[float32](t, 1.5, 2.5, 0.5)
	})
	t.Run("float64", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[float64](t, 10.55, 20.45, 5.0)
	})
	t.Run("CustomInt", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[CustomInt](t, 10, 20, 5)
	})
	t.Run("CustomFloat", func(t *testing.T) {
		t.Parallel()
		testNumericTypeHelper[CustomFloat](t, 1.25, 3.75, 1.0)
	})
}

func TestSearch_BranchCoverage(t *testing.T) {
	t.Parallel()

	t.Run("NilStartDefaultsToRoot", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 5}}

		opt := dijkstra.Search(a, nil, dijkstra.By(b))
		path, cost, ok := opt.PathTo(b)

		if !ok || cost != 5 || len(path) != 2 {
			t.Fatalf("failed nil start fallback check: ok=%v, cost=%d, path=%v", ok, cost, path)
		}
	})

	t.Run("NilEndTraversesEntireGraph", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 1}}
		b.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: c, Weight: 2}}

		opt := dijkstra.Search(a, dijkstra.By(a), nil)

		if _, ok := opt[a]; !ok {
			t.Errorf("expected a in optimum state")
		}
		if _, ok := opt[b]; !ok {
			t.Errorf("expected b in optimum state")
		}
		if _, ok := opt[c]; !ok {
			t.Errorf("expected c in optimum state")
		}
	})

	t.Run("EarlyExitOnEndMatchAtRoot", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 10}}

		opt := dijkstra.Search(a, dijkstra.By(a), dijkstra.By(a))

		if stateA, ok := opt[a]; !ok || !stateA.Visited {
			t.Fatalf("root A should be visited")
		}
		if stateB, ok := opt[b]; ok && stateB.Visited {
			t.Fatalf("node B should not be visited due to early exit at root")
		}
	})

	t.Run("EarlyExitOnEndMatchIntermediate", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 1}, {To: c, Weight: 100}}
		b.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: c, Weight: 1}}

		opt := dijkstra.Search(a, dijkstra.By(a), dijkstra.By(b))

		if stateB, ok := opt[b]; !ok || !stateB.Visited {
			t.Fatalf("node B should be visited")
		}
		if stateC, ok := opt[c]; ok && stateC.Visited {
			t.Fatalf("node C should not be visited due to early exit at B")
		}
	})

	t.Run("LazyDeletionSkipStaleHeapItem", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}

		a.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: b, Weight: 10},
			{To: c, Weight: 1},
		}
		c.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: b, Weight: 1},
		}

		opt := dijkstra.Search(a, dijkstra.By(a), nil)
		path, cost, ok := opt.PathTo(b)

		if !ok || cost != 2 {
			t.Fatalf("expected cost 2 via C, got cost=%d, ok=%v", cost, ok)
		}
		wantPath := []*TestNode[int]{a, c, b}
		if !reflect.DeepEqual(path, wantPath) {
			t.Fatalf("expected path %v, got %v", wantPath, path)
		}
	})

	t.Run("SkipVisitedNeighbors", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}

		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 1}}
		b.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: a, Weight: 1}}

		opt := dijkstra.Search(a, nil, nil)
		if opt[a].Distance != 0 || opt[b].Distance != 1 {
			t.Fatalf("unexpected distances in cycle graph")
		}
	})

	t.Run("IgnoreEqualOrWorsePathToNeighbor", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}

		a.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: b, Weight: 2},
			{To: c, Weight: 1},
		}
		c.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: b, Weight: 5},
		}

		opt := dijkstra.Search(a, nil, nil)
		if opt[b].Distance != 2 || opt[b].ParentOf != a {
			t.Fatalf("worse path via C should have been ignored")
		}
	})

	t.Run("UnreachableNodeInGraph", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		unreachable := &TestNode[int]{Name: "Unreachable"}

		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 1}}

		opt := dijkstra.Search(a, dijkstra.By(a), dijkstra.By(unreachable))
		path, cost, ok := opt.PathTo(unreachable)

		if ok || path != nil || cost != 0 {
			t.Fatalf("unreachable node should return false, nil, 0; got ok=%v, path=%v, cost=%d", ok, path, cost)
		}
	})

	t.Run("MultiSourceStartSelector", func(t *testing.T) {
		t.Parallel()
		root := &TestNode[int]{Name: "Root"}
		s1 := &TestNode[int]{Name: "S1"}
		s2 := &TestNode[int]{Name: "S2"}
		target := &TestNode[int]{Name: "Target"}

		root.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: s1, Weight: 0}, {To: s2, Weight: 0}}
		s1.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: target, Weight: 20}}
		s2.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: target, Weight: 5}}

		multiStart := func(n *TestNode[int]) bool {
			return n == s1 || n == s2
		}

		opt := dijkstra.Search(root, multiStart, dijkstra.By(target))
		path, cost, ok := opt.PathTo(target)

		if !ok || cost != 5 {
			t.Fatalf("expected cost 5 from nearest source S2, got %d", cost)
		}
		wantPath := []*TestNode[int]{s2, target}
		if !reflect.DeepEqual(path, wantPath) {
			t.Fatalf("expected path %v, got %v", wantPath, path)
		}
	})
}

func TestGraphTopologies(t *testing.T) {
	t.Parallel()

	t.Run("SelfLoop", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: a, Weight: 1},
			{To: b, Weight: 3},
		}

		opt := dijkstra.Search(a, nil, nil)
		path, cost, ok := opt.PathTo(b)

		if !ok || cost != 3 || len(path) != 2 {
			t.Fatalf("self loop handling failed: ok=%v, cost=%d, path=%v", ok, cost, path)
		}
	})

	t.Run("ParallelEdges", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{
			{To: b, Weight: 10},
			{To: b, Weight: 2},
			{To: b, Weight: 5},
		}

		opt := dijkstra.Search(a, nil, nil)
		_, cost, ok := opt.PathTo(b)

		if !ok || cost != 2 {
			t.Fatalf("parallel edge resolution failed: expected cost 2, got %d", cost)
		}
	})

	t.Run("ZeroWeightEdges", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}
		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 0}}
		b.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: c, Weight: 0}}

		opt := dijkstra.Search(a, nil, nil)
		path, cost, ok := opt.PathTo(c)

		if !ok || cost != 0 || len(path) != 3 {
			t.Fatalf("zero weight path failed: ok=%v, cost=%d, path=%v", ok, cost, path)
		}
	})

	t.Run("DisconnectedSubgraphInBFS", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}
		c := &TestNode[int]{Name: "C"}
		d := &TestNode[int]{Name: "D"}

		a.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: b, Weight: 1}}
		c.Roads = []dijkstra.Edge[*TestNode[int], int]{{To: d, Weight: 1}}

		opt := dijkstra.Search(a, nil, nil)

		if _, ok := opt[a]; !ok {
			t.Errorf("expected A in state")
		}
		if _, ok := opt[b]; !ok {
			t.Errorf("expected B in state")
		}
		if _, ok := opt[c]; ok {
			t.Errorf("disconnected node C should not be in state")
		}
	})
}

func TestOptimumState_PathTo(t *testing.T) {
	t.Parallel()

	t.Run("MissingKeyReturnsZeroAndFalse", func(t *testing.T) {
		t.Parallel()
		opt := make(dijkstra.OptimumState[*TestNode[int], int])
		missing := &TestNode[int]{Name: "Missing"}

		path, cost, ok := opt.PathTo(missing)
		if ok || path != nil || cost != 0 {
			t.Fatalf("expected false, nil, 0; got ok=%v, path=%v, cost=%d", ok, path, cost)
		}
	})

	t.Run("SingleNodePathToSelf", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		opt := dijkstra.Search(a, nil, nil)

		path, cost, ok := opt.PathTo(a)
		if !ok || cost != 0 || len(path) != 1 || path[0] != a {
			t.Fatalf("path to self failed: ok=%v, cost=%d, path=%v", ok, cost, path)
		}
	})

	t.Run("NilParentInChainBreak", func(t *testing.T) {
		t.Parallel()
		a := &TestNode[int]{Name: "A"}
		b := &TestNode[int]{Name: "B"}

		opt := make(dijkstra.OptimumState[*TestNode[int], int])
		opt[b] = &dijkstra.State[*TestNode[int], int]{
			Node:     b,
			ParentOf: a,
			Distance: 5,
		}

		path, cost, ok := opt.PathTo(b)
		if !ok || cost != 5 {
			t.Fatalf("expected ok=true and cost=5, got ok=%v, cost=%d", ok, cost)
		}
		wantPath := []*TestNode[int]{a, b}
		if !reflect.DeepEqual(path, wantPath) {
			t.Fatalf("expected path %v, got %v", wantPath, path)
		}
	})
}

func TestSelectorHelper_By(t *testing.T) {
	t.Parallel()

	a := &TestNode[int]{Name: "A"}
	b := &TestNode[int]{Name: "B"}

	selector := dijkstra.By(a)

	if !selector(a) {
		t.Errorf("expected By(a) to return true for node a")
	}
	if selector(b) {
		t.Errorf("expected By(a) to return false for node b")
	}
}
