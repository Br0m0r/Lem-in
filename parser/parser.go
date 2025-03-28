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
	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var antCount int
	var rooms []structs.Room
	var tunnels []structs.Tunnel

	if scanner.Scan() {
		countStr := strings.TrimSpace(scanner.Text())
		antCount, err = strconv.Atoi(countStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	var isNextRoomStart, isNextRoomEnd bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			if strings.HasPrefix(line, "##start") {
				isNextRoomStart = true
				continue
			}
			if strings.HasPrefix(line, "##end") {
				isNextRoomEnd = true
				continue
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}
			room := structs.Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: isNextRoomStart,
				IsEnd:   isNextRoomEnd,
			}
			rooms = append(rooms, room)
			isNextRoomStart = false
			isNextRoomEnd = false
			continue
		}
		if strings.Contains(line, "-") {
			roomNames := strings.Split(line, "-")
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			tunnels = append(tunnels, structs.Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]})
			continue
		}
	}

	return antCount, rooms, tunnels, nil
}
