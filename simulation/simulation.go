package simulation

import (
	"fmt"
	"os"
	"strings"

	"lem-in/structs"
)

const NotInjected = -1

// GeneratePathGrid creates a 2D grid visualization for a single simulation path.
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

func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}

type simulationLogger struct {
	terminal strings.Builder
	file     strings.Builder
}

func (l *simulationLogger) WriteTurn(turn int, moves []string, grid string) {
	l.file.WriteString(fmt.Sprintf("TURN %d\n%s\n\n", turn+1, grid))
	l.terminal.WriteString(fmt.Sprintf("Turn %d: %s\n", turn+1, strings.Join(moves, " ")))
}

func initSimStates(paths [][]string, assignment structs.PathAssignment) []structs.PathSim {
	sims := make([]structs.PathSim, len(paths))
	antCounter := 1
	for i, p := range paths {
		count := assignment.AntsPerPath[i]
		positions := make([]int, count)
		for j := range positions {
			positions[j] = NotInjected
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

func processPath(sim structs.PathSim, newPos []int, turnMoves *[]string) {
	pathLen := len(sim.Path)
	if pathLen == 2 {
		for j := 0; j < len(sim.Positions); j++ {
			if sim.Positions[j] == NotInjected {
				newPos[j] = 1
				*turnMoves = append(*turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
				break
			}
		}
	} else {
		for j := len(sim.Positions) - 1; j >= 0; j-- {
			pos := sim.Positions[j]
			if pos == NotInjected {
				if !isOccupied(newPos, 1) {
					newPos[j] = 1
					*turnMoves = append(*turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
				}
			} else if pos < pathLen-1 {
				next := pos + 1
				if next == pathLen-1 || !isOccupied(newPos, next) {
					newPos[j] = next
					*turnMoves = append(*turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
				}
			}
		}
	}
}

func SimulateMultiPath(antCount int, paths [][]string, assignment structs.PathAssignment, extraInfo string) {
	sims := initSimStates(paths, assignment)
	logger := &simulationLogger{}

	outFile, err := os.Create("simulation_output.txt")
	if err != nil {
		fmt.Println("Error creating simulation_output.txt:", err)
		return
	}
	defer outFile.Close()

	logger.file.WriteString(extraInfo + "\n\n")
	turn := 0

	for {
		var turnMoves []string
		var turnGrid strings.Builder

		for idx, sim := range sims {
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)
			processPath(sim, newPos, &turnMoves)
			copy(sims[idx].Positions, newPos)
			turnGrid.WriteString(GeneratePathGrid(sims[idx]) + "\n")
		}

		if len(turnMoves) == 0 {
			break
		}

		logger.WriteTurn(turn, turnMoves, turnGrid.String())
		turn++
	}

	logger.file.WriteString(fmt.Sprintf("Total turns: %d\n", turn))
	_, err = outFile.WriteString(logger.file.String())
	if err != nil {
		fmt.Println("Error writing to simulation_output.txt:", err)
	}

	fmt.Print(logger.terminal.String())
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
