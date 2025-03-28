package structs

// Room represents a room in the ant farm.
type Room struct {
	Name    string // Room name
	X       int    // X coordinate (used for visualization)
	Y       int    // Y coordinate (used for visualization)
	IsStart bool   // True if this is the start room
	IsEnd   bool   // True if this is the end room
}

// Tunnel represents a connection between two rooms.
type Tunnel struct {
	RoomA string // Name of the first room
	RoomB string // Name of the second room
}

// Graph represents the ant farm as an adjacency list.
type Graph struct {
	Rooms     map[string]*Room    // Mapping from room name to Room
	Neighbors map[string][]string // Mapping from room name to adjacent room names
}

// PathAssignment holds the distribution of ants among the available paths.
type PathAssignment struct {
	Paths       [][]string // Each path represented as a slice of room names.
	AntsPerPath []int      // Number of ants assigned to each corresponding path.
}

// PathSim represents the simulation state for a single path.
// It tracks the list of rooms in the path, the positions of ants (by index),
// and the global ant IDs assigned to each ant on the path.
type PathSim struct {
	Path      []string // Sequence of room names forming the path
	Positions []int    // Current position index for each ant (-1 means not yet injected)
	AntIDs    []int    // Global ant IDs assigned to ants on this path
}
