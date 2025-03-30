package structs

// Room represents a room in the ant farm.
// Created during parsing (in parser.ParseInputFile), Rooms are used to build the Graph.
// They are later used in visualization and simulation to represent individual nodes.
type Room struct {
	Name    string // Room name.
	X       int    // X coordinate (used for visualization).
	Y       int    // Y coordinate (used for visualization).
	IsStart bool   // True if this is the start room.
	IsEnd   bool   // True if this is the end room.
}

// Tunnel represents a connection between two rooms.
// Tunnels are created during parsing (in parser.ParseInputFile) from input definitions.
// They are used in building the Graph to establish links between Rooms.
type Tunnel struct {
	RoomA string // Name of the first room.
	RoomB string // Name of the second room.
}

// Graph represents the ant farm as an adjacency list.
// The Graph is built from Rooms and Tunnels in BuildGraph.
// It is used for path finding (in FindMultiplePaths) to determine available routes.
type Graph struct {
	Rooms     map[string]*Room    // Mapping from room name to Room.
	Neighbors map[string][]string // Mapping from room name to adjacent room names.
}

// PathAssignment holds the distribution of ants among the available paths.
// It is created by AssignAnts after paths are found.
// This struct is used in simulation (SimulateMultiPath) and visualization to understand how ants are distributed.
type PathAssignment struct {
	Paths       [][]string // Each path is represented as a slice of room names.
	AntsPerPath []int      // Number of ants assigned to each corresponding path.
}

// PathSim represents the simulation state for a single path.
// It is initialized in initSimState (simulation.go) for each path.
// It tracks the sequence of rooms in the path, the current positions of ants (with -1 meaning not yet injected),
// and the unique ant IDs. This struct is updated turn by turn in SimulateMultiPath.
type PathSim struct {
	Path      []string // Sequence of room names forming the path.
	Positions []int    // Current position index for each ant (-1 means not yet injected).
	AntIDs    []int    // Global ant IDs assigned to ants on this path.
}
