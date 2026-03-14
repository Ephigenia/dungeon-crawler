package game

// Player holds the player's state and stats.
type Player struct {
	X, Y    int
	HP      int
	MaxHP   int
	Attack  int
	Defense int
}

func newPlayer(x, y int) *Player {
	return &Player{
		X:       x,
		Y:       y,
		HP:      30,
		MaxHP:   30,
		Attack:  5,
		Defense: 2,
	}
}

// IsAlive returns true if the player has HP remaining.
func (p *Player) IsAlive() bool {
	return p.HP > 0
}

// TakeDamage reduces HP by the incoming attack minus defense, minimum 1.
func (p *Player) TakeDamage(attack int) {
	dmg := attack - p.Defense
	if dmg < 1 {
		dmg = 1
	}
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}
