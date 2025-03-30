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

// Run is the entry point for the application workflow.
func Run() {
	// Check that at least one command-line argument (the input file path) is provided.
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}

	// inputFile holds the file path provided as the first command-line argument.
	inputFile := os.Args[1]

	// ParseInputFile reads the input file and returns the ant count, rooms, tunnels, and error if any.
	antCount, rooms, tunnels, err := parser.ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate that there is at least one ant.
	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Verify that both a start room and an end room exist.
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

	// Build the graph of the ant farm using rooms and tunnels.
	antFarmGraph, err := graph.BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Find multiple edge-disjoint paths from the start to the end room.
	paths, err := graph.FindMultiplePaths(antFarmGraph)
	if err != nil || len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Use a greedy algorithm to assign ants to the found paths.
	assignment := scheduling.AssignAnts(antCount, paths)

	// Generate extra information about the simulation (input data summary, paths, etc.)
	extraInfo := visualizer.PrintExtraInfo(antCount, rooms, tunnels, paths, assignment)

	// Run the simulation of ant movements along the paths.
	simulation.SimulateMultiPath(antCount, paths, assignment, extraInfo)
}
