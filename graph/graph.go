package graph

import (
    "errors"
    "fmt"

    "lem-in/structs"
)

// BuildGraph creates a graph from rooms and tunnels.
func BuildGraph(roomList []structs.Room, connections []structs.Tunnel) (*structs.Graph, error) {
    graphData := &structs.Graph{
        Rooms:     make(map[string]*structs.Room),
        Neighbors: make(map[string][]string),
    }
    for i := range roomList {
        currentRoom := roomList[i]
        graphData.Rooms[currentRoom.Name] = &currentRoom
    }
    for _, tunnel := range connections {
        if _, found := graphData.Rooms[tunnel.RoomA]; !found {
            return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomA)
        }
        if _, found := graphData.Rooms[tunnel.RoomB]; !found {
            return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomB)
        }
        graphData.Neighbors[tunnel.RoomA] = append(graphData.Neighbors[tunnel.RoomA], tunnel.RoomB)
        graphData.Neighbors[tunnel.RoomB] = append(graphData.Neighbors[tunnel.RoomB], tunnel.RoomA)
    }
    return graphData, nil
}

// FindMultiplePaths finds disjoint paths from start to end using repeated BFS.
func FindMultiplePaths(graphData *structs.Graph) ([][]string, error) {
    var startRoom, endRoom string
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

    // copy neighbors to allow edge removal
    connectionsCopy := make(map[string][]string)
    for roomName, neighbors := range graphData.Neighbors {
        newList := make([]string, len(neighbors))
        copy(newList, neighbors)
        connectionsCopy[roomName] = newList
    }

    var foundPaths [][]string
    for {
        path, pathFound := bfs(connectionsCopy, startRoom, endRoom)
        if !pathFound {
            break
        }
        foundPaths = append(foundPaths, path)
        removePathEdges(connectionsCopy, path)
    }

    if len(foundPaths) == 0 {
        return nil, errors.New("ERROR: no valid paths found")
    }
    return foundPaths, nil
}

// bfs performs a breadth-first search for a single path.
func bfs(connections map[string][]string, startRoom, endRoom string) ([]string, bool) {
    queue := []string{startRoom}
    visited := make(map[string]bool)
    parentRoom := make(map[string]string)
    visited[startRoom] = true

    for len(queue) > 0 {
        currentRoom := queue[0]
        queue = queue[1:]
        if currentRoom == endRoom {
            var path []string
            for room := endRoom; room != ""; room = parentRoom[room] {
                path = append([]string{room}, path...)
            }
            return path, true
        }
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

// removePathEdges removes edges used by a path to prevent reuse.
func removePathEdges(connections map[string][]string, roomPath []string) {
    for i := 0; i < len(roomPath)-1; i++ {
        fromRoom, toRoom := roomPath[i], roomPath[i+1]
        connections[fromRoom] = removeEdge(connections[fromRoom], toRoom)
    }
}

// removeEdge removes a target room from a list of connections.
func removeEdge(connectionList []string, roomNameToRemove string) []string {
    newList := connectionList[:0]
    for _, roomName := range connectionList {
        if roomName != roomNameToRemove {
            newList = append(newList, roomName)
        }
    }
    return newList
}
