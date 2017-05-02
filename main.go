package main

import (
	"fmt"
	gv "github.com/awalterschulze/gographviz"
	"math/rand"
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
	for i := 0; i < TOTAL_NODES; i++ {
		graph.AddNode("G", strconv.Itoa(i), nil)
	}
	for _, node := range graph.Nodes.Nodes {
		node.Attrs.Add("style", "filled")
		node.Attrs.Add("fillcolor", getColor())
	}
	for i := 0; i < MAX_EDGES; i++ {
		other := int(rand.Int31n(TOTAL_NODES))
		for other == i%TOTAL_NODES {
			other = int(rand.Int31n(TOTAL_NODES))
		}
		graph.AddEdge(strconv.Itoa(i%TOTAL_NODES),
			strconv.Itoa(other), false, nil)
	}
	fmt.Println(graph)
}
