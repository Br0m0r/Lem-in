package graph

import (
	"errors"
	"fmt"

	"lem-in/structs"
)

// BuildGraph creates a graph (map) of the ant farm using the list of rooms and tunnels.
// Example with example00.txt:
//
//	Parsed Rooms: "0" (start, 0,3), "2" (2,5), "3" (4,0), "1" (end, 8,3)
//	Parsed Tunnels: "0-2", "2-3", "3-1"
//
// After building, the graph will have:
//
//	Rooms map: keys "0", "2", "3", "1" mapping to their respective room data.
//	Neighbors map: "0": ["2"], "2": ["0", "3"], "3": ["2", "1"], "1": ["3"]
func BuildGraph(roomList []structs.Room, connections []structs.Tunnel) (*structs.Graph, error) {
	// Create a new graph with empty maps.
	// "graphData" will store all room details and the list of their direct connections.
	graphData := &structs.Graph{
		Rooms:     make(map[string]*structs.Room), // Keys will be room names.
		Neighbors: make(map[string][]string),      // Each room name will map to a list of connected room names.
	}

	// Add each room to the graph.
	// For each room from the parsed roomList (from example00.txt, these are "0", "2", "3", "1"),
	// we store the room data using the room's name as the key.
	for i := range roomList {
		currentRoom := roomList[i]
		graphData.Rooms[currentRoom.Name] = &currentRoom
	}

	// Process each tunnel and update the list of connected rooms.
	// For each tunnel (e.g., "0-2", "2-3", "3-1") in the parsed connections:
	for _, tunnel := range connections {
		// Check that both room names in the tunnel exist in the Rooms map.
		if _, found := graphData.Rooms[tunnel.RoomA]; !found {
			return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomA)
		}
		if _, found := graphData.Rooms[tunnel.RoomB]; !found {
			return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomB)
		}
		// Since tunnels work in both directions, add each room to the other's neighbor list.
		// For tunnel "0-2": add "2" as a neighbor of "0", and "0" as a neighbor of "2".
		graphData.Neighbors[tunnel.RoomA] = append(graphData.Neighbors[tunnel.RoomA], tunnel.RoomB)
		graphData.Neighbors[tunnel.RoomB] = append(graphData.Neighbors[tunnel.RoomB], tunnel.RoomA)
	}

	// At this point, using our example00.txt:
	//   Rooms map: "0" -> Room "0", "2" -> Room "2", "3" -> Room "3", "1" -> Room "1"
	//   Neighbors map: "0": ["2"], "2": ["0", "3"], "3": ["2", "1"], "1": ["3"]

	return graphData, nil
}

// FindMultiplePaths finds all separate paths (without reusing tunnels) from the start room to the end room
// using a breadth-first search (BFS) approach.
// Example with example00.txt:
//
//	Start room is "0" (as flagged by IsStart) and end room is "1" (flagged as IsEnd).
//	With our graph, bfs will find the path: ["0", "2", "3", "1"].
//	Then, the tunnel connections used in that path will be removed.
func FindMultiplePaths(graphData *structs.Graph) ([][]string, error) {
	var startRoom, endRoom string
	// Identify the start and end room names from the Rooms map.
	for name, roomData := range graphData.Rooms {
		if roomData.IsStart {
			startRoom = name
		}
		if roomData.IsEnd {
			endRoom = name
		}
	}
	if startRoom == "" || endRoom == "" {
		return nil, errors.New("ERROR: missing start or end room")
	}

	// Make a copy of the Neighbors map (connections) so that we can modify it as we remove used tunnels.
	// For example, the copy initially is:
	//   "0": ["2"], "2": ["0", "3"], "3": ["2", "1"], "1": ["3"]
	connectionsCopy := make(map[string][]string)
	for roomName, connectedList := range graphData.Neighbors {
		newList := make([]string, len(connectedList))
		copy(newList, connectedList)
		connectionsCopy[roomName] = newList
	}

	var foundPaths [][]string

	// Repeatedly search for a new path using BFS.
	// For our example, the first (and only) BFS call will find ["0", "2", "3", "1"].
	for {
		path, pathFound := bfs(connectionsCopy, startRoom, endRoom)
		if !pathFound {
			break
		}
		// Add the found path to our list.
		foundPaths = append(foundPaths, path)
		// Remove the tunnels used in this path so they cannot be used again.
		removePathEdges(connectionsCopy, path)
	}

	if len(foundPaths) == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	return foundPaths, nil
}

// bfs performs a breadth-first search to find one path from the startRoom to the endRoom.
// Example with our graph from example00.txt:
//
//	It starts at "0", visits "2", then "3", and finally "1", reconstructing the path as ["0", "2", "3", "1"].
func bfs(connections map[string][]string, startRoom, endRoom string) ([]string, bool) {
	queue := []string{startRoom}          // Start with the startRoom (e.g., "0").
	visited := make(map[string]bool)      // To track visited rooms.
	parentRoom := make(map[string]string) // To record how we reached each room.
	visited[startRoom] = true

	for len(queue) > 0 {
		currentRoom := queue[0]
		queue = queue[1:]
		// If we've reached the end room (e.g., "1"), rebuild the path.
		if currentRoom == endRoom {
			var path []string
			// Reconstruct the path by going backwards from endRoom using the parentRoom map.
			for room := endRoom; room != ""; room = parentRoom[room] {
				// Prepend each room to build the path in order.
				path = append([]string{room}, path...)
			}
			return path, true
		}
		// Check every room directly connected to the current room.
		for _, nextRoom := range connections[currentRoom] {
			if !visited[nextRoom] {
				visited[nextRoom] = true
				parentRoom[nextRoom] = currentRoom
				queue = append(queue, nextRoom)
			}
		}
	}
	return nil, false
}

// removePathEdges removes the tunnels used in a given path from the connections map.
// This prevents the same tunnel from being used in another path.
// Example with path ["0", "2", "3", "1"]:
//
//	For "0" -> "2": remove "2" from the list of neighbors for "0".
//	For "2" -> "3": remove "3" from the list of neighbors for "2".
//	For "3" -> "1": remove "1" from the list of neighbors for "3".
//
// The reverse connections (e.g., "2" still has "0") are kept.
func removePathEdges(connections map[string][]string, roomPath []string) {
	for i := 0; i < len(roomPath)-1; i++ {
		fromRoom, toRoom := roomPath[i], roomPath[i+1]
		connections[fromRoom] = removeEdge(connections[fromRoom], toRoom)
	}
}

// removeEdge removes a specific room name from a list of connected rooms.
// Example: removeEdge(["2"], "2") will return an empty list, removing the connection.
func removeEdge(connectionList []string, roomNameToRemove string) []string {
	newList := connectionList[:0] // Reuse the same underlying array.
	for _, roomName := range connectionList {
		if roomName != roomNameToRemove {
			newList = append(newList, roomName)
		}
	}
	return newList
}
