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
	colors = [...]string{"red", "orange", "yellow", "green", "pink", "gold", "chocolate", "blue", "white", "gray"}
)

const (
	TOTAL_NODES                   = 42
	MAX_CONTINENTS                = 10
	MIN_CONTINENTS                = 5
	MIN_TERRITORIES_PER_CONTINENT = 3
)

func getColor() string {
	return colors[rand.Int31n(int32(len(colors)))]
}

type Board struct {
	numContinents   int
	continentCounts []int
	vizgraph        *gv.Graph
	solvgraph       *alg.Graph
}

func (b *Board) allocateContinents() {
	// ensure each continent has the minimum quantity of territories
	for i := range b.continentCounts {
		b.continentCounts[i] += MIN_TERRITORIES_PER_CONTINENT
	}
	continentsNeedingAllocation := TOTAL_NODES - b.numContinents*MIN_TERRITORIES_PER_CONTINENT
	// randomly allocate territories
	for continentsNeedingAllocation > 0 {
		b.continentCounts[rand.Int31n(int32(b.numContinents))]++
		continentsNeedingAllocation--
	}
	fmt.Fprintln(os.Stderr, b.continentCounts)
}

func (b *Board) createGraph() {
	graph := gv.NewGraph()
	solver := alg.NewGraph(TOTAL_NODES)
	for i := 0; i < TOTAL_NODES; i++ {
		graph.AddNode("G", strconv.Itoa(i), nil)
	}
	currentNode := 0
	for i := 0; i < b.numContinents; i++ {
		currentColor := colors[i]
		for k := 0; k < b.continentCounts[i]; k++ {
			node := graph.Nodes.Nodes[currentNode]
			currentNode++
			node.Attrs.Add("style", "filled")
			node.Attrs.Add("fillcolor", currentColor)
		}
	}
	b.vizgraph = graph
	b.solvgraph = solver
}

func (b *Board) assignEdges() {
	const RETRIES = 5
	var other, current int
	currentContinentOffset := 0
	for i := 0; i < b.numContinents; i++ {
		numExtraEdges := chooseNumberExtraEdges(b.continentCounts[i])
	PerContinentLoop:
		for k := 0; k < numExtraEdges; k++ {
			current = int(rand.Int31n(int32(b.continentCounts[i]))) + currentContinentOffset
			other = int(rand.Int31n(int32(b.continentCounts[i]))) + currentContinentOffset
			for m := 0; m < RETRIES && (current == other || b.solvgraph.HasEdge(current, other)); m++ {
				other = int(rand.Int31n(int32(b.continentCounts[i]))) + currentContinentOffset
				if m == RETRIES-1 {
					continue PerContinentLoop
				}
			}
			fmt.Fprintf(os.Stderr, "i: %d k: %d current: %d other: %d\n", i, k, current, other)
			b.AddEdge(current, other)
		}
		currentContinentOffset += b.continentCounts[i]
	}
}

// makeEdgeCycles ensures that a continent starts with all territories
// connected in a cycle
func (b *Board) makeEdgeCycles() {
	var other, current int
	currentContinentOffset := 0
	for i := 0; i < b.numContinents; i++ {
		//start with a cycle
		for k := 0; k < b.continentCounts[i]; k++ {
			current = k + currentContinentOffset
			other = ((k + 1) % b.continentCounts[i]) + currentContinentOffset
			b.AddEdge(current, other)
		}
		currentContinentOffset += b.continentCounts[i]
	}
}

func MakeBoard(size int) *Board {
	// Create an array to track the distribution of territories between continents
	numContinents := rand.Int31n(MAX_CONTINENTS-MIN_CONTINENTS) + MIN_CONTINENTS
	continentCounts := make([]int, numContinents)
	board := &Board{
		numContinents:   int(numContinents),
		continentCounts: continentCounts,
	}
	board.allocateContinents()
	board.createGraph()
	board.makeEdgeCycles()
	board.assignEdges()
	return board
}

// chooseNumberExtraEdges selects a random number of additional edges to add to a
// continent between 2 and (n^2-3n)/2, which is the maximum number of unique
// edges in the continent
func chooseNumberExtraEdges(numTerritories int) int {
	switch numTerritories {
	case 3:
		return 0
	case 4:
		return 1
	}
	return int(rand.Int31n(int32(((numTerritories*numTerritories)-3*numTerritories)/2)-2)) + 2
}
func (b *Board) AddEdge(u, v int) {
	b.vizgraph.AddEdge(strconv.Itoa(u), strconv.Itoa(v), false, nil)
	b.solvgraph.InsertEdge(u, v, 1)
}

func (b *Board) String() string {
	return b.vizgraph.String()
}

func main() {
	rand.Seed(time.Now().Unix())
	board := MakeBoard(TOTAL_NODES)
	fmt.Println(board)
	/*
		s, e, d := solver.FindDiameter()
		fmt.Fprintf(os.Stderr, "Diameter: %d\nStart: %d\nEnd: %d\n", d, s, e)
	*/
}
