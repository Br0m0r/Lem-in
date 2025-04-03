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

// ParseInputFile reads an input file and returns:
//   - The total number of ants,
//   - A list of rooms,
//   - A list of tunnels (connections between rooms).
//
// Example using example00.txt:
//
//	File content:
//	  4
//	  ##start
//	  0 0 3
//	  2 2 5
//	  3 4 0
//	  ##end
//	  1 8 3
//	  0-2
//	  2-3
//	  3-1
//
//	From this file:
//	  - The ant count is 4.
//	  - Rooms parsed:
//	      "0" with coordinates (0,3), marked as start.
//	      "2" with coordinates (2,5).
//	      "3" with coordinates (4,0).
//	      "1" with coordinates (8,3), marked as end.
//	  - Tunnels parsed:
//	      "0-2", "2-3", "3-1".
func ParseInputFile(filePath string) (int, []structs.Room, []structs.Tunnel, error) {
	// Open the input file.
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	// Ensure the file is closed after reading.
	defer file.Close()

	// Create a scanner to read the file line by line.
	scanner := bufio.NewScanner(file)

	// Variables to store our results.
	var antTotal int                // Will hold the number of ants.
	var roomList []structs.Room     // Will accumulate room definitions.
	var tunnelList []structs.Tunnel // Will accumulate tunnel definitions.

	// Create a map to check for duplicate room positions.
	// We use "x,y" as a key to ensure no two rooms share the same coordinates.
	positionMap := make(map[string]bool)

	// Read the first line to get the ant count.
	// For example00.txt, the first line is "4".
	if scanner.Scan() {
		antStr := strings.TrimSpace(scanner.Text())
		antTotal, err = strconv.Atoi(antStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	// Flags to mark the next room as the start or the end.
	var nextRoomIsStart, nextRoomIsEnd bool

	// Process the remaining lines.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue // Skip empty lines.
		}
		// Handle comment lines:
		// Lines beginning with '#' are either comments or commands.
		// For example, "##start" and "##end" mark the following room.
		if line[0] == '#' {
			if strings.HasPrefix(line, "##start") {
				nextRoomIsStart = true // Mark the next room as the starting room.
				continue
			}
			if strings.HasPrefix(line, "##end") {
				nextRoomIsEnd = true // Mark the next room as the ending room.
				continue
			}
			// Ignore other comments.
			continue
		}
		// If the line has exactly 3 parts, it's a room definition.
		// In example00.txt, the room definitions are like "0 0 3" or "2 2 5".
		parts := strings.Fields(line)
		if len(parts) == 3 {
			// Convert the coordinate strings into integers.
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}
			// Create a unique key for this room's position.
			posKey := fmt.Sprintf("%d,%d", x, y)
			if _, exists := positionMap[posKey]; exists {
				return 0, nil, nil, errors.New("ERROR: duplicate room coordinates")
			}
			positionMap[posKey] = true // Mark these coordinates as used.

			// Create the room with the provided data.
			// For example, for "0 0 3": Name is "0", X is 0, Y is 3, and it is marked as start if flagged.
			newRoom := structs.Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: nextRoomIsStart,
				IsEnd:   nextRoomIsEnd,
			}
			// Add the room to our room list.
			roomList = append(roomList, newRoom)

			// Reset the start/end flags after using them.
			nextRoomIsStart = false
			nextRoomIsEnd = false
			continue
		}
		// If the line contains a hyphen, it's a tunnel definition.
		// In example00.txt, these are "0-2", "2-3", and "3-1".
		if strings.Contains(line, "-") {
			roomNames := strings.Split(line, "-")
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			newTunnel := structs.Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]}
			tunnelList = append(tunnelList, newTunnel)
			continue
		}
	}
	// Return the ant count, room list, and tunnel list.
	return antTotal, roomList, tunnelList, nil
}
