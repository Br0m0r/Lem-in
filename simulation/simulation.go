package simulation

import (
	"fmt"
	"lem-in/structs"
	"lem-in/visualizer"
	"strings"
)

// 1. initSimState
// Initializes the simulation state for each path by setting ant positions to -1 (not injected)
// and assigning unique ant IDs.
func initSimState(paths [][]string, assignment structs.PathAssignment) []structs.PathSim {
	sims := make([]structs.PathSim, len(paths))
	antCounter := 1 // Global ant counter.
	for i, p := range paths {
		count := assignment.AntsPerPath[i] // Number of ants assigned to this path.
		positions := make([]int, count)
		for j := range positions {
			positions[j] = -1 // -1 means the ant is not yet on the path.
		}
		antIDs := make([]int, count) // Assign unique IDs.
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		sims[i] = structs.PathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}
	return sims
}

// 2. isOccupied
// Checks if any ant is already in the room at index 'pos' (used to determine if a move is allowed).
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}

// 3. processTurn
// Processes one simulation turn for all paths.
// It updates ant positions, collects the move strings, and returns a grid visualization.
// Returns a slice of move strings (for terminal output) and the grid visualization string.
func processTurn(sims []structs.PathSim) (turnMoves []string, turnGrid string) {
	var gridBuilder strings.Builder

	// Process each simulation path.
	for idx := range sims {
		sim := &sims[idx]
		pathLen := len(sim.Path)
		newPos := make([]int, len(sim.Positions))
		copy(newPos, sim.Positions)

		if pathLen == 2 {
			// Direct path: try to inject one ant.
			for j := 0; j < len(sim.Positions); j++ {
				if sim.Positions[j] == -1 {
					newPos[j] = 1
					turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					break // Only one ant per turn.
				}
			}
		} else {
			// Multi-step path: process ants in reverse order to avoid collisions.
			for j := len(sim.Positions) - 1; j >= 0; j-- {
				if sim.Positions[j] == -1 {
					// Injection: check if the first move is free.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 {
					next := sim.Positions[j] + 1
					// Move forward if the next room is free or it's the end.
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}
		}
		// Update the simulation state.
		copy(sim.Positions, newPos)
		// Build the grid visualization using the visualizer's GeneratePathGrid.
		gridBuilder.WriteString(visualizer.GeneratePathGrid(*sim) + "\n")
	}
	return turnMoves, gridBuilder.String()
}

// 4. SimulateMultiPath
// The main simulation function. It initializes the state, repeatedly processes turns until no moves occur,
// and then outputs the simulation results using visualizer functions.
func SimulateMultiPath(antCount int, paths [][]string, assignment structs.PathAssignment, extraInfo string) {
	// Initialize simulation state.
	sims := initSimState(paths, assignment)

	var turnOutputs []string     // Detailed grid output for each turn.
	var terminalOutputs []string // Minimal move output for each turn.
	turn := 0

	// Process simulation turns until no ant moves.
	for {
		moves, grid := processTurn(sims)
		if len(moves) == 0 {
			break
		}
		turnOutputs = append(turnOutputs, grid)
		terminalOutputs = append(terminalOutputs, strings.Join(moves, " "))
		turn++
	}

	// Output the simulation results.
	err := visualizer.WriteSimulationOutput("simulation_output.txt", extraInfo, turnOutputs, turn)
	if err != nil {
		fmt.Println("Error writing simulation output:", err)
	}
	visualizer.PrintTerminalOutput(terminalOutputs)
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
