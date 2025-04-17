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

// ParseInputFile reads ant count, rooms, and tunnels from an input file.
func ParseInputFile(filePath string) (int, []structs.Room, []structs.Tunnel, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// read number of ants
	var antTotal int
	if scanner.Scan() {
		antStr := strings.TrimSpace(scanner.Text())
		antTotal, err = strconv.Atoi(antStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	var roomList []structs.Room
	var tunnelList []structs.Tunnel
	positionMap := make(map[string]bool)
	var nextRoomIsStart, nextRoomIsEnd bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "##start") {
			nextRoomIsStart = true
			continue
		}
		if strings.HasPrefix(line, "##end") {
			nextRoomIsEnd = true
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}
			posKey := fmt.Sprintf("%d,%d", x, y)
			if positionMap[posKey] {
				return 0, nil, nil, errors.New("ERROR: duplicate room coordinates")
			}
			positionMap[posKey] = true
			newRoom := structs.Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: nextRoomIsStart,
				IsEnd:   nextRoomIsEnd,
			}
			roomList = append(roomList, newRoom)
			nextRoomIsStart = false
			nextRoomIsEnd = false
			continue
		}
		if strings.Contains(line, "-") {
			names := strings.Split(line, "-")
			tunnelList = append(tunnelList, structs.Tunnel{RoomA: names[0], RoomB: names[1]})
		}
	}

	return antTotal, roomList, tunnelList, nil
}
