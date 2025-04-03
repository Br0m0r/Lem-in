package simulation

import (
	"fmt"
	"lem-in/structs"
	"lem-in/visualizer"
	"strings"
)

// initSimulation sets up the simulation state for each path.
// It creates a PathSim for every available path, initializing each ant's starting position and assigning unique IDs.
//
// Variables:
//
//	simStates: a slice that will hold a PathSim for each path in pathList.
//	antIDCounter: a counter starting at 1, ensuring every ant gets a unique identifier.
//	path: one of the paths from pathList, for example ["0", "2", "3", "1"] from example00.txt.
//	antCountForPath: the number of ants assigned to this path, from assignment.AntsPerPath.
//	positions: a slice of integers representing each ant's current position in the path.
//	           - Initially, every value is -1, meaning the ant hasn't started moving.
//	antIDs: a slice of integers holding unique IDs for each ant on the path.
func initSimulation(pathList [][]string, assignment structs.PathAssignment) []structs.PathSim {
	// Create a slice with a length equal to the number of paths.
	simStates := make([]structs.PathSim, len(pathList))
	antIDCounter := 1 // Start assigning ant IDs from 1.

	// Loop over each path in pathList.
	for i, path := range pathList {
		// Get the number of ants for the current path from the assignment.
		antCountForPath := assignment.AntsPerPath[i]

		// Create a slice for positions; each ant is initially set to -1 (has not entered the path).
		positions := make([]int, antCountForPath)
		for j := range positions {
			// Explain: positions[j] = -1 means ant j has not been injected into the path.
			positions[j] = -1
		}

		// Create a slice for ant IDs.
		antIDs := make([]int, antCountForPath)
		for j := 0; j < antCountForPath; j++ {
			// Assign a unique ID to each ant.
			antIDs[j] = antIDCounter
			antIDCounter++ // Increment for the next ant.
		}

		// Save the simulation state for this path.
		// For example, if path is ["0", "2", "3", "1"] and there are 4 ants,
		// then simStates[i] will be:
		//   Path: ["0", "2", "3", "1"]
		//   Positions: [-1, -1, -1, -1]
		//   AntIDs: [1, 2, 3, 4]
		simStates[i] = structs.PathSim{
			Path:      path,
			Positions: positions,
			AntIDs:    antIDs,
		}
	}
	return simStates
}

// isRoomOccupied checks whether any ant is currently at the specified room index in the path.
// Parameters:
//
//	positions: slice of current positions for each ant along a path.
//	roomIndex: the index in the path to check (for example, index 1 corresponds to room "2" in ["0","2","3","1"]).
//
// Returns:
//
//	True if at least one ant is at that index; otherwise, false.
func isRoomOccupied(positions []int, roomIndex int) bool {
	// Loop over every ant's current position.
	for _, pos := range positions {
		// If an ant's position equals roomIndex, that room is occupied.
		if pos == roomIndex {
			return true
		}
	}
	return false
}

// processTurn simulates one turn of ant movements on all paths.
// It attempts to inject new ants into the path or move ants already in the path forward.
// Parameters:
//
//	simStates: a slice containing the current simulation state for each path.
//
// Returns:
//
//	moveDescriptions: a slice of strings describing each move in this turn (e.g., "L1-2").
//	gridVisualization: a string that visually represents the state of all paths after this turn.
func processTurn(simStates []structs.PathSim) (moveDescriptions []string, gridVisualization string) {
	var gridBuilder strings.Builder // To build the visual grid output.

	// Loop over each simulation state (each path) in simStates.
	for idx := range simStates {
		// simState is a pointer to the current path's simulation state.
		simState := &simStates[idx]
		// pathLength is the total number of rooms in the current path.
		// For example, for path ["0","2","3","1"], pathLength is 4.
		pathLength := len(simState.Path)

		// newPositions is a temporary slice to hold updated ant positions for this turn.
		// We copy the current positions into newPositions to calculate changes.
		newPositions := make([]int, len(simState.Positions))
		copy(newPositions, simState.Positions)

		// Check if the path consists of only 2 rooms (start and end).
		if pathLength == 2 {
			// For a direct connection, we try to inject one ant from the start.
			// Loop over each ant assigned to this path.
			for j := 0; j < len(simState.Positions); j++ {
				// If an ant has not started (position is -1), inject it.
				if simState.Positions[j] == -1 {
					// Set the ant's position to 1 (which is the index for the end room in a 2-room path).
					newPositions[j] = 1
					// Record the move.
					// For example, "L1-<end>" where simState.Path[1] is the end room.
					moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1]))
					// Only one ant is injected per turn on a direct path.
					break
				}
			}
		} else {
			// For longer paths (like ["0","2","3","1"]), process ants in reverse order.
			// Processing in reverse ensures that ants ahead move first, which can free up space for those behind.
			for j := len(simState.Positions) - 1; j >= 0; j-- {
				// If an ant hasn't started moving (position is -1), try to inject it into the path.
				if simState.Positions[j] == -1 {
					// Check if the first room after the start (index 1) is free.
					if !isRoomOccupied(newPositions, 1) {
						// Move the ant into the path at index 1.
						newPositions[j] = 1
						// Create a move description such as "L1-2" meaning ant 1 moves to room "2" (which is simState.Path[1]).
						moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1]))
					}
				} else if simState.Positions[j] < pathLength-1 {
					// For ants already in the path but not yet at the end, attempt to move them forward.
					// nextIndex is the next room index along the path.
					nextIndex := simState.Positions[j] + 1
					// Check if the next room is free or if the next room is the end room.
					if nextIndex == pathLength-1 || !isRoomOccupied(newPositions, nextIndex) {
						// Move the ant to the next room.
						newPositions[j] = nextIndex
						// Record the move, e.g., "L2-3" meaning ant 2 moves to room at index nextIndex.
						moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[nextIndex]))
					}
				}
				// Detailed note: each iteration checks one ant. The reverse loop ensures that ants further along the path move first,
				// avoiding blocking situations where an ant behind cannot move because the one in front hasn't moved.
			}
		}
		// After processing all ants for the current path, update simState.Positions with the newPositions computed.
		copy(simState.Positions, newPositions)
		// Use the visualizer to generate a grid view for the current path.
		// The grid shows each room and, in parentheses, which ants are present.
		gridBuilder.WriteString(visualizer.GeneratePathGrid(*simState) + "\n")
	}
	// Return the collected move descriptions and the complete grid visualization for this turn.
	return moveDescriptions, gridBuilder.String()
}

// SimulateMultiPath runs the simulation until no ant moves occur.
// It outputs a detailed grid visualization to "simulation_output.txt" and prints a summary of moves to the terminal.
//
// Variables:
//
//	antTotal: total number of ants (from example00.txt, this is 4).
//	pathList: the available paths (e.g., [["0", "2", "3", "1"]]).
//	assignment: indicates how ants are distributed among the paths (e.g., [4] if all 4 ants use the one path).
//	headerInfo: a text header containing input details and path information, created by the visualizer.
//	simStates: simulation state for each path (each element contains the path, positions, and ant IDs).
//	gridOutputs: slice of grid visualizations for each turn.
//	moveOutputs: slice of move descriptions for each turn.
//	turnCount: the number of simulation turns executed.
func SimulateMultiPath(antTotal int, pathList [][]string, assignment structs.PathAssignment, headerInfo string) {
	// Initialize simulation state using our paths and ant assignment.
	simStates := initSimulation(pathList, assignment)

	// gridOutputs will store the grid view (detailed simulation state) for each turn.
	var gridOutputs []string
	// moveOutputs will store the move summary (simple move info) for each turn.
	var moveOutputs []string
	turnCount := 0 // Initialize turn counter.

	// Run simulation turns repeatedly until no moves are made.
	for {
		// processTurn returns:
		//   moves: a slice of strings describing the moves this turn.
		//   grid: a string visualization of the current state of the path(s).
		moves, grid := processTurn(simStates)

		// If no moves were made, all ants have reached the end; exit the loop.
		if len(moves) == 0 {
			break
		}

		// Append the grid visualization for this turn to gridOutputs.
		gridOutputs = append(gridOutputs, grid)
		// Combine the move descriptions into a single string and append to moveOutputs.
		moveOutputs = append(moveOutputs, strings.Join(moves, " "))
		turnCount++ // Increment the turn counter.
		// Detailed note: Each loop iteration represents one simulation turn where ants are moved.
	}

	// Write the complete simulation output (header, each turn's grid, and total turn count) to "simulation_output.txt".
	err := visualizer.WriteSimulationOutput("simulation_output.txt", headerInfo, gridOutputs, turnCount)
	if err != nil {
		fmt.Println("Error writing simulation output:", err)
	}
	// Print the simple move summary for each turn to the terminal.
	visualizer.PrintTerminalOutput(moveOutputs)
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
