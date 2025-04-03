package scheduling

import (
	"lem-in/structs"
	"sort"
)

// AssignAnts distributes the ants among available paths using a greedy method.
// It calculates a "cost" for each path (path length + ants already assigned - 1)
// and assigns each ant to the path with the lowest cost.
// Parameters:
//   - antTotal (int): Total number of ants.
//   - pathList ([][]string): A list of paths (each path is a list of room names).
//
// Returns:
//   - structs.PathAssignment: Contains the list of paths and the number of ants assigned to each.
func AssignAnts(antTotal int, pathList [][]string) structs.PathAssignment {
	numPaths := len(pathList)
	antsOnPath := make([]int, numPaths)

	for i := 0; i < antTotal; i++ {
		// Local structure to hold the cost for a path.
		type pathCost struct {
			index int
			cost  int
		}

		costList := make([]pathCost, numPaths)
		// Calculate cost for each path: path length + ants already assigned - 1.
		for j := 0; j < numPaths; j++ {
			cost := len(pathList[j]) + antsOnPath[j] - 1
			costList[j] = pathCost{index: j, cost: cost}
		}

		// Sort paths by cost (lowest cost first).
		sort.Slice(costList, func(a, b int) bool {
			return costList[a].cost < costList[b].cost
		})

		// Assign the ant to the path with the lowest cost.
		antsOnPath[costList[0].index]++
	}
	return structs.PathAssignment{Paths: pathList, AntsPerPath: antsOnPath}
}
