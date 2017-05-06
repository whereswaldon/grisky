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

var colors = [...]string{"red", "orange", "yellow", "green", "pink", "gold", "chocolate"}

const (
	TOTAL_NODES = 42
	MAX_EDGES   = TOTAL_NODES * 2
)

func getColor() string {
	return colors[rand.Int31n(int32(len(colors)))]
}

func main() {
	rand.Seed(time.Now().Unix())
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
