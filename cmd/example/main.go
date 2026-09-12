package main

import (
	"fmt"
	"strings"

	"github.com/alexandremahdhaoui/dijkstra"
)

type City struct {
	Name     string
	Region   string
	IsHub    bool
	Outbound []dijkstra.Edge[*City, float64]
}

func (c *City) Edges() []dijkstra.Edge[*City, float64] {
	return c.Outbound
}

func main() {
	nyc := &City{Name: "New York", Region: "NA", IsHub: true}
	lon := &City{Name: "London", Region: "EU", IsHub: true}
	fra := &City{Name: "Frankfurt", Region: "EU", IsHub: true}
	tyo := &City{Name: "Tokyo", Region: "APAC", IsHub: true}
	sgp := &City{Name: "Singapore", Region: "APAC", IsHub: true}
	syd := &City{Name: "Sydney", Region: "APAC", IsHub: false}
	mre := &City{Name: "Mars Outpost", Region: "SPACE", IsHub: false}

	nyc.Outbound = []dijkstra.Edge[*City, float64]{
		{To: lon, Weight: 5585.0},
		{To: fra, Weight: 6200.0},
	}
	lon.Outbound = []dijkstra.Edge[*City, float64]{
		{To: fra, Weight: 650.0},
		{To: sgp, Weight: 10880.0},
	}
	fra.Outbound = []dijkstra.Edge[*City, float64]{
		{To: tyo, Weight: 9330.0},
	}
	sgp.Outbound = []dijkstra.Edge[*City, float64]{
		{To: tyo, Weight: 5320.0},
		{To: syd, Weight: 6300.0},
	}
	tyo.Outbound = []dijkstra.Edge[*City, float64]{
		{To: syd, Weight: 7800.0},
	}

	printHeader("GENERIC DIJKSTRA LIBRARY SHOWCASE")

	printSection("1. Point-to-Point Route Search (NYC -> Sydney)")
	opt1 := dijkstra.Search(nyc, dijkstra.By(nyc), dijkstra.By(syd))
	path1, dist1, ok1 := opt1.PathTo(syd)
	printResult(path1, dist1, ok1)

	printSection("2. Full Network Dispatch Table (From NYC to All Destinations)")
	opt2 := dijkstra.Search(nyc, dijkstra.By(nyc), nil)
	allCities := []*City{lon, fra, tyo, sgp, syd, mre}
	for _, target := range allCities {
		path, dist, ok := opt2.PathTo(target)
		fmt.Printf("  Destination: %-12s | ", target.Name)
		if ok {
			var names []string
			for _, p := range path {
				names = append(names, p.Name)
			}
			fmt.Printf("Distance: %8.1f km | Path: %s\n", dist, strings.Join(names, " -> "))
		} else {
			fmt.Printf("Status: UNREACHABLE\n")
		}
	}

	printSection("3. Multi-Source Search (Shortest connection from ANY EU node to Tokyo)")
	euStartSelector := func(c *City) bool {
		return c.Region == "EU"
	}
	opt3 := dijkstra.Search(nyc, euStartSelector, dijkstra.By(tyo))
	path3, dist3, ok3 := opt3.PathTo(tyo)
	printResult(path3, dist3, ok3)

	printSection("4. Dynamic Attribute Matching (Search by Property: First APAC Location)")
	apacSelector := func(c *City) bool {
		return c.Region == "APAC"
	}
	opt4 := dijkstra.Search(nyc, dijkstra.By(nyc), apacSelector)
	for _, c := range []*City{sgp, tyo, syd} {
		if path, dist, ok := opt4.PathTo(c); ok {
			printResult(path, dist, ok)
			break
		}
	}

	printSection("5. Edge Case: Unreachable Node Handling")
	opt5 := dijkstra.Search(nyc, dijkstra.By(nyc), dijkstra.By(mre))
	path5, dist5, ok5 := opt5.PathTo(mre)
	printResult(path5, dist5, ok5)
}

func printHeader(title string) {
	fmt.Println("================================================================================")
	fmt.Printf("  %s\n", title)
	fmt.Println("================================================================================")
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

func printResult(path []*City, dist float64, ok bool) {
	if !ok {
		fmt.Println("  Result: No reachable path found")
		return
	}
	var names []string
	for _, p := range path {
		names = append(names, p.Name)
	}
	fmt.Printf("  Result: Path Found (%d hops, Total Distance: %.1f km)\n", len(path)-1, dist)
	fmt.Printf("  Route : %s\n", strings.Join(names, " -> "))
}
