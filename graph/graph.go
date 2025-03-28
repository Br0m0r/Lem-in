package graph

import (
	"fmt"

	"lem-in/structs"
)

// BuildGraph constructs a graph from rooms and tunnels.
// It builds a map of room names to Room structs and an adjacency list.
func BuildGraph(rooms []structs.Room, tunnels []structs.Tunnel) (*structs.Graph, error) {
	g := &structs.Graph{
		Rooms:     make(map[string]*structs.Room),
		Neighbors: make(map[string][]string),
	}

	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	for _, tunnel := range tunnels {
		if _, ok := g.Rooms[tunnel.RoomA]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomA)
		}
		if _, ok := g.Rooms[tunnel.RoomB]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomB)
		}
		g.Neighbors[tunnel.RoomA] = append(g.Neighbors[tunnel.RoomA], tunnel.RoomB)
		g.Neighbors[tunnel.RoomB] = append(g.Neighbors[tunnel.RoomB], tunnel.RoomA)
	}

	return g, nil
}

// ---------------------
// Max-Flow Functions
// ---------------------

// BuildFlowNetwork creates a flow network from the graph.
// Each edge gets a capacity of 1.
// func BuildFlowNetwork(g *structs.Graph) *structs.FlowNetwork {
// 	network := &structs.FlowNetwork{
// 		Adjacency: make(map[string][]*structs.FlowEdge),
// 	}
// 	for node := range g.Rooms {
// 		network.Nodes = append(network.Nodes, node)
// 	}
// 	for room, neighbors := range g.Neighbors {
// 		for _, neighbor := range neighbors {
// 			edge := &structs.FlowEdge{From: room, To: neighbor, Capacity: 1, Flow: 0}
// 			network.Adjacency[room] = append(network.Adjacency[room], edge)
// 		}
// 	}
// 	return network
// }

// EdmondsKarp runs the max-flow algorithm to find the maximum number of edge-disjoint paths.
// Each found augmenting path increases the flow by 1.
// func EdmondsKarp(network *structs.FlowNetwork, start, end string) int {
// 	maxFlow := 0
// 	for {
// 		parent := make(map[string]*structs.FlowEdge)
// 		queue := []string{start}
// 		for len(queue) > 0 && parent[end] == nil {
// 			current := queue[0]
// 			queue = queue[1:]
// 			for _, edge := range network.Adjacency[current] {
// 				residual := edge.Capacity - edge.Flow
// 				if residual > 0 && parent[edge.To] == nil && edge.To != start {
// 					parent[edge.To] = edge
// 					queue = append(queue, edge.To)
// 					if edge.To == end {
// 						break
// 					}
// 				}
// 			}
// 		}
// 		if parent[end] == nil {
// 			break
// 		}
// 		for node := end; node != start; {
// 			edge := parent[node]
// 			edge.Flow += 1
// 			node = edge.From
// 		}
// 		maxFlow++
// 	}
// 	return maxFlow
// }

// ExtractPaths retrieves all edge-disjoint paths with flow from start to end.
// It removes used flow as paths are extracted.
func ExtractPaths(g *structs.Graph) ([][]string, error) {
	var start, end string

	// Find start and end rooms
	for name, room := range g.Rooms {
		if room.IsStart {
			start = name
		}
		if room.IsEnd {
			end = name
		}
	}

	// If start or end room is missing, return an error
	if start == "" || end == "" {
		return nil, fmt.Errorf("start or end room is missing")
	}

	// BFS queue (each element is a path)
	queue := [][]string{{start}}
	// Store selected paths
	var selectedPaths [][]string
	// Flag for used second rooms
	usedSecondRooms := make(map[string]bool)

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:] // Dequeue

		// Check if we reached the end
		if path[len(path)-1] == end {
			if len(path) < 3 { // Ensure valid path with second room
				continue
			}

			secondRoom := path[1] // Second room after "start"

			// If second room is already used, skip this path
			if usedSecondRooms[secondRoom] {
				continue
			}

			// Mark second room as used and save this path
			usedSecondRooms[secondRoom] = true
			selectedPaths = append(selectedPaths, path)

			continue
		}

		// Explore neighbors
		current := path[len(path)-1]
		for _, neighbor := range g.Neighbors[current] {
			if !contains(path, neighbor) { // Avoid revisiting nodes
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}
	fmt.Println(selectedPaths)
	return selectedPaths, nil
}

// Helper function: Check if a path contains a room
func contains(path []string, room string) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
}

// FindMultiplePaths finds all edge-disjoint paths from start to end using the max-flow approach.
// func FindMultiplePaths(g *structs.Graph) ([][]string, error) {

// }
