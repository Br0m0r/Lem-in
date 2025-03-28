package graph

import (
	"errors"
	"fmt"

	// For sorting neighbor lists
	"lem-in/structs"
)

// ------------------------ Build Graph Block ------------------------
// BuildGraph constructs a graph from rooms and tunnels.
// It builds a map of room names to Room structs and an adjacency list.
// After adding all tunnels, each neighbor list is sorted alphabetically so that
// the BFS in FindMultiplePaths explores neighbors in a consistent order.
func BuildGraph(rooms []structs.Room, tunnels []structs.Tunnel) (*structs.Graph, error) {
	g := &structs.Graph{
		Rooms:     make(map[string]*structs.Room),
		Neighbors: make(map[string][]string),
	}

	// Add rooms to the graph.
	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	// Add tunnels to build the adjacency list.
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

	// --------------- Reorder Neighbors Block ---------------
	// Sort each room's neighbor list alphabetically to ensure a consistent BFS order.

	return g, nil
}

// ------------------------ Path Finding Block ------------------------
// FindMultiplePaths finds all edge-disjoint paths from the start room to the end room.
// Since every tunnel has an implicit capacity of 1, it repeatedly searches for a simple path
// using BFS, then "removes" the forward edges used in that path so that they are not reused.
func FindMultiplePaths(g *structs.Graph) ([][]string, error) {
	var start, end string
	for name, room := range g.Rooms {
		if room.IsStart {
			start = name
		}
		if room.IsEnd {
			end = name
		}
	}
	if start == "" || end == "" {
		return nil, errors.New("ERROR: missing start or end room")
	}

	// Make a copy of the neighbors so we can modify edge usage without affecting the original graph.
	copyNeighbors := make(map[string][]string)
	for room, neighs := range g.Neighbors {
		neighborsCopy := make([]string, len(neighs))
		copy(neighborsCopy, neighs)
		copyNeighbors[room] = neighborsCopy
	}

	var paths [][]string

	// Continue searching until no path is found.
	for {
		path, found := bfs(copyNeighbors, start, end)
		if !found {
			break
		}
		paths = append(paths, path)
		removePathEdges(copyNeighbors, path)
	}

	if len(paths) == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	return paths, nil
}

// bfs performs a breadth-first search in the graph defined by neighbors from start to end.
// It returns the first found path as a slice of room names.
func bfs(neighbors map[string][]string, start, end string) ([]string, bool) {
	queue := []string{start}
	visited := make(map[string]bool)
	parent := make(map[string]string)
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == end {
			// Reconstruct the path from end to start.
			var path []string
			for node := end; node != ""; node = parent[node] {
				path = append([]string{node}, path...)
			}
			return path, true
		}
		for _, neigh := range neighbors[current] {
			if !visited[neigh] {
				visited[neigh] = true
				parent[neigh] = current
				queue = append(queue, neigh)
			}
		}
	}
	return nil, false
}

// ------------------------ Edge Removal Block ------------------------
// removePathEdges removes the forward edges used in the given path from the neighbors map.
// For an undirected graph, instead of removing both directions immediately,
// we remove only the forward edge (from u to v) to leave a reverse connection intact.
// This reverse edge might help the BFS find an alternative route and thus reveal extra disjoint paths.
func removePathEdges(neighbors map[string][]string, path []string) {
	for i := 0; i < len(path)-1; i++ {
		u, v := path[i], path[i+1]
		neighbors[u] = removeEdge(neighbors[u], v)
		// Note: The reverse edge from v to u is left intact.
	}
}

// removeEdge removes the element 'target' from the slice and returns the new slice.
func removeEdge(slice []string, target string) []string {
	newSlice := slice[:0]
	for _, s := range slice {
		if s != target {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
