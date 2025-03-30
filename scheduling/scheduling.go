package scheduling

import (
	"lem-in/structs"
	"sort"
)

// 1. AssignAnts
// Distributes the ants among the available paths using a greedy algorithm.
// It calculates an "effective cost" for each path and assigns each ant to the path with the lowest cost.
//
// Arguments:
//   - antCount (int): The total number of ants to be assigned.
//   - paths ([][]string): A slice of available paths, where each path is a slice of room names.
//
// Returns:
//   - structs.PathAssignment: A structure containing:
//   - Paths: the original paths,
//   - AntsPerPath: a slice of integers where each element represents the number of ants assigned to the corresponding path.
func AssignAnts(antCount int, paths [][]string) structs.PathAssignment {
	// numPaths holds the total number of available paths.
	numPaths := len(paths)
	// antsPerPath is a slice to count the ants assigned to each path. Initially, all counts are 0.
	antsPerPath := make([]int, numPaths)

	// Process each ant one by one.
	for i := 0; i < antCount; i++ {
		// Define a local structure to hold the cost information for a path.
		// 'index' is the path's index, and 'cost' is its effective cost.
		type pathCost struct {
			index int
			cost  int
		}

		// pathCosts will hold the cost for each available path.
		pathCosts := make([]pathCost, numPaths)
		// For each path, calculate its effective cost.
		// Effective cost = (path length) + (number of ants already assigned) - 1.
		// The subtraction of 1 adjusts for the starting room not counting as a "move".
		for j := 0; j < numPaths; j++ {
			cost := len(paths[j]) + antsPerPath[j] - 1
			pathCosts[j] = pathCost{index: j, cost: cost}
		}

		// Sort the pathCosts slice so that the path with the smallest cost comes first.
		sort.Slice(pathCosts, func(a, b int) bool {
			return pathCosts[a].cost < pathCosts[b].cost
		})
		// Assign the current ant to the path with the lowest effective cost.
		antsPerPath[pathCosts[0].index]++
	}
	// Return the assignment containing the paths and the count of ants assigned to each path.
	return structs.PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
