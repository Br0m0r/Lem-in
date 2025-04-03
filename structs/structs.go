package structs

// Room holds the information for a single room in the ant farm.
// For example, using example00.txt:
//   - Room "0" is defined as "0 0 3" and marked as start.
//   - Room "2" is defined as "2 2 5".
//   - Room "3" is defined as "3 4 0".
//   - Room "1" is defined as "1 8 3" and marked as end.
type Room struct {
	Name    string // The room's name, e.g., "0", "2", "3", or "1".
	X       int    // The X coordinate (used for visualization), e.g., 0 for room "0".
	Y       int    // The Y coordinate (used for visualization), e.g., 3 for room "0".
	IsStart bool   // True if this room is the starting room (e.g., room "0").
	IsEnd   bool   // True if this room is the ending room (e.g., room "1").
}

// Tunnel represents a connection between two rooms.
// In example00.txt, tunnels are defined by lines like "0-2", "2-3", and "3-1".
// This struct stores the names of the two connected rooms.
type Tunnel struct {
	RoomA string // The name of the first room in the tunnel, e.g., "0" in tunnel "0-2".
	RoomB string // The name of the second room, e.g., "2" in tunnel "0-2".
}

// Graph represents the ant farm as a whole.
// It is built from the list of rooms and tunnels parsed from the input.
// The Graph is used to find paths from the start room to the end room.
type Graph struct {
	Rooms map[string]*Room // A map that links room names to their data.
	// For example, after parsing example00.txt, this map will have entries like:
	//   "0" → Room{Name: "0", X: 0, Y: 3, IsStart: true, IsEnd: false}
	//   "2" → Room{Name: "2", X: 2, Y: 5, IsStart: false, IsEnd: false}, etc.
	Neighbors map[string][]string // A map where each room name points to a list of directly connected room names.
	// For example, for example00.txt, the Neighbors map will be:
	//   "0" → ["2"], "2" → ["0", "3"], "3" → ["2", "1"], "1" → ["3"]
}

// PathAssignment holds the result of assigning ants to the available paths.
// After finding paths in the graph, the scheduling algorithm uses this struct.
// In example00.txt, if one path ["0", "2", "3", "1"] is found, this struct will record:
//   - Paths: a list containing that single path.
//   - AntsPerPath: how many ants are assigned to that path.
type PathAssignment struct {
	Paths [][]string // Each element is a path represented as a slice of room names.
	// For example, it might be: [["0", "2", "3", "1"]]
	AntsPerPath []int // The number of ants assigned to each corresponding path.
	// For example, if there are 4 ants, it might be: [4]
}

// PathSim represents the simulation state for a single path during the ant movement simulation.
// It tracks the order of rooms in the path, the current position of each ant on that path,
// and the unique ID for each ant.
type PathSim struct {
	Path []string // The sequence of room names that form the path.
	// For example, a path might be: ["0", "2", "3", "1"]
	Positions []int // The current room index for each ant on the path.
	// A value of -1 indicates the ant hasn't started moving along the path.
	AntIDs []int // Unique identifiers for the ants on this path.
}
