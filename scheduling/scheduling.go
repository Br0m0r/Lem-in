package scheduling

import (
	"lem-in/structs"
	"sort"
)

// AssignAnts distributes the ants among the available paths using a simple greedy algorithm.
// It calculates an "effective cost" for each path, which is defined as:
//
//	cost = (path length) + (number of ants already assigned) - 1
//
// The subtraction of 1 is to adjust for the fact that the starting room does not count as a move.
// Then, it assigns each ant to the path with the lowest cost.
//
// For example, using example00.txt:
//   - There are 4 ants.
//   - The only available path is ["0", "2", "3", "1"].
//   - The length of the path is 4 (there are 4 rooms).
//   - Initially, no ants are assigned, so the cost is: 4 + 0 - 1 = 3.
//   - Since there is only one path, every ant will be assigned to it, resulting in AntsPerPath: [4].
//
// Parameters:
//
//	antCount (int): The total number of ants, e.g., 4 from example00.txt.
//	paths ([][]string): A slice of available paths. In example00.txt, this might be [["0", "2", "3", "1"]].
//
// Returns:
//
//	structs.PathAssignment: Contains the original paths and a slice indicating how many ants are assigned to each path.
func AssignAnts(antCount int, paths [][]string) structs.PathAssignment {
	// numPaths holds the total number of available paths.
	numPaths := len(paths)

	// antsPerPath will count how many ants are assigned to each path.
	// For our example00.txt with one path, this slice starts as [0].
	antsPerPath := make([]int, numPaths)

	// Process each ant one by one.
	// For example, we have 4 ants so the loop will run 4 times.
	for i := 0; i < antCount; i++ {
		// Define a local structure to hold the cost for each path.
		// 'index' is the position of the path in the slice.
		// 'cost' is the computed effective cost for that path.
		type pathCost struct {
			index int
			cost  int
		}

		// Create a slice to store the cost for each available path.
		pathCosts := make([]pathCost, numPaths)

		// Calculate the effective cost for each path.
		// For each path, cost = (length of the path) + (ants already assigned) - 1.
		// In example00.txt, for our only path, initially cost = 4 (length) + 0 - 1 = 3.
		for j := 0; j < numPaths; j++ {
			cost := len(paths[j]) + antsPerPath[j] - 1
			pathCosts[j] = pathCost{index: j, cost: cost}
		}

		// Sort the pathCosts slice so that the path with the smallest cost comes first.
		// With one path, no real sorting is needed, but this is useful when there are multiple paths.
		sort.Slice(pathCosts, func(a, b int) bool {
			return pathCosts[a].cost < pathCosts[b].cost
		})

		// Assign the current ant to the path with the lowest effective cost.
		// For example, every ant (all 4 ants) will be assigned to our only path.
		antsPerPath[pathCosts[0].index]++
	}

	// Return the assignment containing:
	//   - Paths: the original paths slice (e.g., [["0", "2", "3", "1"]])
	//   - AntsPerPath: the number of ants assigned to each path. In our example, this will be [4].
	return structs.PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
