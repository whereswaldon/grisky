package main

import (
	"fmt"
	gv "github.com/awalterschulze/gographviz"
	"github.com/whereswaldon/slijkstra/alg"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	colors = [...]string{"red", "orange", "yellow", "green", "pink", "gold", "chocolate"}
)

const (
	TOTAL_NODES                   = 42
	MAX_EDGES                     = TOTAL_NODES * 2
	MAX_CONTINENTS                = 10
	MIN_CONTINENTS                = 5
	MIN_TERRITORIES_PER_CONTINENT = 3
)

func getColor() string {
	return colors[rand.Int31n(int32(len(colors)))]
}

func main() {
	rand.Seed(time.Now().Unix())
	// Create an array to track the distribution of territories between continents
	numContinents := rand.Int31n(MAX_CONTINENTS-MIN_CONTINENTS) + MIN_CONTINENTS
	continentCounts := make([]int, numContinents)
	continentsNeedingAllocation := TOTAL_NODES - numContinents*MIN_TERRITORIES_PER_CONTINENT
	// randomly allocate territories
	for continentsNeedingAllocation > 0 {
		continentCounts[rand.Int31n(numContinents)]++
		continentsNeedingAllocation--
	}
	// ensure each continent has the minimum quantity of territories
	for i := range continentCounts {
		continentCounts[i] += MIN_TERRITORIES_PER_CONTINENT
	}
	fmt.Fprintln(os.Stderr, continentCounts)
	graph := gv.NewGraph()
	solver := alg.NewGraph(TOTAL_NODES)
	for i := 0; i < TOTAL_NODES; i++ {
		graph.AddNode("G", strconv.Itoa(i), nil)
	}
	for _, node := range graph.Nodes.Nodes {
		node.Attrs.Add("style", "filled")
		node.Attrs.Add("fillcolor", getColor())
	}
	var other, current int
	for i := 0; i < MAX_EDGES; i++ {
		other = int(rand.Int31n(TOTAL_NODES))
		current = i % TOTAL_NODES
		for other == current {
			other = int(rand.Int31n(TOTAL_NODES))
		}
		graph.AddEdge(strconv.Itoa(current),
			strconv.Itoa(other), false, nil)
		solver.InsertEdge(current, other, 1)
	}
	fmt.Println(graph)
	s, e, d := solver.FindDiameter()
	fmt.Fprintf(os.Stderr, "Diameter: %d\nStart: %d\nEnd: %d\n", d, s, e)
}
