package graph

import (
	"errors"
	"fmt"

	"lem-in/structs"
)

// 1. BuildGraph
// Constructs a graph from the given rooms and tunnels.
// Arguments:
//   - rooms ([]structs.Room): A slice of room structures with name, coordinates, and start/end flags.
//   - tunnels ([]structs.Tunnel): A slice of tunnel structures linking two room names.
//
// Returns:
//   - *structs.Graph: A pointer to a Graph structure containing:
//   - Rooms: a map (string → *Room) of room names to room data.
//   - Neighbors: a map (string → []string) representing the adjacency list.
//   - error: an error if any tunnel references an unknown room.
func BuildGraph(rooms []structs.Room, tunnels []structs.Tunnel) (*structs.Graph, error) {
	// Create a new Graph with empty maps.
	g := &structs.Graph{
		Rooms:     make(map[string]*structs.Room),
		Neighbors: make(map[string][]string),
	}

	// Add each room to the graph.
	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	// Add tunnels to build the adjacency list.
	for _, tunnel := range tunnels {
		// Verify that both endpoints exist.
		if _, ok := g.Rooms[tunnel.RoomA]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomA)
		}
		if _, ok := g.Rooms[tunnel.RoomB]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomB)
		}
		// Since the graph is undirected, add both directions.
		g.Neighbors[tunnel.RoomA] = append(g.Neighbors[tunnel.RoomA], tunnel.RoomB)
		g.Neighbors[tunnel.RoomB] = append(g.Neighbors[tunnel.RoomB], tunnel.RoomA)
	}

	// (Optional: Sorting neighbor lists could be done here to ensure consistent BFS order.)

	return g, nil
}

// 2. FindMultiplePaths
// Finds all edge-disjoint paths from the start to the end room using a BFS-based approach.
// Arguments:
//   - g (*structs.Graph): The graph constructed from rooms and tunnels.
//
// Returns:
//   - [][]string: A slice of paths, where each path is a slice of room names from start to end.
//   - error: An error if no valid paths are found or if start/end are missing.
func FindMultiplePaths(g *structs.Graph) ([][]string, error) {
	var start, end string
	// Identify start and end room names.
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

	// Make a copy of the neighbors so that we can modify it while finding paths.
	copyNeighbors := make(map[string][]string)
	for room, neighs := range g.Neighbors {
		neighborsCopy := make([]string, len(neighs))
		copy(neighborsCopy, neighs)
		copyNeighbors[room] = neighborsCopy
	}

	var paths [][]string

	// Repeatedly search for a path using BFS.
	for {
		path, found := bfs(copyNeighbors, start, end)
		if !found {
			break
		}
		// Append the found path.
		paths = append(paths, path)
		// Remove the forward edges of the found path so they are not reused.
		removePathEdges(copyNeighbors, path)
	}

	if len(paths) == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	return paths, nil
}

// 3. bfs
// Performs a breadth-first search from start to end using the provided neighbors map.
// Arguments:
//   - neighbors (map[string][]string): The graph's adjacency list (which may be modified).
//   - start (string): The starting room name.
//   - end (string): The target room name.
//
// Returns:
//   - []string: The found path as a slice of room names (if found).
//   - bool: A flag indicating if a path was found (true) or not (false).
func bfs(neighbors map[string][]string, start, end string) ([]string, bool) {
	queue := []string{start}          // Initialize queue with the start node.
	visited := make(map[string]bool)  // Track visited nodes.
	parent := make(map[string]string) // Map to reconstruct the path.
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		// If we reach the end, reconstruct the path.
		if current == end {
			var path []string
			for node := end; node != ""; node = parent[node] {
				// Prepend node to path.
				path = append([]string{node}, path...)
			}
			return path, true
		}
		// Explore unvisited neighbors.
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

// 4. removePathEdges
// Removes the forward edges used in the given path from the neighbors map,
// so that subsequent BFS searches cannot reuse these edges.
// Arguments:
//   - neighbors (map[string][]string): The mutable adjacency list copy.
//   - path ([]string): The path found by bfs, as a slice of room names.
func removePathEdges(neighbors map[string][]string, path []string) {
	// For each consecutive pair of nodes in the path, remove the forward edge.
	for i := 0; i < len(path)-1; i++ {
		u, v := path[i], path[i+1]
		neighbors[u] = removeEdge(neighbors[u], v)
		// Note: The reverse edge from v to u is left intact.
	}
}

// 5. removeEdge
// Removes the target string from a slice of strings and returns the updated slice.
// Arguments:
//   - slice ([]string): The slice from which to remove the element.
//   - target (string): The element to remove.
//
// Returns:
//   - []string: The new slice with the target removed.
func removeEdge(slice []string, target string) []string {
	newSlice := slice[:0] // Use the same underlying array.
	for _, s := range slice {
		if s != target {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice
}
