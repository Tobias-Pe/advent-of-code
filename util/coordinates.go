package util

import "fmt"

type Coordinate struct {
	I, J int
}

func (c Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.I, c.J)
}

func (c Coordinate) GetNeighbours() []Coordinate {
	return []Coordinate{{c.I + 1, c.J}, {c.I - 1, c.J}, {c.I, c.J + 1}, {c.I, c.J - 1}}
}

func (c Coordinate) IsValid(iMax, jMax int) bool {
	return c.I >= 0 && c.J >= 0 && c.I < iMax && c.J < jMax
}

func (c Coordinate) Add(dI, dJ int) Coordinate {
	return Coordinate{c.I + dI, c.J + dJ}
}

func (c Coordinate) Left() Coordinate {
	return Coordinate{c.I, c.J - 1}
}

func (c Coordinate) Right() Coordinate {
	return Coordinate{c.I, c.J + 1}
}

func (c Coordinate) Up() Coordinate {
	return Coordinate{c.I - 1, c.J}
}

func (c Coordinate) Down() Coordinate {
	return Coordinate{c.I + 1, c.J}
}
