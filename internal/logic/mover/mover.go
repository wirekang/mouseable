package mover

import (
	"math"

	"github.com/wirekang/mouseable/internal/di"
)

type Mover struct {
	speed      float64
	maxSpeed   float64
	direction  di.Direction
	directions []di.Direction
	isLocked   bool
}

func (m *Mover) SetMaxSpeed(v int) {
	m.maxSpeed = float64(v)
	if m.speed > float64(v) {
		m.speed = float64(v)
	}
}

func (m *Mover) SetSpeed(v float64) {
	if v > m.maxSpeed {
		m.speed = m.maxSpeed
		return
	}

	m.speed = v
}

func (m *Mover) AddSpeedIfDirection(v float64) {
	if m.direction == 0 {
		return
	}
	m.SetSpeed(m.speed + v)
}

func (m *Mover) AddDirection(dir di.Direction) {
	m.directions = append(m.directions, dir)
	m.calcDirection()
}

func (m *Mover) RemoveDirection(dir di.Direction) {
	rst := make([]di.Direction, len(m.directions))
	var find bool
	for i := range m.directions {
		if !find && m.directions[i] == dir {
			find = true
			continue
		}

		rst[i] = m.directions[i]
	}
	m.directions = rst
	m.calcDirection()
}
func (m *Mover) SetDirection(d di.Direction) {
	m.direction = d
}

func (m *Mover) Direction() di.Direction {
	return m.direction
}

func (m *Mover) Vector() (r VectorInt) {
	if m.direction == 0 {
		return
	}

	f := dirVectorMap[m.direction]
	r.X = int(math.Round(m.speed * f.X))
	r.Y = int(math.Round(m.speed * f.Y))
	return
}

func (m *Mover) calcDirection() {
	m.direction = 0
	for i := range m.directions {
		m.direction = m.direction | m.directions[i]
	}

	if m.direction == 0 {
		m.speed = 0
	}
}

type VectorInt struct {
	X, Y int
}

type VectorFloat struct {
	X, Y float64
}

const r = 1 / math.Sqrt2

var dirVectorMap = map[di.Direction]VectorFloat{
	di.DirectionUp:                       {+0, -1},
	di.DirectionDown:                     {+0, +1},
	di.DirectionRight:                    {+1, +0},
	di.DirectionRight | di.DirectionUp:   {+r, -r},
	di.DirectionRight | di.DirectionDown: {+r, +r},
	di.DirectionLeft:                     {-1, +0},
	di.DirectionLeft | di.DirectionUp:    {-r, -r},
	di.DirectionLeft | di.DirectionDown:  {-r, +r},
}
