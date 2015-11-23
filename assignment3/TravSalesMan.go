package main

import (
	"container/heap"
	"fmt"
	"math"
)

func DijkstraAlgoCall(start string, end string, Graph []Edge) (nodeOrder []string) {

	directed := true
	findAll := true

	// construct linked representation of example data
	allNodes, startNode, endNode := linkGraph(Graph, directed, start, end)
	if directed {
		fmt.Print("Directed")
	} else {
		fmt.Print("Undirected")
	}
	fmt.Printf(" graph with %d nodes, %d edges\n", len(allNodes), len(Graph))
	if startNode == nil {
		fmt.Printf("start node %q not found in graph\n", start)
		return
	}
	if findAll {
		endNode = nil
	} else if endNode == nil {
		fmt.Printf("end node %q not found in graph\n", end)
		return
	}

	// run Dijkstra's shortest path algorithm
	paths := dijkstra(allNodes, startNode, endNode)
	fmt.Println("Shortest path(s):")
	//var nodeOrder []string
	nodeOrder = make([]string, len(allNodes))
	k := 0
	for _, p := range paths {

		nodeOrder[k] = p.path[len(p.path)-1]
		k++
		//, "length", p.cost)
	}

	fmt.Println(nodeOrder)

	return nodeOrder
}

// node and neighbor structs hold data useful for the heap-optimized
// Dijkstra's shortest path algorithm
type node struct {
	vert string     // vertex name
	tent float64    // tentative cost
	prev *node      // previous node in shortest path back to start
	done bool       // true when tent and prev represent shortest path
	nbs  []neighbor // edges from this vertex
	rx   int        // heap.Remove index
}

type neighbor struct {
	nd   *node   // node corresponding to vertex
	cost float64 // cost to this node (from whatever node references this)
}

// linkGraph constructs a linked representation of the input graph,
// with additional fields needed by the shortest path algorithm.
//
// Return value allNodes will contain all nodes found in the input graph,
// even ones not reachable from the start node.
// Return values startNode, endNode will be nil if the specified start or
// end node names are not found in the graph.
func linkGraph(Graph []Edge, directed bool,
	start, end string) (allNodes []*node, startNode, endNode *node) {

	all := make(map[string]*node)
	// one pass over graph to collect nodes and link neighbors
	for _, e := range Graph {
		n1 := all[e.vert1]
		n2 := all[e.vert2]
		// add previously unseen nodes
		if n1 == nil {
			n1 = &node{vert: e.vert1}
			all[e.vert1] = n1
		}
		if n2 == nil {
			n2 = &node{vert: e.vert2}
			all[e.vert2] = n2
		}
		// link neighbors
		n1.nbs = append(n1.nbs, neighbor{n2, e.cost})
		if !directed {
			n2.nbs = append(n2.nbs, neighbor{n1, e.cost})
		}
	}
	allNodes = make([]*node, len(all))
	var n int
	for _, nd := range all {
		allNodes[n] = nd
		n++
	}
	return allNodes, all[start], all[end]
}

// return type
type Path struct {
	path []string
	cost float64
}

// dijkstra is a heap-enhanced version of Dijkstra's shortest path algorithm.
//
// If endNode is specified, only a single path is returned.
// If endNode is nil, paths to all nodes are returned.
//
// Note input allNodes is needed to efficiently accomplish WP steps 1 and 2.
// This initialization could be done in linkGraph, but is done here to more
// closely follow the WP algorithm.
func dijkstra(allNodes []*node, startNode, endNode *node) (pl []Path) {
	// WP steps 1 and 2.
	for _, nd := range allNodes {
		nd.tent = math.MaxInt32
		nd.done = false
		nd.prev = nil
		nd.rx = -1
	}
	current := startNode
	current.tent = 0
	var unvis ndList

	for {
		// WP step 3: update tentative costs to neighbors
		for _, nb := range current.nbs {
			if nd := nb.nd; !nd.done {
				if d := current.tent + nb.cost; d < nd.tent {
					nd.tent = d
					nd.prev = current
					if nd.rx < 0 {
						heap.Push(&unvis, nd)
					} else {
						heap.Fix(&unvis, nd.rx)
					}
				}
			}
		}
		// WP step 4: mark current node visited, record path and cost
		current.done = true
		if endNode == nil || current == endNode {
			// record path and cost for return value
			cost := current.tent
			// recover path by tracing prev links,
			var p []string
			for ; current != nil; current = current.prev {
				p = append(p, current.vert)
			}
			// then reverse list
			for i := (len(p) + 1) / 2; i > 0; i-- {
				p[i-1], p[len(p)-i] = p[len(p)-i], p[i-1]
			}
			pl = append(pl, Path{p, cost}) // pl is return value
			// WP step 5 (case of end node reached)
			if endNode != nil {
				return
			}
		}
		if len(unvis) == 0 {
			break // WP step 5 (case of no more reachable nodes)
		}
		// WP step 6: new current is node with smallest tentative cost
		current = heap.Pop(&unvis).(*node)
	}
	return
}

// ndList implements container/heap
type ndList []*node

func (n ndList) Len() int           { return len(n) }
func (n ndList) Less(i, j int) bool { return n[i].tent < n[j].tent }
func (n ndList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
	n[i].rx = i
	n[j].rx = j
}
func (n *ndList) Push(x interface{}) {
	nd := x.(*node)
	nd.rx = len(*n)
	*n = append(*n, nd)
}
func (n *ndList) Pop() interface{} {
	s := *n
	last := len(s) - 1
	r := s[last]
	*n = s[:last]
	r.rx = -1
	return r
}
