package simulation

import (
	"fmt"
	"os"
	"strings"

	"lem-in/structs"
)

// ------------------------ Grid Visualization Block ------------------------
// GeneratePathGrid creates a 2D grid visualization for a single simulation path.
// It collects all ant IDs present in each room so that every finished ant is proudly displayed.
// For example, if room "3" has L10, L11, and L12, it will print: [ 3 (L10, L11, L12) ]
func GeneratePathGrid(sim structs.PathSim) string {
	var sb strings.Builder
	// Loop through each room in the path.
	for i, room := range sim.Path {
		var antLabels []string
		// For the current room, collect all ant IDs that have reached it.
		for j, pos := range sim.Positions {
			if pos == i {
				antLabels = append(antLabels, fmt.Sprintf("L%d", sim.AntIDs[j]))
			}
		}
		// If there are any ants in the room, print the room with all ant IDs.
		if len(antLabels) > 0 {
			sb.WriteString(fmt.Sprintf("[ %s (%s) ]", room, strings.Join(antLabels, ", ")))
		} else {
			// Otherwise, just print the room name.
			sb.WriteString(fmt.Sprintf("[ %s ]", room))
		}
		// Add an arrow between rooms (except after the last room).
		if i < len(sim.Path)-1 {
			sb.WriteString(" ---> ")
		}
	}
	return sb.String()
}

// ------------------------ Utility Functions Block ------------------------
// isOccupied returns true if any ant occupies the given room index.
// This is our straightforward check for occupancy—if one ant is there, the position is blocked.
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}

// ------------------------ Simulation Initialization Block ------------------------
// SimulateMultiPath simulates ant movement along multiple paths concurrently.
// It prints minimal move information to the terminal (so you can see which ant moves where each turn)
// and writes a detailed 2D grid visualization (with all ants in their respective rooms, including finished ones)
// along with extra info to simulation_output.txt.
func SimulateMultiPath(antCount int, paths [][]string, assignment structs.PathAssignment, extraInfo string) {
	// ------------------- Build Simulation State for Each Path -------------------
	// Each path gets its own simulation state which includes:
	// - The list of room names forming the path,
	// - The positions of ants along the path (initialized to -1, meaning not yet injected),
	// - Unique ant IDs for each ant on the path.
	sims := make([]structs.PathSim, len(paths))
	antCounter := 1 // Global ant counter.
	for i, p := range paths {
		count := assignment.AntsPerPath[i] // Number of ants assigned to this path.
		positions := make([]int, count)    // Their positions; -1 indicates not injected.
		for j := range positions {
			positions[j] = -1 // Initialize all ants as not injected.
		}
		antIDs := make([]int, count) // Unique IDs for each ant on this path.
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		// Create the simulation state for this path.
		sims[i] = structs.PathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}

	// ------------------- Output File Initialization Block -------------------
	// Create or overwrite the output file for detailed simulation output.
	outFile, err := os.Create("simulation_output.txt")
	if err != nil {
		fmt.Println("Error creating simulation_output.txt:", err)
		return
	}
	defer outFile.Close()

	var terminalBuilder strings.Builder // To build our minimal terminal output.
	var fileBuilder strings.Builder     // To build our detailed file output.

	// Write extra information (input data, summary, etc.) at the top of the file.
	fileBuilder.WriteString(extraInfo)
	fileBuilder.WriteString("\n\n")

	// ------------------------ Simulation Turn Loop Block ------------------------
	turn := 0
	for {
		var turnMoves []string       // Collect minimal move info for this turn.
		var turnGrid strings.Builder // Build grid visualization for this turn.

		// ------------------- Process Each Simulation Path -------------------
		for idx, sim := range sims {
			pathLen := len(sim.Path)
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)

			// ------------------- Direct Paths Handling (Path Length == 2) -------------------
			// For direct paths, where there are only two rooms (start and end),
			// we inject one ant per turn since the end is always free.
			if pathLen == 2 {
				// Loop over ants and inject one ant if not already injected.
				for j := 0; j < len(sim.Positions); j++ {
					if sim.Positions[j] == -1 {
						newPos[j] = 1 // Move ant from start (index 0) to end (index 1).
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
						// Inject only one ant per turn.
						break
					}
				}
				// Note: No movement is needed for ants already injected; they remain in the end room.
			} else {
				// ------------------- Multi-step Paths Handling -------------------
				// Process ants in reverse order so that those closer to the end move first.
				for j := len(sim.Positions) - 1; j >= 0; j-- {
					if sim.Positions[j] == -1 {
						// Injection: if the ant is not yet injected, check if the next room (index 1) is free.
						if !isOccupied(newPos, 1) {
							newPos[j] = 1
							turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
						}
					} else if sim.Positions[j] < pathLen-1 {
						next := sim.Positions[j] + 1
						// Move forward if the next room is free.
						// For intermediate moves, standard occupancy check is applied.
						if next == pathLen-1 || !isOccupied(newPos, next) {
							newPos[j] = next
							turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
						}
					}
				}
			}
			// Update the simulation state for this path.
			copy(sims[idx].Positions, newPos)
			// Append the grid visualization for this path to our turn grid.
			turnGrid.WriteString(GeneratePathGrid(sim) + "\n")
		}

		// ------------------- End Simulation Check -------------------
		// If no moves occurred during this turn, the simulation is complete.
		if len(turnMoves) == 0 {
			break
		}

		// ------------------- Record Turn Output -------------------
		// Write the current turn header and grid visualization to the detailed output.
		fileBuilder.WriteString(fmt.Sprintf("TURN %d\n", turn+1))
		fileBuilder.WriteString(turnGrid.String())
		fileBuilder.WriteString("\n")
		// Append minimal turn moves to the terminal output.
		terminalBuilder.WriteString(fmt.Sprintf("Turn %d: %s\n", turn+1, strings.Join(turnMoves, " ")))
		turn++
	}

	// ------------------- Finalization Block -------------------
	// Write the total number of turns to the output file.
	fileBuilder.WriteString(fmt.Sprintf("Total turns: %d\n", turn))
	_, err = outFile.WriteString(fileBuilder.String())
	if err != nil {
		fmt.Println("Error writing to simulation_output.txt:", err)
	}

	// Print the minimal moves to the terminal.
	fmt.Print(terminalBuilder.String())
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
