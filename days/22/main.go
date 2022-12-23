package main

import (
	//	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Direction int

const (
	EAST  Direction = 0
	SOUTH Direction = 1
	WEST  Direction = 2
	NORTH Direction = 3
)

type Instruction struct {
	distance int
	turn     string
}

type RoomBorder struct {
	room *Room
	face Direction
}

type FoldableRoomBorder struct {
	room         *Room
	face         Direction
	originalFace Direction
}

type Room struct {
	offsetX     int
	offsetY     int
	layout      []string
	connections map[Direction]RoomBorder
}

type Board struct {
	height       int
	width        int
	rooms        [][]*Room
	instructions []Instruction
	roomSize     int
}

type Player struct {
	x      int
	y      int
	facing Direction
}

var faceNames = map[Direction]string{
	NORTH: "NORTH",
	EAST:  "EAST",
	SOUTH: "SOUTH",
	WEST:  "WEST",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])

	board, player := parseInput(input, 50)
	board.fold()
	board.execute(&player)

	password := (player.y+1)*1000 + (player.x+1)*4 + int(player.facing)
	fmt.Printf("password: %d\n", password)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string, roomSize int) (Board, Player) {
	components := strings.Split(input, "\n\n")

	rooms := parseRooms(components[0], roomSize)
	instructions := parseInstructions(components[1])

	// find player starting position
	firstRoomX := 0
	for rooms[0][firstRoomX] == nil {
		firstRoomX++
	}

	return Board{
			height:       len(rooms) * roomSize,
			width:        len(rooms[0]) * roomSize,
			rooms:        rooms,
			instructions: instructions,
			roomSize:     roomSize,
		}, Player{
			x:      firstRoomX * roomSize,
			y:      0,
			facing: EAST,
		}
}

func parseRooms(input string, roomSize int) [][]*Room {
	lines := strings.Split(input, "\n")

	rooms := make([][]*Room, len(lines)/roomSize)

	maxLineLength := 0
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) > maxLineLength {
			maxLineLength = len(lines[i])
		}
	}

	for i := 0; i < len(rooms); i++ {
		rooms[i] = make([]*Room, maxLineLength/roomSize)
	}

	for roomY := 0; roomY < len(lines); roomY += roomSize {
		topLine := lines[roomY]
		start := len(topLine) - len(strings.TrimLeft(topLine, " "))
		end := len(topLine)
		roomsOnLine := (end - start) / roomSize

		for i := 0; i < roomsOnLine; i++ {
			roomX := start + (i * roomSize)
			roomLayout := make([]string, roomSize)
			for yOffset := 0; yOffset < roomSize; yOffset++ {
				roomLayout[yOffset] = lines[roomY+yOffset][roomX : roomX+roomSize]
			}
			room := &Room{
				offsetX:     roomX,
				offsetY:     roomY,
				layout:      roomLayout,
				connections: make(map[Direction]RoomBorder),
			}
			// ...
			rooms[roomY/roomSize][roomX/roomSize] = room
		}
	}

	return rooms
}

func parseInstructions(input string) []Instruction {
	r := regexp.MustCompile(`(\d+)[LR]?`)

	res := r.FindAllString(input, -1)

	instructions := make([]Instruction, len(res))
	for i, match := range res {
		lastChar := match[len(match)-1:]

		var instr Instruction
		if lastChar == "L" || lastChar == "R" {
			distance, _ := strconv.Atoi(match[:len(match)-1])
			instr = Instruction{distance, lastChar}
		} else {
			distance, _ := strconv.Atoi(match)
			instr = Instruction{distance, ""}
		}
		instructions[i] = instr
	}

	return instructions
}

func (b *Board) fold() {
	// The general idea:
	// Starting with a 1D list of outer edges: we can recursively join adjacent
	// edges that form a concave right angle (creating a "warp" between faces) by
	// "rotating" one edge into the other (and by extension, all the edges that
	// follow) by 90 degrees. This should eventually leave us with an empty list.
	//
	// We will know that the angle is concave by walking the edges in a clockwise
	// direction, and looking for counter-clockwise direction changes.

	edges := b.generatePerimeter()

	edgeIndex := 0
	for len(edges) > 0 {

		edge1Index := edgeIndex
		edge2Index := edgeIndex + 1
		// Walk the perimeter and find edges that turn counter-clockwise
		edge1 := edges[edge1Index]
		edge2 := edges[edge2Index]
		turn := ((edge1.face - edge2.face + 4) % 4)
		if turn == 1 {
			// this is a counter-clockwise turn, join the two edges
			edge1.room.connections[edge1.originalFace] = RoomBorder{edge2.room, edge2.originalFace}
			edge2.room.connections[edge2.originalFace] = RoomBorder{edge1.room, edge1.originalFace}

			edges = append(edges[:edge1Index], edges[edge1Index+2:]...)

			if len(edges) > 0 {

				// rotate the remaining edges
				for i := edgeIndex; i < len(edges); i++ {
					edges[i].face = (edges[i].face - 1 + 4) % 4
				}

				edgeIndex = edgeIndex % len(edges)

			}
		} else {
			edgeIndex++
		}
		if edgeIndex >= len(edges)-1 {
			edgeIndex = 0
		}
	}
}

func (b *Board) generatePerimeter() []FoldableRoomBorder {
	edges := make([]FoldableRoomBorder, 0)

	seen := make(map[FoldableRoomBorder]bool)

	var currentRoom *Room
	for i := 0; i < len(b.rooms[0]); i++ {
		if b.rooms[0][i] != nil {
			currentRoom = b.rooms[0][i]
			break
		}
	}

	edge := FoldableRoomBorder{currentRoom, NORTH, NORTH}
	edges = append(edges, edge)
	seen[FoldableRoomBorder{currentRoom, NORTH, NORTH}] = true

	currentFace := NORTH

	for {
		// Try the next face on the current room. if it's an inner face, walk one
		// room in the direction of the next face
		nextFace := (currentFace + 1) % 4
		nextRoom := currentRoom
		var nextX, nextY int

		switch nextFace {
		case NORTH:
			{
				// Check the room above (at the top left corner)
				nextX = currentRoom.offsetX
				nextY = currentRoom.offsetY - 1
			}
		case WEST:
			{
				// Check the room to the left (at the bottom left corner)
				nextX = currentRoom.offsetX - 1
				nextY = currentRoom.offsetY + b.roomSize - 1
			}
		case SOUTH:
			{
				// Check the room below (at the bottom right corner)
				nextX = currentRoom.offsetX + b.roomSize - 1
				nextY = currentRoom.offsetY + b.roomSize
			}
		case EAST:
			{
				// Check the room to the right (at the top right corner)
				nextX = currentRoom.offsetX + b.roomSize
				nextY = currentRoom.offsetY
			}
		}

		// if the position is inside a room, it's an inner face
		if b.isInBounds(nextX, nextY) {
			nextFace = currentFace
			nextRoom, _, _ = b.getRoomAt(nextX, nextY)
			// check one more point, to determine if it's an interior corner. if it
			// is, we rotate the face counter-clockwise and test the room again
			switch nextFace {
			case NORTH:
				{
					// Check for a room diagonally up and to the right
					nextY -= 1
				}
			case EAST:
				{
					// Check for a room diagonally down and to the right
					nextX += 1
				}
			case SOUTH:
				{
					// Check for a room diagonally down and to the left
					nextY += 1
				}
			case WEST:
				{
					// Check for a room diagonally up and to the left
					nextX -= 1
				}
			}
			if b.isInBounds(nextX, nextY) {
				nextFace = (nextFace - 1 + 4) % 4
				nextRoom, _, _ = b.getRoomAt(nextX, nextY)
			}
		}

		edge := FoldableRoomBorder{nextRoom, nextFace, nextFace}
		if seen[edge] {
			// The perimeter is complete, we can exit
			break
		}

		edges = append(edges, edge)
		seen[edge] = true
		currentRoom = nextRoom
		currentFace = nextFace
	}

	return edges
}

func (b *Board) execute(player *Player) {
	for _, instruction := range b.instructions {
		b.movePlayer(player, instruction)
	}
}

func (b *Board) movePlayer(player *Player, instruction Instruction) {
	isInRoomBounds := func(x, y int) bool {
		return x >= 0 && x < b.roomSize && y >= 0 && y < b.roomSize
	}
	// assume the player always starts at a valid position on the board (on an
	// open room tile, in-bounds)
	room, roomX, roomY := b.getRoomAt(player.x, player.y)
	facing := player.facing
	nextX := roomX
	nextY := roomY
	nextRoom := room
	nextFacing := facing

	for i := 0; i < instruction.distance; i++ {
		switch facing {
		case EAST:
			{
				nextX++
			}
		case SOUTH:
			{
				nextY++
			}
		case WEST:
			{
				nextX--
			}
		case NORTH:
			{
				nextY--
			}
		}
		if !isInRoomBounds(nextX, nextY) {
			// Traverse to the connected room. First, test if there is a room
			// there naturally. If not, look up the connected face
			nextRoom, nextX, nextY = b.getRoomAt(room.offsetX+nextX, room.offsetY+nextY)
			if nextRoom == nil {
				connection := room.connections[facing]
				nextRoom = connection.room
				newPlayer := b.findConnectedPositionAndOrientation(connection, Player{roomX, roomY, facing})
				nextX = newPlayer.x
				nextY = newPlayer.y
				nextFacing = newPlayer.facing
			}
		}
		switch nextRoom.getTileAt(nextX, nextY) {
		case '#':
			{
				// we hit an obstacle, stop moving
				nextX = roomX
				nextY = roomY
				nextFacing = facing
				nextRoom = room
				i = instruction.distance
			}
		}

		// commit the movement (if there is any)
		roomY = nextY
		roomX = nextX
		room = nextRoom
		facing = nextFacing
	}

	player.x = room.offsetX + roomX
	player.y = room.offsetY + roomY

	switch instruction.turn {
	case "L":
		{
			player.facing = (facing - 1 + 4) % 4
		}
	case "R":
		{
			player.facing = (facing + 1) % 4
		}
	}
}

// // find the room at the given y level, wrapping back to the top if y exceeds
// // the total board height
// func (b *Board) getRoom(y int) (*Room, int) {
// 	boardY := (y + b.height) % b.height
// 	for i := 0; i < len(b.rooms); i++ {
// 		room := &b.rooms[i]
// 		if boardY >= room.offsetY && boardY < room.offsetY+room.height {
// 			return room, boardY - room.offsetY
// 		}
// 	}
// 	panic(errors.New("room not found somehow?"))
// }

func (b *Board) isInBounds(x, y int) bool {
	if y < 0 || y >= b.height || x < 0 || x >= b.width {
		return false
	}
	room, _, _ := b.getRoomAt(x, y)
	return room != nil
}

func (r *Room) getTileAt(x, y int) byte {
	return r.layout[y][x]
}

func (b *Board) findConnectedPositionAndOrientation(connection RoomBorder, player Player) Player {
	outgoingDirection := (connection.face + 2) % 4

	newX := player.x
	newY := player.y

	// There's probably an elegant mathematical way to compute these, but it's late and I'm lazy
	switch player.facing {
	case EAST:
		{
			switch connection.face {
			case NORTH:
				{
					newX = b.roomSize - 1 - player.y
					newY = 0
				}
			case EAST:
				{
					newY = b.roomSize - 1 - player.y
				}
			case SOUTH:
				{
					newX = player.y
					newY = b.roomSize - 1
				}
			case WEST:
				{
					newX = 0
				}
			}
		}
	case SOUTH:
		{
			switch connection.face {
			case EAST:
				{
					newX = b.roomSize - 1
					newY = player.x
				}
			case SOUTH:
				{
					newX = b.roomSize - 1 - player.x
				}
			case WEST:
				{
					newX = 0
					newY = b.roomSize - 1 - player.x
				}
			case NORTH:
				{
					newY = 0
				}
			}
		}
	case WEST:
		{
			switch connection.face {
			case NORTH:
				{
					newX = player.y
					newY = 0
				}
			case WEST:
				{
					newY = b.roomSize - 1 - player.y
				}
			case SOUTH:
				{
					newX = b.roomSize - 1 - player.y
					newY = b.roomSize - 1
				}
			case EAST:
				{
					newX = b.roomSize - 1
				}
			}
		}
	case NORTH:
		{
			switch connection.face {
			case EAST:
				{
					newX = b.roomSize - 1
					newY = b.roomSize - 1 - player.x
				}
			case NORTH:
				{
					newX = b.roomSize - 1 - player.x
				}
			case WEST:
				{
					newX = 0
					newY = player.x
				}
			case SOUTH:
				{
					newY = b.roomSize - 1
				}
			}
		}
	}

	return Player{newX, newY, outgoingDirection}
}

func (b *Board) getRoomAt(x, y int) (*Room, int, int) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return nil, 0, 0
	}
	return b.rooms[y/b.roomSize][x/b.roomSize], x % b.roomSize, y % b.roomSize
}

func (b *Board) print(player Player) string {
	var printout string
	// print headers. assume all headers are 3 digits at most
	rangeAxisHeight := len(fmt.Sprint(b.width)) + 1
	depthAxisLength := len(fmt.Sprint(b.height)) + 1
	for y := 0; y < rangeAxisHeight; y++ {
		var line string
		for x := -(depthAxisLength + 1); x < b.width+1; x++ {
			xStr := fmt.Sprint(x)
			if x >= 0 && x%5 == 0 {
				digit := len(xStr) - (rangeAxisHeight - y)
				if digit >= 0 {
					line += string(xStr[digit])
				} else {
					line += " "
				}
			} else {
				line += " "
			}
		}
		printout += fmt.Sprintln(line)
	}

	for y := 0; y < b.height; y++ {
		var line string
		currentDepthSize := len(fmt.Sprint(y))
		for i := 0; i < depthAxisLength-currentDepthSize; i++ {
			line += " "
		}
		line += fmt.Sprintf("%d ", y)

		for x := 0; x < b.width; x++ {
			room, roomX, roomY := b.getRoomAt(x, y)

			if room == nil {
				line += " "
				continue
			}

			if player.y == y && player.x == x {
				playerRune := "?"
				switch player.facing {
				case EAST:
					{
						playerRune = ">"
					}
				case SOUTH:
					{
						playerRune = "v"
					}
				case WEST:
					{
						playerRune = "<"
					}
				case NORTH:
					{
						playerRune = "^"
					}
				}
				line += playerRune
			} else {
				line += string(room.getTileAt(roomX, roomY))
			}
		}
		printout += fmt.Sprintln(line)
	}

	return printout
}
