# Lem-in

A digital ant farm simulation that finds the quickest way to move ants from start to end room.

## Objective

Create a program that reads an ant farm description from a file and simulates ants moving through tunnels to find the optimal path with minimum moves.

## Usage

```bash
go run . <input_file>
```

Example:
```bash
go run . examples/example01.txt
```

## Input Format

The input file contains:
1. **Number of ants** (first line)
2. **Rooms** with coordinates: `room_name x_coord y_coord`
   - `##start` marks the starting room
   - `##end` marks the ending room
3. **Tunnels** connecting rooms: `room1-room2`

### Example Input:
```
3
##start
1 23 3
2 16 7
3 16 3
##end
0 9 5
1-3
3-0
2-1
```

## Output Format

The program outputs:
1. **Input data** (echoed back)
2. **Movement turns** showing ant movements: `Lx-room_name`

### Example Output:
```
3
##start
1 23 3
2 16 7
3 16 3
##end
0 9 5
1-3
3-0
2-1

L1-3 L2-1
L1-0 L2-3 L3-1
L2-0 L3-3
L3-0
```

## Rules

- Ants start at `##start` and must reach `##end`
- Each room can hold only one ant (except start/end rooms)
- Each tunnel can only be used once per turn
- Room names cannot start with 'L' or '#' and must have no spaces
- Only standard Go packages allowed

## Error Handling

The program handles invalid input with error message: `ERROR: invalid data format`

Common error cases:
- Invalid number of ants
- Missing start or end room
- Duplicate rooms
- Invalid room coordinates
- Links to unknown rooms
- No path between start and end

## Project Structure

```
lem-in/
├── main.go                # Entry point
├── app/                   # Application logic
├── parser/                # Input file parsing
├── graph/                 # Graph construction and pathfinding
├── scheduling/            # Ant scheduling algorithm
├── simulation/            # Movement simulation
├── visualizer/            # Visualization output
└── structs/               # Shared data structures
```

