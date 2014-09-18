package main

type Point struct {
	X int
	Y int
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

func (d Direction) Reverse() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	default:
		return North
	}
}

type Exits map[Direction]*Room

type Room struct {
	Name        string
	Description string
	Position    Point
	Exits       Exits
	Dungeon     *Dungeon
}

type Dungeon struct {
	Rooms map[Point]*Room
}

// OpenExit creates a connection from the receiver's Room to another Room. If
// a Room already exists in the Direction specified, Exits on both Rooms are
// linked. If no Room exists, a new one is created and added to the Dungeon.
// If an exit is already open in that Direction, returns false.
func (r *Room) OpenExit(direction Direction) bool {
	if _, present := r.Exits[direction]; present {
		return false
	}

	changeX := 0
	changeY := 0

	switch direction {
	case North:
		changeY = -1
	case South:
		changeY = 1
	case East:
		changeX = 1
	case West:
		changeX = -1
	default:
		return false
	}

	destinationPosition := Point{
		X: r.Position.X + changeX,
		Y: r.Position.Y + changeY,
	}

	if destinationRoom, roomPresent := r.Dungeon.Rooms[destinationPosition]; roomPresent {
		r.Exits[direction] = destinationRoom
		destinationRoom.Exits[direction.Reverse()] = r
	}

	return true
}
