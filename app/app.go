package app

import (
	"fmt"
	"os"

	"lem-in/graph"
	"lem-in/parser"
	"lem-in/scheduling"
	"lem-in/simulation"
	"lem-in/visualizer"
)

// Run is the application entry point.
func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Parse input, build graph, assign ants, etc.
	antCount, rooms, tunnels, err := parser.ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	var startFound, endFound bool
	for _, room := range rooms {
		if room.IsStart {
			startFound = true
		}
		if room.IsEnd {
			endFound = true
		}
	}
	if !startFound || !endFound {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	antFarmGraph, err := graph.BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	paths, err := graph.FindMultiplePaths(antFarmGraph)
	if err != nil || len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	assignment := scheduling.AssignAnts(antCount, paths)
	extraInfo := visualizer.PrintExtraInfo(antCount, rooms, tunnels, paths, assignment)

	// Run the simulation.
	simulation.SimulateMultiPath(antCount, paths, assignment, extraInfo)
}
