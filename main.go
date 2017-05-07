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
	MAX_EDGES                     = TOTAL_NODES * 2
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
	var other, current int
	for i := 0; i < MAX_EDGES; i++ {
		other = int(rand.Int31n(TOTAL_NODES))
		current = i % TOTAL_NODES
		for other == current {
			other = int(rand.Int31n(TOTAL_NODES))
		}
		b.vizgraph.AddEdge(strconv.Itoa(current),
			strconv.Itoa(other), false, nil)
		b.solvgraph.InsertEdge(current, other, 1)
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
	board.assignEdges()
	return board
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
