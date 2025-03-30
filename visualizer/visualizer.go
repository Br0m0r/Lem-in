package visualizer

import (
	"fmt"
	"os"
	"strings"

	"lem-in/structs"
)

// 1. buildInputData
// Returns a string with the raw input data (ant count, room definitions, tunnels).
func buildInputData(antCount int, rooms []structs.Room, tunnels []structs.Tunnel) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", antCount))
	for _, room := range rooms {
		if room.IsStart {
			sb.WriteString("##start\n")
		}
		if room.IsEnd {
			sb.WriteString("##end\n")
		}
		sb.WriteString(fmt.Sprintf("%s %d %d\n", room.Name, room.X, room.Y))
	}
	for _, tunnel := range tunnels {
		sb.WriteString(fmt.Sprintf("%s-%s\n", tunnel.RoomA, tunnel.RoomB))
	}
	return sb.String()
}

// 2. buildSummary
// Returns a simple summary of simulation parameters (number of ants, rooms, tunnels, start and end rooms).
func buildSummary(antCount int, rooms []structs.Room, tunnels []structs.Tunnel) string {
	var sb strings.Builder
	sb.WriteString("----------- Summary -----------\n")
	sb.WriteString(fmt.Sprintf("Number of ants: %d\n", antCount))
	sb.WriteString(fmt.Sprintf("Number of rooms: %d\n", len(rooms)))
	sb.WriteString(fmt.Sprintf("Number of tunnels: %d\n", len(tunnels)))
	var startRoom, endRoom string
	for _, room := range rooms {
		if room.IsStart {
			startRoom = room.Name
		}
		if room.IsEnd {
			endRoom = room.Name
		}
	}
	sb.WriteString(fmt.Sprintf("Start room: %s\n", startRoom))
	sb.WriteString(fmt.Sprintf("End room: %s\n", endRoom))
	return sb.String()
}

// 3. buildAllPaths
// Returns a string that lists all the paths found.
func buildAllPaths(paths [][]string) string {
	var sb strings.Builder
	sb.WriteString("---------- All Found Paths ----------\n")
	sb.WriteString(fmt.Sprintf("Number of possible paths: %d\n", len(paths)))
	for i, p := range paths {
		sb.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(p, " -> ")))
	}
	return sb.String()
}

// 4. buildSelectedPaths
// Returns a string that lists the paths selected after ant assignment.
func buildSelectedPaths(assignment structs.PathAssignment) string {
	var sb strings.Builder
	sb.WriteString("---------- Selected Paths ---------- \n")
	for i, p := range assignment.Paths {
		sb.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(p, " -> ")))
	}
	return sb.String()
}

// 5. PrintExtraInfo
// Combines raw input data, summary, all paths, and selected paths into one extra info string.
// This is typically invoked at the start of the simulation.
func PrintExtraInfo(antCount int, rooms []structs.Room, tunnels []structs.Tunnel, paths [][]string, assignment structs.PathAssignment) string {
	var sb strings.Builder

	sb.WriteString(buildInputData(antCount, rooms, tunnels))
	sb.WriteString("\n")
	sb.WriteString(buildSummary(antCount, rooms, tunnels))
	sb.WriteString("\n")
	sb.WriteString(buildAllPaths(paths))
	sb.WriteString("\n")
	sb.WriteString(buildSelectedPaths(assignment))
	sb.WriteString("\n")

	return sb.String()
}

// 6. GeneratePathGrid
// Creates a grid visualization of a simulation path showing each room and the ants present.
// This function is called during each simulation turn.
func GeneratePathGrid(sim structs.PathSim) string {
	var sb strings.Builder
	for i, room := range sim.Path {
		var antLabels []string
		for j, pos := range sim.Positions {
			if pos == i {
				antLabels = append(antLabels, fmt.Sprintf("L%d", sim.AntIDs[j]))
			}
		}
		if len(antLabels) > 0 {
			sb.WriteString(fmt.Sprintf("[ %s (%s) ]", room, strings.Join(antLabels, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("[ %s ]", room))
		}
		if i < len(sim.Path)-1 {
			sb.WriteString(" ---> ")
		}
	}
	return sb.String()
}

// 7. WriteSimulationOutput
// Writes the complete simulation output (extra info and turn-by-turn details) to a file.
// This is invoked after the simulation loop ends.
func WriteSimulationOutput(filename string, extraInfo string, turnOutputs []string, totalTurns int) error {
	var sb strings.Builder
	// Write the extra info header.
	sb.WriteString(extraInfo)
	sb.WriteString("\n\n")
	// Append each turn's output.
	for i, turnOutput := range turnOutputs {
		sb.WriteString(fmt.Sprintf("TURN %d\n", i+1))
		sb.WriteString(turnOutput)
		sb.WriteString("\n")
	}
	// Write the total number of turns.
	sb.WriteString(fmt.Sprintf("Total turns: %d\n", totalTurns))
	// Write to file.
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// 8. PrintTerminalOutput
// Prints the minimal move information for each turn to the terminal.
// This is called after the simulation loop to display moves.
func PrintTerminalOutput(terminalOutputs []string) {
	for i, moves := range terminalOutputs {
		fmt.Printf("Turn %d: %s\n", i+1, moves)
	}
}
