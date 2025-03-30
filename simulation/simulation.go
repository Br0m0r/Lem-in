package simulation

import (
	"fmt"
	"lem-in/structs"
	"lem-in/visualizer"
	"strings"
)

// 1. initSimState
// Initializes the simulation state for each path.
// Arguments:
//   - paths ([][]string): Each path is a slice of room names.
//   - assignment (structs.PathAssignment): Contains the paths and a slice (AntsPerPath)
//     indicating how many ants should go on each path.
//
// Returns:
//   - []structs.PathSim: A slice of simulation state objects for each path.
//
// Variables:
//   - sims: Slice of PathSim, where each PathSim holds:
//   - Path: the slice of room names for that path,
//   - Positions: a slice of int (each initialized to -1, meaning the ant is not yet on the path),
//   - AntIDs: a slice of unique ant identifiers.
//   - antCounter: An integer used to assign unique IDs to each ant.
func initSimState(paths [][]string, assignment structs.PathAssignment) []structs.PathSim {
	sims := make([]structs.PathSim, len(paths))
	antCounter := 1 // Global counter for ant IDs.
	for i, p := range paths {
		count := assignment.AntsPerPath[i] // Number of ants for this path.
		positions := make([]int, count)
		for j := range positions {
			positions[j] = -1 // -1 indicates that the ant has not yet been injected into the path.
		}
		antIDs := make([]int, count)
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		sims[i] = structs.PathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}
	return sims
}

// 2. isOccupied
// Checks whether any ant is currently occupying a given room position (index) in a path.
// Arguments:
//   - positions ([]int): The slice representing the current positions of ants.
//   - pos (int): The room index to check.
//
// Returns:
//   - bool: true if any ant is at that index; false otherwise.
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}

// 3. processTurn
// Processes a single simulation turn for all paths by attempting to inject or move ants forward.
// Arguments:
//   - sims ([]structs.PathSim): The simulation state for each path.
//
// Returns:
//   - turnMoves ([]string): A slice of strings indicating moves made this turn (e.g., "L3-room").
//   - turnGrid (string): A string representation (grid visualization) of all paths after this turn.
//
// Variables:
//   - gridBuilder: A strings.Builder used to accumulate the grid visualization for this turn.
//   - newPos: A temporary slice used to store updated positions before committing them to the simulation state.
func processTurn(sims []structs.PathSim) (turnMoves []string, turnGrid string) {
	var gridBuilder strings.Builder

	// Iterate through each simulation path.
	for idx := range sims {
		sim := &sims[idx]
		pathLen := len(sim.Path)
		newPos := make([]int, len(sim.Positions))
		copy(newPos, sim.Positions)

		if pathLen == 2 {
			// Direct path (start and end only): try to inject one ant.
			for j := 0; j < len(sim.Positions); j++ {
				if sim.Positions[j] == -1 {
					newPos[j] = 1
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					break // Only one ant is injected per turn.
				}
			}
		} else {
			// For multi-step paths: process ants in reverse order to avoid blocking.
			for j := len(sim.Positions) - 1; j >= 0; j-- {
				if sim.Positions[j] == -1 {
					// Injection: if the first move (from start to room index 1) is free.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 {
					next := sim.Positions[j] + 1
					// Move forward if next room is free or if it's the end room.
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}
		}
		// Update the simulation state's positions for this path.
		copy(sim.Positions, newPos)
		// Generate a grid visualization of this path using visualizer.GeneratePathGrid.
		gridBuilder.WriteString(visualizer.GeneratePathGrid(*sim) + "\n")
	}
	return turnMoves, gridBuilder.String()
}

// 4. SimulateMultiPath
// Orchestrates the entire simulation process. It initializes the simulation state,
// repeatedly processes turns until no moves occur, and outputs the simulation results.
// Arguments:
//   - antCount (int): Total number of ants in the simulation.
//   - paths ([][]string): The available paths (each is a slice of room names).
//   - assignment (structs.PathAssignment): Contains the ant distribution for each path.
//   - extraInfo (string): A header string with simulation input and summary details.
//
// Returns:
//   - None (the function outputs results to file and terminal).
//
// Variables:
//   - sims: The simulation state (slice of PathSim) returned by initSimState.
//   - turnOutputs: A slice of strings representing the grid visualization for each turn.
//   - terminalOutputs: A slice of strings for minimal move outputs per turn.
//   - turn: The counter for the number of turns processed.
func SimulateMultiPath(antCount int, paths [][]string, assignment structs.PathAssignment, extraInfo string) {
	// Initialize the simulation state.
	sims := initSimState(paths, assignment)

	var turnOutputs []string     // To accumulate detailed grid outputs per turn.
	var terminalOutputs []string // To accumulate move information per turn.
	turn := 0

	// Continue processing turns until no ant moves occur.
	for {
		moves, grid := processTurn(sims)
		// If no moves were made during this turn, the simulation is complete.
		if len(moves) == 0 {
			break
		}
		turnOutputs = append(turnOutputs, grid)
		terminalOutputs = append(terminalOutputs, strings.Join(moves, " "))
		turn++
	}

	// Output the simulation results:
	// Write a detailed output file and print minimal moves to the terminal.
	err := visualizer.WriteSimulationOutput("simulation_output.txt", extraInfo, turnOutputs, turn)
	if err != nil {
		fmt.Println("Error writing simulation output:", err)
	}
	visualizer.PrintTerminalOutput(terminalOutputs)
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
