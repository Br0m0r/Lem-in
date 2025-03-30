package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/structs"
)

// ParseInputFile reads the input file and returns the ant count, rooms, tunnels, and an error if any.
// It expects the first line to be the ant count, then room definitions (with "##start" and "##end"),
// followed by tunnel definitions.
func ParseInputFile(filename string) (int, []structs.Room, []structs.Tunnel, error) {
	// 1. Open the input file.
	//    - filename (string): the path to the input file.
	//    - Returns an *os.File handle if successful.
	file, err := os.Open(filename)
	if err != nil {
		// If the file cannot be opened, return an error.
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	// Ensure the file is closed once processing is done.
	defer file.Close()

	// 2. Create a scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	// 3. Declare variables to store results.
	var antCount int             // Will hold the number of ants (parsed from the first line).
	var rooms []structs.Room     // A slice to accumulate room definitions.
	var tunnels []structs.Tunnel // A slice to accumulate tunnel definitions.

	// 4. Create a map to check for duplicate room coordinates.
	//    - Key format is "x,y" (e.g., "10,20").
	coordinateMap := make(map[string]bool)

	// 5. Read the first line for the ant count.
	if scanner.Scan() {
		countStr := strings.TrimSpace(scanner.Text())
		// Convert the ant count from string to integer.
		antCount, err = strconv.Atoi(countStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	// 6. Variables to flag if the next room should be marked as start or end.
	var isNextRoomStart, isNextRoomEnd bool

	// 7. Process the rest of the file line by line.
	for scanner.Scan() {
		// Remove any surrounding whitespace.
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines.
		if len(line) == 0 {
			continue
		}
		// 8. Handle comment lines:
		//    - Lines starting with "#" are comments or commands.
		if line[0] == '#' {
			// If the line is "##start", flag the next room as the start.
			if strings.HasPrefix(line, "##start") {
				isNextRoomStart = true
				continue
			}
			// If the line is "##end", flag the next room as the end.
			if strings.HasPrefix(line, "##end") {
				isNextRoomEnd = true
				continue
			}
			// Other comments are ignored.
			continue
		}
		// 9. Split the line into parts.
		parts := strings.Fields(line)
		// If there are exactly 3 parts, it is treated as a room definition.
		if len(parts) == 3 {
			// Convert the second and third parts into integers (x and y coordinates).
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}

			// 10. Check for duplicate coordinates.
			//     Create a key using the format "x,y".
			coordKey := fmt.Sprintf("%d,%d", x, y)
			if _, exists := coordinateMap[coordKey]; exists {
				return 0, nil, nil, errors.New("ERROR: duplicate room coordinates")
			}
			// Record these coordinates.
			coordinateMap[coordKey] = true

			// 11. Create a room structure with the name, coordinates, and flags for start/end.
			room := structs.Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: isNextRoomStart,
				IsEnd:   isNextRoomEnd,
			}
			// Append the room to the rooms slice.
			rooms = append(rooms, room)
			// Reset the flags for the next room.
			isNextRoomStart = false
			isNextRoomEnd = false
			// Move to the next line.
			continue
		}
		// 12. If the line contains a hyphen, treat it as a tunnel definition.
		if strings.Contains(line, "-") {
			// Split the line by "-" into two room names.
			roomNames := strings.Split(line, "-")
			// There must be exactly 2 parts for a valid tunnel.
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			// Create a tunnel structure and add it to the tunnels slice.
			tunnels = append(tunnels, structs.Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]})
			continue
		}
	}
	// 13. Return the parsed ant count, rooms, tunnels, and nil error if all goes well.
	return antCount, rooms, tunnels, nil
}
